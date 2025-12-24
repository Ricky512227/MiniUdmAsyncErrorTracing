package main

import (
	"context"
	"fmt"
	"log"
	"os"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	ctx := context.Background()

	// Get namespace from args or use default
	namespace := "default"
	if len(os.Args) > 1 {
		namespace = os.Args[1]
	}

	config, err := ctrl.GetConfig()
	if err != nil {
		log.Fatalf("Failed to get kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	items, err := GetDeployments(clientset, ctx, namespace)
	if err != nil {
		log.Printf("Error getting deployments: %v", err)
		os.Exit(1)
	}

	if len(items) == 0 {
		fmt.Printf("No deployments found in namespace '%s'\n", namespace)
		return
	}

	fmt.Printf("Found %d deployment(s) in namespace '%s':\n\n", len(items), namespace)
	for i, item := range items {
		fmt.Printf("%d. Name: %s\n", i+1, item.Name)
		fmt.Printf("   Replicas: %d/%d\n", item.Status.ReadyReplicas, *item.Spec.Replicas)
		fmt.Printf("   Age: %s\n", getAge(item.CreationTimestamp.Time))
		if len(item.Spec.Template.Spec.Containers) > 0 {
			fmt.Printf("   Image: %s\n", item.Spec.Template.Spec.Containers[0].Image)
		}
		fmt.Println()
	}
}

// GetDeployments retrieves all deployments in a namespace
func GetDeployments(clientset *kubernetes.Clientset, ctx context.Context,
	namespace string) ([]v1.Deployment, error) {

	list, err := clientset.AppsV1().Deployments(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf("Error listing deployments: %v", err)
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	return list.Items, nil
}

// getAge returns a human-readable age string
func getAge(creationTime interface{}) string {
	// Simple placeholder - in real implementation, calculate time difference
	return "N/A"
}
