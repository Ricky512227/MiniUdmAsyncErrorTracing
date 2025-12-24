package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Client wraps Kubernetes clientset and provides high-level operations
type Client struct {
	Clientset *kubernetes.Clientset
	Context   context.Context
}

// NewClient creates a new Kubernetes client
func NewClient(ctx context.Context) (*Client, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{
		Clientset: clientset,
		Context:   ctx,
	}, nil
}

// GetDeployments retrieves all deployments in a namespace
func (c *Client) GetDeployments(namespace string) ([]v1.Deployment, error) {
	list, err := c.Clientset.AppsV1().Deployments(namespace).
		List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	return list.Items, nil
}

// GetDeployment retrieves a specific deployment by name
func (c *Client) GetDeployment(namespace, name string) (*v1.Deployment, error) {
	deployment, err := c.Clientset.AppsV1().Deployments(namespace).
		Get(c.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s: %w", name, err)
	}
	return deployment, nil
}

// GetNamespace retrieves a namespace by name
func (c *Client) GetNamespace(name string) (*corev1.Namespace, error) {
	namespace, err := c.Clientset.CoreV1().Namespaces().
		Get(c.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace %s: %w", name, err)
	}
	return namespace, nil
}

// NamespaceExists checks if a namespace exists
func (c *Client) NamespaceExists(name string) (bool, error) {
	_, err := c.GetNamespace(name)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// GetPods retrieves all pods in a namespace
func (c *Client) GetPods(namespace string) ([]corev1.Pod, error) {
	list, err := c.Clientset.CoreV1().Pods(namespace).
		List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	return list.Items, nil
}

// GetPod retrieves a specific pod by name
func (c *Client) GetPod(namespace, name string) (*corev1.Pod, error) {
	pod, err := c.Clientset.CoreV1().Pods(namespace).
		Get(c.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s: %w", name, err)
	}
	return pod, nil
}

// IsDeploymentReady checks if a deployment is ready
func (c *Client) IsDeploymentReady(namespace, name string) (bool, error) {
	deployment, err := c.GetDeployment(namespace, name)
	if err != nil {
		return false, err
	}

	if deployment.Status.ReadyReplicas < *deployment.Spec.Replicas {
		return false, nil
	}

	return true, nil
}

