package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/internal/logger"
	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/config"
	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/kubernetes"
	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	namespace string
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "list-deployments",
	Short: "List Kubernetes deployments in a namespace",
	Long:  `List all deployments in the specified Kubernetes namespace with their status and details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Initialize logger
		if err := logger.InitDevelopment(cfg.Logging.Level); err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}
		defer logger.Sync()

		// Use namespace from flag or config
		ns := namespace
		if ns == "" {
			ns = cfg.Kubernetes.Namespace
		}

		logger.Logger.Info("Listing deployments", zap.String("namespace", ns))

		// Create Kubernetes client
		ctx := context.Background()
		k8sClient, err := kubernetes.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}

		// Get deployments
		deployments, err := k8sClient.GetDeployments(ns)
		if err != nil {
			return fmt.Errorf("failed to get deployments: %w", err)
		}

		if len(deployments) == 0 {
			fmt.Printf("No deployments found in namespace '%s'\n", ns)
			return nil
		}

		fmt.Printf("\nFound %d deployment(s) in namespace '%s':\n\n", len(deployments), ns)
		for i, deployment := range deployments {
			fmt.Printf("%d. Name: %s\n", i+1, deployment.Name)
			
			replicas := int32(0)
			if deployment.Spec.Replicas != nil {
				replicas = *deployment.Spec.Replicas
			}
			
			fmt.Printf("   Replicas: %d/%d\n", deployment.Status.ReadyReplicas, replicas)
			fmt.Printf("   Age: %s\n", utils.GetAge(deployment.CreationTimestamp.Time))
			
			if len(deployment.Spec.Template.Spec.Containers) > 0 {
				fmt.Printf("   Image: %s\n", deployment.Spec.Template.Spec.Containers[0].Image)
			}
			
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (default: from config)")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

