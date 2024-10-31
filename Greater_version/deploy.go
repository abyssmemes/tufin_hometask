// deploy.go

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy applications to the Kubernetes cluster",
	Run:   runDeployCmd,
}

func runDeployCmd(cmd *cobra.Command, args []string) {
	fmt.Println("Deploying applications...")

	clientset, config, err := getKubernetesClient()
	if err != nil {
		fmt.Printf("Failed to create Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// Deploy MySQL
	fmt.Println("Deploying MySQL...")
	if err := deployMySQL(clientset, config); err != nil {
		fmt.Printf("Failed to deploy MySQL: %v\n", err)
	} else {
		fmt.Println("MySQL deployed successfully.")
	}

	// Deploy WordPress
	fmt.Println("Deploying WordPress...")
	if err := deployWordPress(clientset, config); err != nil {
		fmt.Printf("Failed to deploy WordPress: %v\n", err)
	} else {
		fmt.Println("WordPress deployed successfully.")
	}

	fmt.Println("Deployment completed.")
}

func deployMySQL(clientset *kubernetes.Clientset, config *rest.Config) error {
	deploymentPath := filepath.Join("deployments", "mysql-deployment.yaml")
	servicePath := filepath.Join("deployments", "mysql-service.yaml")

	if err := applyYAML(config, deploymentPath); err != nil {
		return err
	}
	if err := applyYAML(config, servicePath); err != nil {
		return err
	}
	return nil
}

func deployWordPress(clientset *kubernetes.Clientset, config *rest.Config) error {
	deploymentPath := filepath.Join("deployments", "wordpress-deployment.yaml")
	servicePath := filepath.Join("deployments", "wordpress-service.yaml")

	if err := applyYAML(config, deploymentPath); err != nil {
		return err
	}
	if err := applyYAML(config, servicePath); err != nil {
		return err
	}
	return nil
}
