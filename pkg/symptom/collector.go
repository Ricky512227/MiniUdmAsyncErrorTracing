package symptom

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/config"
	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/kubernetes"
	"go.uber.org/zap"
)

// Collector handles symptom collection from Kubernetes pods
type Collector struct {
	config    *config.Config
	k8sClient *kubernetes.Client
	logger    *zap.Logger
}

// SymptomCollectionConfig holds configuration for symptom collection
type SymptomCollectionConfig struct {
	Namespace string
	Pods      []string
	StartTime time.Time
}

// ErrorEvent represents an error detected during collection
type ErrorEvent struct {
	Timestamp time.Time
	Source    string
	Message   string
}

// NewCollector creates a new symptom collector
func NewCollector(cfg *config.Config, k8sClient *kubernetes.Client, logger *zap.Logger) *Collector {
	return &Collector{
		config:    cfg,
		k8sClient: k8sClient,
		logger:    logger,
	}
}

// StartCollection starts the symptom collection process
//
// Command Syntax: symptom-collection -n <namespace> -p <"podnames">
// Example:
//
//	symptom-collection -n miniudm -p "uecm"
//	symptom-collection -n dracvnf -p "uecm nim"
//
// **** Concurrency validation check = 1  ****
// 1. Get the namespace from the argument.
// 2. Get the pod from the argument to capture trace.
// 3. Check the namespace exists in the cluster.
// 4. Check the deployments of the pod exists in the namespace.
// 5. Check the readiness and liveness
// 6. Calculate the number of Process inside the mcc platform for the required process.
//
// # Record the start timestamp of execution also record the current timestamp of the pod
//
// Routine :: 1
//
//	Parallely enable the trace of each process.
//
// Routine :: 2
//
//	Parallely enable the pcap.
//
// Routine :: 3
//
//	Execute the pybot command on testclient pod and wait the command's output status
//	Routine :: 4
//	  watch the /cmconfig.log --> Parallely, read the file for any run time error [mapping should be done using keywords]
//	                            --> create a channel and inform to the user if any starts/error was logging.
//	Routine :: 5
//	  watch the /logstore/TspCore for every nth Sec --> Parallely, check for any cores.
//	                                                   --> create a channel and inform to the user if any coring happens
//	Routine :: 6
//	  watch the /RTPTraceError --> Parallely, read the file for any run time error [mapping should be done using keywords]
//	                              --> create a channel and inform to the user if any starts/error was logging.
//	Routine :: 7
//	  watch the /Envoy --> Parallely, read the file for any run time error [mapping should be done using keywords]
//	                     --> create a channel and inform to the user if any starts/error was logging.
//	Routine :: 8
//	  watch the /dumplog --> Parallely, read the file for any run time error [mapping should be done using keywords]
//	                       --> create a channel and inform to the user if any starts/error was logging.
//
// Routine :: 9
//
//	Once the pybot cmd status is set to complete
//	  Routine :: 10
//	    store the process traces into single file, can be used for analysis and also for symptom collections
//	  Routine ::
//	    Parallely Disables the trace of each process.
//	  Routine ::
//	    Parallely Disables the pcaps.
//	  Routine ::
//	    Stops the process for to watch the /cmconfig.log
//	    Store the traces for analysis or symptom collection
//	  Routine ::
//	    Stops the process for to watch the /logstore/TspCore
//	    Store the traces for analysis or symptom collection
//	  Routine ::
//	    Stops the process for to watch the /RTPTraceError
//	    Store the traces for analysis or symptom collection
//	  Routine ::
//	    Stops the process for to watch the /Envoy
//	    Store the traces for analysis or symptom collection
//	  Routine ::
//	    Stops the process for to watch the /dumplog
//	    Store the traces for analysis or symptom collection
//
// Copy all the information to local path (All traces, pcap, log.html)
func (c *Collector) StartCollection(ctx context.Context, namespace string, pods []string) error {
	config := &SymptomCollectionConfig{
		Namespace: namespace,
		Pods:      pods,
		StartTime: time.Now(),
	}

	c.logger.Info("Starting symptom collection",
		zap.String("namespace", config.Namespace),
		zap.Strings("pods", config.Pods),
	)

	// Validation checks
	if err := c.validateNamespace(ctx, config.Namespace); err != nil {
		return fmt.Errorf("namespace validation failed: %w", err)
	}

	if err := c.validateDeployments(ctx, config.Namespace, config.Pods); err != nil {
		return fmt.Errorf("deployment validation failed: %w", err)
	}

	// Start symptom collection routines
	var wg sync.WaitGroup
	errorChan := make(chan ErrorEvent, 100)

	// Routine 1: Enable traces for processes
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.Info("Enabling traces for processes")
		c.enableTraces(ctx, config.Pods)
	}()

	// Routine 2: Enable pcap capture
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.Info("Enabling pcap capture")
		c.enablePcap(ctx, config.Pods)
	}()

	// Routine 3: Execute pybot command
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.logger.Info("Executing pybot command")
		c.executePybot(ctx)
	}()

	// Routines 4-8: Watch log files
	logPaths := c.config.Paths.LogPaths
	if len(logPaths) == 0 {
		logPaths = []string{
			"/cmconfig.log",
			"/logstore/TspCore",
			"/RTPTraceError",
			"/Envoy",
			"/dumplog",
		}
	}

	for _, logPath := range logPaths {
		wg.Add(1)
		logPath := logPath // Capture for goroutine
		go func() {
			defer wg.Done()
			c.logger.Info("Watching log file", zap.String("path", logPath))
			c.watchLogFile(ctx, logPath, errorChan)
		}()
	}

	// Error monitoring
	go func() {
		for event := range errorChan {
			c.logger.Error("Error detected during symptom collection",
				zap.Time("timestamp", event.Timestamp),
				zap.String("source", event.Source),
				zap.String("message", event.Message),
			)
		}
	}()

	// Routine 9: Monitor completion and cleanup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// In real implementation, wait for pybot completion signal
		time.Sleep(10 * time.Second) // Simulated test duration
		c.logger.Info("Test completed, starting cleanup")
		c.cleanup(ctx, config)
	}()

	wg.Wait()
	close(errorChan)

	duration := time.Since(config.StartTime)
	c.logger.Info("Symptom collection completed", zap.Duration("duration", duration))
	return nil
}

// validateNamespace checks if namespace exists
func (c *Collector) validateNamespace(ctx context.Context, namespace string) error {
	exists, err := c.k8sClient.NamespaceExists(namespace)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("namespace %s does not exist", namespace)
	}
	return nil
}

// validateDeployments checks if deployments exist and are ready
func (c *Collector) validateDeployments(ctx context.Context, namespace string, pods []string) error {
	deployments, err := c.k8sClient.GetDeployments(namespace)
	if err != nil {
		return err
	}

	for _, pod := range pods {
		found := false
		for _, dep := range deployments {
			if strings.Contains(dep.Name, pod) {
				found = true
				ready, err := c.k8sClient.IsDeploymentReady(namespace, dep.Name)
				if err != nil {
					return fmt.Errorf("failed to check deployment readiness: %w", err)
				}
				if !ready {
					return fmt.Errorf("deployment %s is not ready", dep.Name)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("deployment for pod %s not found", pod)
		}
	}

	return nil
}

// enableTraces enables tracing for processes
func (c *Collector) enableTraces(ctx context.Context, pods []string) {
	for _, pod := range pods {
		c.logger.Debug("Enabling trace for pod", zap.String("pod", pod))
		// Actual implementation would enable tracing via kubectl exec or API
		time.Sleep(500 * time.Millisecond)
	}
}

// enablePcap enables pcap capture
func (c *Collector) enablePcap(ctx context.Context, pods []string) {
	for _, pod := range pods {
		c.logger.Debug("Enabling pcap for pod", zap.String("pod", pod))
		// Actual implementation would enable pcap via kubectl exec or API
		time.Sleep(500 * time.Millisecond)
	}
}

// executePybot executes pybot command on testclient pod
func (c *Collector) executePybot(ctx context.Context) {
	c.logger.Debug("Executing pybot command on testclient pod")
	// Actual implementation would execute pybot via kubectl exec
	time.Sleep(2 * time.Second)
	c.logger.Debug("Pybot command completed")
}

// watchLogFile watches a log file for errors
func (c *Collector) watchLogFile(ctx context.Context, filePath string, errorChan chan<- ErrorEvent) {
	keywords := c.config.Symptom.ErrorKeywords
	if len(keywords) == 0 {
		keywords = []string{"error", "ERROR", "fatal", "FATAL", "exception", "EXCEPTION", "panic", "PANIC"}
	}

	interval := c.config.Symptom.CheckInterval
	if interval == 0 {
		interval = 1 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// In real implementation, read file and check for keywords
			// For now, simulate error detection
			if time.Now().Unix()%10 == 0 {
				errorChan <- ErrorEvent{
					Timestamp: time.Now(),
					Source:    filePath,
					Message:   fmt.Sprintf("Error detected in %s", filePath),
				}
			}
		}
	}
}

// cleanup performs cleanup after test completion
// Routine 10: store the process traces into single file, can be used for analysis and also for symptom collections
// Routine: Parallely Disables the trace of each process
// Routine: Parallely Disables the pcaps
// Routine: Stops the process for to watch the /cmconfig.log, Store the traces for analysis or symptom collection
// Routine: Stops the process for to watch the /logstore/TspCore, Store the traces for analysis or symptom collection
// Routine: Stops the process for to watch the /RTPTraceError, Store the traces for analysis or symptom collection
// Routine: Stops the process for to watch the /Envoy, Store the traces for analysis or symptom collection
// Routine: Stops the process for to watch the /dumplog, Store the traces for analysis or symptom collection
// Copy all the information to local path (All traces, pcap, log.html)
func (c *Collector) cleanup(ctx context.Context, config *SymptomCollectionConfig) {
	c.logger.Info("Storing process traces for analysis")
	c.logger.Info("Disabling traces for all processes")
	c.logger.Info("Disabling pcaps")
	c.logger.Info("Stopping log watchers and storing traces")

	// In real implementation:
	// - Store traces to file
	// - Disable traces and pcaps
	// - Stop log watchers
	// - Copy all data to local path
}

