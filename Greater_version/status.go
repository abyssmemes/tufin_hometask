package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the status of pods in the default namespace",
	Run:   runStatusCmd,
}

func runStatusCmd(cmd *cobra.Command, args []string) {
	fmt.Println("Getting pod status in default namespace...")

	clientset, _, err := getKubernetesClient()
	if err != nil {
		fmt.Printf("Failed to create Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	ctx := context.TODO()

	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Failed to list pods: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%-30s %-20s\n", "Pod Name", "Status")
	for _, pod := range pods.Items {
		fmt.Printf("%-30s %-20s\n", pod.Name, string(pod.Status.Phase))
	}
}
