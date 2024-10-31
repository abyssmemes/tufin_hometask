package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func int32Ptr(i int32) *int32 { return &i }

func main() {
	var rootCmd = &cobra.Command{Use: "tufin"}

	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(statusCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var clusterCmd = &cobra.Command{
	Use:   "cluster [minikube]",
	Short: "Deploy a Kubernetes cluster (k3s by default, or minikube if specified)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
	},
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
	case "windows":
		fmt.Println("Please install minikube manually from https://minikube.sigs.k8s.io/docs/start/")
		return fmt.Errorf("automatic installation not supported on Windows")
	default:
		fmt.Println("Unsupported OS. Please install minikube manually.")
		return fmt.Errorf("unsupported OS")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy MySQL and WordPress pods",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deploying MySQL and WordPress pods...")

		// Build the config from kubeconfig file
		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			fmt.Println("Cannot find kubeconfig")
			os.Exit(1)
		}

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("Failed to build kubeconfig: %v\n", err)
			os.Exit(1)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Printf("Failed to create clientset: %v\n", err)
			os.Exit(1)
		}

		ctx := context.TODO()

		// Deploy MySQL
		fmt.Println("Deploying MySQL...")
		mysqlDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "mysql",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "mysql",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "mysql",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "mysql",
								Image: "mysql:latest",
								Env: []corev1.EnvVar{
									{
										Name:  "MYSQL_ROOT_PASSWORD",
										Value: "password",
									},
								},
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 3306,
									},
								},
							},
						},
					},
				},
			},
		}

		_, err = clientset.AppsV1().Deployments("default").Create(ctx, mysqlDeployment, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Failed to create MySQL deployment: %v\n", err)
		} else {
			fmt.Println("MySQL deployment created.")
		}

		// Create a Service for MySQL
		mysqlService := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "mysql",
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": "mysql",
				},
				Ports: []corev1.ServicePort{
					{
						Port:     3306,
						Protocol: corev1.ProtocolTCP,
					},
				},
			},
		}

		_, err = clientset.CoreV1().Services("default").Create(ctx, mysqlService, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Failed to create MySQL service: %v\n", err)
		} else {
			fmt.Println("MySQL service created.")
		}

		// Deploy WordPress
		fmt.Println("Deploying WordPress...")
		wordpressDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "wordpress",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "wordpress",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "wordpress",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "wordpress",
								Image: "wordpress:4.8-apache",
								Env: []corev1.EnvVar{
									{
										Name:  "WORDPRESS_DB_HOST",
										Value: "mysql",
									},
									{
										Name:  "WORDPRESS_DB_PASSWORD",
										Value: "password",
									},
								},
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 80,
									},
								},
							},
						},
					},
				},
			},
		}

		_, err = clientset.AppsV1().Deployments("default").Create(ctx, wordpressDeployment, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Failed to create WordPress deployment: %v\n", err)
		} else {
			fmt.Println("WordPress deployment created.")
		}

		// Create a Service for WordPress
		wordpressService := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "wordpress",
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": "wordpress",
				},
				Ports: []corev1.ServicePort{
					{
						Port:     80,
						Protocol: corev1.ProtocolTCP,
					},
				},
				Type: corev1.ServiceTypeNodePort,
			},
		}

		_, err = clientset.CoreV1().Services("default").Create(ctx, wordpressService, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Failed to create WordPress service: %v\n", err)
		} else {
			fmt.Println("WordPress service created.")
		}

		fmt.Println("Deployment completed.")
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the status of pods in the default namespace",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting pod status in default namespace...")

		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			fmt.Println("Cannot find kubeconfig")
			os.Exit(1)
		}

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("Failed to build kubeconfig: %v\n", err)
			os.Exit(1)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Printf("Failed to create clientset: %v\n", err)
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
	},
}
