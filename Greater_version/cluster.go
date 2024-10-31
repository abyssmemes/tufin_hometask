package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster [minikube]",
	Short: "Deploy a Kubernetes cluster (k3s by default, or minikube if specified)",
	Args:  cobra.MaximumNArgs(1),
	Run:   runClusterCmd,
}

func runClusterCmd(cmd *cobra.Command, args []string) {
	var clusterType string
	if len(args) > 0 && strings.ToLower(args[0]) == "minikube" {
		clusterType = "minikube"
	} else {
		clusterType = "k3s"
	}

	fmt.Printf("Deploying %s Kubernetes cluster...\n", clusterType)

	if clusterType == "minikube" {
		deployMinikubeCluster()
	} else {
		deployK3sCluster()
	}
}

func deployK3sCluster() {
	// Check if k3s is installed
	if _, err := exec.LookPath("k3s"); err != nil {
		fmt.Println("k3s is not installed. Attempting to install k3s...")
		// Install k3s
		if err := installK3s(); err != nil {
			fmt.Printf("Failed to install k3s: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("k3s installed successfully.")
	}

	// Start the k3s server
	fmt.Println("Starting k3s server...")
	startCmd := exec.Command("sudo", "k3s", "server", "--write-kubeconfig", "~/.kube/config", "--write-kubeconfig-mode", "644")
	startCmd.Stdout = os.Stdout
	startCmd.Stderr = os.Stderr
	if err := startCmd.Start(); err != nil {
		fmt.Printf("Failed to start k3s server: %v\n", err)
		os.Exit(1)
	}

	// Wait a few seconds to ensure k3s is up and running
	fmt.Println("Waiting for k3s to be ready...")
	time.Sleep(15 * time.Second)

	fmt.Println("k3s cluster started successfully.")
}

func deployMinikubeCluster() {
	// Check if minikube is installed
	if _, err := exec.LookPath("minikube"); err != nil {
		fmt.Println("minikube is not installed. Attempting to install minikube...")
		// Install minikube
		if err := installMinikube(); err != nil {
			fmt.Printf("Failed to install minikube: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("minikube installed successfully.")
	}

	// Start minikube
	fmt.Println("Starting minikube...")
	startCmd := exec.Command("minikube", "start")
	startCmd.Stdout = os.Stdout
	startCmd.Stderr = os.Stderr
	if err := startCmd.Run(); err != nil {
		fmt.Printf("Failed to start minikube: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("minikube cluster started successfully.")
}

// Function to install k3s
func installK3s() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("sh", "-c", "curl -sfL https://get.k3s.io | sh -")
	default:
		fmt.Println("Automatic installation of k3s is only supported on Linux.")
		return fmt.Errorf("unsupported OS")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Function to install minikube
func installMinikube() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("sh", "-c", `
            curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 &&
            sudo install minikube /usr/local/bin/
        `)
	case "darwin":
		cmd = exec.Command("brew", "install", "minikube")
	default:
		fmt.Println("Please install minikube manually from https://minikube.sigs.k8s.io/docs/start/")
		return fmt.Errorf("automatic installation not supported on this OS")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
