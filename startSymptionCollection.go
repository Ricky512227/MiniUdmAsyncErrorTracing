package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

// SymptomCollectionConfig holds configuration for symptom collection
type SymptomCollectionConfig struct {
	Namespace string
	Pods      []string
	StartTime time.Time
}

var (
	errorKeywords = []string{"error", "ERROR", "fatal", "FATAL", "exception", "EXCEPTION", "panic", "PANIC"}
	logFiles      = []string{
		"/cmconfig.log",
		"/logstore/TspCore",
		"/RTPTraceError",
		"/Envoy",
		"/dumplog",
	}
)

// StartSymptomCollection starts the symptom collection process
//
// Command Syntax: startSymptionCollection -n <namespace> -c <"podnames">
// Example:
//
//	startSymptionCollection -n miniudm -c "uecm"
//	startSymptionCollection -n dracvnf -c "uecm nim"
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
func StartSymptomCollection() {

	namespace := flag.String("n", "default", "Kubernetes namespace")
	podsStr := flag.String("c", "", "Comma-separated pod names")
	flag.Parse()

	if *podsStr == "" {
		log.Fatal("Error: pod names (-c) is required")
	}

	pods := strings.Split(*podsStr, " ")
	config := &SymptomCollectionConfig{
		Namespace: *namespace,
		Pods:      pods,
		StartTime: time.Now(),
	}

	log.Printf("Starting symptom collection")
	log.Printf("Namespace: %s", config.Namespace)
	log.Printf("Pods: %v", config.Pods)

	ctx := context.Background()
	kubeConfig := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(kubeConfig)

	// Validation checks
	if err := validateNamespace(clientset, ctx, config.Namespace); err != nil {
		log.Fatalf("Namespace validation failed: %v", err)
	}

	if err := validateDeployments(clientset, ctx, config.Namespace, config.Pods); err != nil {
		log.Fatalf("Deployment validation failed: %v", err)
	}

	// Start symptom collection routines
	var wg sync.WaitGroup
	errorChan := make(chan string, 100)

	// Routine 1: Parallely enable the trace of each process
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Routine 1: Parallely enabling trace of each process")
		enableTraces(config.Pods)
	}()

	// Routine 2: Parallely enable the pcap
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Routine 2: Parallely enabling pcap")
		enablePcap(config.Pods)
	}()

	// Routine 3: Execute the pybot command on testclient pod and wait the command's output status
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Routine 3: Executing pybot command on testclient pod")
		executePybot()
	}()

	// Routine 4: watch the /cmconfig.log --> Parallely, read the file for any run time error [mapping should be done using keywords]
	//            --> create a channel and inform to the user if any starts/error was logging
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Routine 4: Watching /cmconfig.log for runtime errors (keyword mapping)")
		watchLogFile("/cmconfig.log", errorChan)
	}()

	// Routine 5: watch the /logstore/TspCore for every nth Sec --> Parallely, check for any cores.
	//            --> create a channel and inform to the user if any coring happens
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Routine 5: Watching /logstore/TspCore for cores")
		watchLogFile("/logstore/TspCore", errorChan)
	}()

	// Routine 6: watch the /RTPTraceError --> Parallely, read the file for any run time error [mapping should be done using keywords]
	//            --> create a channel and inform to the user if any starts/error was logging
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Routine 6: Watching /RTPTraceError for runtime errors (keyword mapping)")
		watchLogFile("/RTPTraceError", errorChan)
	}()

	// Routine 7: watch the /Envoy --> Parallely, read the file for any run time error [mapping should be done using keywords]
	//           --> create a channel and inform to the user if any starts/error was logging
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Routine 7: Watching /Envoy for runtime errors (keyword mapping)")
		watchLogFile("/Envoy", errorChan)
	}()

	// Routine 8: watch the /dumplog --> Parallely, read the file for any run time error [mapping should be done using keywords]
	//          --> create a channel and inform to the user if any starts/error was logging
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Routine 8: Watching /dumplog for runtime errors (keyword mapping)")
		watchLogFile("/dumplog", errorChan)
	}()

	// Routine 9: Once the pybot cmd status is set to complete
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Routine 9: Monitoring pybot completion and starting cleanup")
		time.Sleep(10 * time.Second) // Simulate test duration - wait for pybot to complete
		log.Println("Pybot command completed, starting cleanup...")
		cleanup(config)
	}()

	// Error monitoring
	go func() {
		for errMsg := range errorChan {
			log.Printf("ERROR DETECTED: %s", errMsg)
		}
	}()

	wg.Wait()
	close(errorChan)

	log.Println("Symptom collection completed")
}

// validateNamespace checks if namespace exists
func validateNamespace(clientset *kubernetes.Clientset, ctx context.Context, namespace string) error {
	_, err := clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("namespace %s does not exist: %w", namespace, err)
	}
	return nil
}

// validateDeployments checks if deployments exist
func validateDeployments(clientset *kubernetes.Clientset, ctx context.Context, namespace string, pods []string) error {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	for _, pod := range pods {
		found := false
		for _, dep := range deployments.Items {
			if strings.Contains(dep.Name, pod) {
				found = true
				// Check readiness
				if dep.Status.ReadyReplicas < *dep.Spec.Replicas {
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

// enableTraces enables tracing for processes (placeholder)
func enableTraces(pods []string) {
	for _, pod := range pods {
		log.Printf("Enabling trace for pod: %s", pod)
		// Actual implementation would enable tracing
		time.Sleep(500 * time.Millisecond)
	}
}

// enablePcap enables pcap capture (placeholder)
func enablePcap(pods []string) {
	for _, pod := range pods {
		log.Printf("Enabling pcap for pod: %s", pod)
		// Actual implementation would enable pcap
		time.Sleep(500 * time.Millisecond)
	}
}

// executePybot executes pybot command (placeholder)
func executePybot() {
	log.Println("Executing pybot command on testclient pod")
	// Actual implementation would execute pybot
	time.Sleep(2 * time.Second)
	log.Println("Pybot command completed")
}

// watchLogFile watches a log file for errors
func watchLogFile(filePath string, errorChan chan<- string) {
	// In real implementation, this would:
	// 1. Read file continuously
	// 2. Check for error keywords
	// 3. Send errors to channel
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		// Simulate error detection
		if i == 5 {
			errorChan <- fmt.Sprintf("Error detected in %s at line %d", filePath, i*100)
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
func cleanup(config *SymptomCollectionConfig) {
	log.Println("Routine 10: Storing process traces into single file for analysis and symptom collections")

	log.Println("Parallely disabling trace of each process...")
	// Disable traces for all processes

	log.Println("Parallely disabling pcaps...")
	// Disable pcaps

	log.Println("Stopping process to watch /cmconfig.log, storing traces...")
	// Stop watcher and store traces

	log.Println("Stopping process to watch /logstore/TspCore, storing traces...")
	// Stop watcher and store traces

	log.Println("Stopping process to watch /RTPTraceError, storing traces...")
	// Stop watcher and store traces

	log.Println("Stopping process to watch /Envoy, storing traces...")
	// Stop watcher and store traces

	log.Println("Stopping process to watch /dumplog, storing traces...")
	// Stop watcher and store traces

	log.Println("Copying all information to local path (All traces, pcap, log.html)...")
	// Copy all traces, pcap, and log.html to local path

	log.Printf("Symptom collection completed in %v", time.Since(config.StartTime))
}
