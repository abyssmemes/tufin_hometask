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

	// List of applications to deploy
	applications := []string{"mysql", "wordpress"}

	// Loop over applications and deploy each one
	for _, app := range applications {
		fmt.Printf("Deploying %s...\n", app)
		if err := deploy(clientset, config, app); err != nil {
			fmt.Printf("Failed to deploy %s: %v\n", app, err)
		} else {
			fmt.Printf("%s deployed successfully.\n", app)
		}
	}

	fmt.Println("Deployment completed.")
}

func deploy(clientset *kubernetes.Clientset, config *rest.Config, appName string) error {
	deploymentPath := filepath.Join("deployments", fmt.Sprintf("%s-deployment.yaml", appName))
	servicePath := filepath.Join("deployments", fmt.Sprintf("%s-service.yaml", appName))

	if err := applyYAML(config, deploymentPath); err != nil {
		return fmt.Errorf("error applying deployment YAML for %s: %v", appName, err)
	}
	if err := applyYAML(config, servicePath); err != nil {
		return fmt.Errorf("error applying service YAML for %s: %v", appName, err)
	}
	return nil
}
