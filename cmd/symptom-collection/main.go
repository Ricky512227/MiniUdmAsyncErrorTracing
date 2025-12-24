package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/internal/logger"
	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/config"
	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/kubernetes"
	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/symptom"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	namespace string
	pods      string
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "symptom-collection",
	Short: "Start symptom collection for Kubernetes pods",
	Long: `Start symptom collection process that:
- Enables traces for processes
- Enables pcap capture
- Executes test commands
- Monitors log files for errors
- Collects and stores traces for analysis`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if pods == "" {
			return fmt.Errorf("pod names (-p) are required")
		}

		// Load configuration
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Initialize logger
		if err := logger.Init(cfg.Logging.Level); err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}
		defer logger.Sync()

		// Use namespace from flag or config
		ns := namespace
		if ns == "" {
			ns = cfg.Kubernetes.Namespace
		}

		// Parse pod names
		podList := strings.Fields(pods)

		logger.Logger.Info("Starting symptom collection",
			zap.String("namespace", ns),
			zap.Strings("pods", podList),
		)

		// Create Kubernetes client
		ctx := context.Background()
		k8sClient, err := kubernetes.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}

		// Create collector
		collector := symptom.NewCollector(cfg, k8sClient, logger.Logger)

		// Start collection
		if err := collector.StartCollection(ctx, ns, podList); err != nil {
			return fmt.Errorf("symptom collection failed: %w", err)
		}

		logger.Logger.Info("Symptom collection completed successfully")
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (default: from config)")
	rootCmd.Flags().StringVarP(&pods, "pods", "p", "", "Space-separated pod names (required)")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	
	rootCmd.MarkFlagRequired("pods")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Logger.Error("Command failed", zap.Error(err))
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

