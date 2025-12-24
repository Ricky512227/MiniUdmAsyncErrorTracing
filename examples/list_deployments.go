package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/kubernetes"
)

// Example demonstrates how to list deployments programmatically
func main() {
	ctx := context.Background()

	// Create Kubernetes client
	client, err := kubernetes.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Get deployments from default namespace
	deployments, err := client.GetDeployments("default")
	if err != nil {
		log.Fatalf("Failed to get deployments: %v", err)
	}

	// Print deployment information
	fmt.Printf("Found %d deployment(s):\n\n", len(deployments))
	for i, deployment := range deployments {
		fmt.Printf("%d. %s\n", i+1, deployment.Name)
		fmt.Printf("   Namespace: %s\n", deployment.Namespace)
		fmt.Printf("   Replicas: %d/%d\n",
			deployment.Status.ReadyReplicas,
			*deployment.Spec.Replicas,
		)
		fmt.Println()
	}
}

