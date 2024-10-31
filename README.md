# Tufin Go Client

Tufin is a Go-based command-line client designed to simplify the deployment and management of local Kubernetes clusters and applications. It provides functionalities for:

	•	Cluster Deployment: Deploy a local Kubernetes cluster using either k3s or minikube.
	•	Application Deployment: Deploy MySQL and WordPress applications within the cluster.
	•	Status Monitoring: Check the status of deployed pods in the default namespace.

## Project Structure

The project is organized into two main directories, each representing a different version of the client:

1. ### Simple Version

Directory: Simple version

	•	Description: This version contains the basic implementation of the client with all functionalities defined within a single main.go file.
	•	Characteristics:
	•	All commands and functions are located in one file.
	•	Easier to read for small-scale applications.
	•	Suitable for understanding the basic flow and logic of the application.
	•	Usage:
	•	Provides the core functionalities of cluster deployment, application deployment, and status checking.

2. ### Greater Version

Directory: Greater version

	•	Description: This is the enhanced and refactored version of the client, featuring improved code organization and maintainability through separation of concerns.
	•	Characteristics:
	•	Functions are separated into multiple files based on their purpose:
	•	main.go: Entry point of the application.
	•	cluster.go: Functions related to cluster deployment.
	•	deploy.go: Functions for deploying applications.
	•	status.go: Functions for checking the status of deployments.
	•	kubeclient.go: Kubernetes client initialization.
	•	applyyaml.go: Functions to apply YAML configurations.
	•	Deployment configurations are separated into YAML files located in the deployments/ directory.
	•	Improved code readability and maintainability.
	•	Follows Go best practices and design patterns.
	•	Usage:
	•	Provides the same core functionalities as the simple version.
	•	Easier to extend and maintain due to modular code structure.