# Tufin Go Client

**Tufin** is a Go-based command-line client designed to simplify the deployment and management of a local Kubernetes cluster and applications. It provides the following functionalities:

- **Cluster Deployment**: Deploy a local Kubernetes cluster using either **k3s** or **minikube**.
- **Application Deployment**: Deploy MySQL and WordPress pods connected together within the cluster.
- **Status Monitoring**: Check the status of the deployed pods in the default namespace.

**Note**: Adding **minikube** support was essential for testing on macOS, as installing and running **k3s** directly on macOS is not straightforward and often requires additional tools like virtual machines or Docker Desktop.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Commands Overview](#commands-overview)
    - [`tufin cluster`](#tufin-cluster)
    - [`tufin deploy`](#tufin-deploy)
    - [`tufin status`](#tufin-status)
- [Usage Instructions](#usage-instructions)
    - [Deploying a Kubernetes Cluster](#deploying-a-kubernetes-cluster)
        - [Option 1: Deploying a k3s Cluster (Default)](#option-1-deploying-a-k3s-cluster-default)
        - [Option 2: Deploying a Minikube Cluster](#option-2-deploying-a-minikube-cluster)
    - [Deploying Applications](#deploying-applications)
    - [Checking Pod Status](#checking-pod-status)
- [Code Structure and Explanation](#code-structure-and-explanation)
- [References](#references)


---

## Prerequisites

Before using the **tufin** client, ensure you have the following installed on your system:

- **Go Programming Language**: [Install Go](https://golang.org/dl/)
- **Git**: [Install Git](https://git-scm.com/downloads)
- **Docker**: Required for pulling container images and running containers.
- **kubectl** (optional): For verifying cluster status and debugging.
- **k3s** Cluster tool.
- **minikube** Alternative Cluster tool.
---

## Installation

1. **Clone the Repository**:

```bash
   git clone https://github.com/abyssmemes/tufin_hometask.git
   cd tufin_hometask/Greater_version
```

2.	**Install Dependencies**:

Ensure you have the necessary Go packages:

```bash
    go get github.com/spf13/cobra@v1.6.1
    go get k8s.io/client-go@v0.27.4
```

3.	**Build the Executable**:

Compile the code to create the tufin executable:
```bash
   go build -o tufin main.go cluster.go deploy.go status.go kubeclient.go applyyaml.go
```


## Commands Overview

The tufin client provides three main commands:
```bash
tufin cluster.go
```
Usage: tufin cluster [minikube]

**Description**: Deploys a local Kubernetes cluster.

**Options**:

No arguments: Deploys a k3s cluster (default).

minikube: Deploys a minikube cluster.


```bash
tufin deploy
```
Usage: tufin deploy

**Description**: 
    Deploys MySQL and WordPress pods into the cluster, connecting WordPress to MySQL.
```bash
tufin status
```
Usage: tufin status

**Description**: Prints a status table containing the pod names and their current status in the default namespace.



## Usage Instructions

Deploying a Kubernetes Cluster

Option 1: Deploying a k3s Cluster (Default)

./tufin cluster

	•	Notes:
	•	The script will check if k3s is installed.
	•	If not installed, it will attempt to install k3s automatically (Linux only).
	•	Requires root or sudo privileges for installation and starting the cluster.
	•	Not recommended for macOS due to installation complexities.

Option 2: Deploying a Minikube Cluster

./tufin cluster minikube

	•	Notes:
	•	The script will check if minikube is installed.
	•	If not installed, it will attempt to install minikube automatically.
	•	Linux: Downloads and installs the minikube binary.
	•	macOS: Installs via Homebrew.
	•	Windows: Prompts the user to install manually.
	•	Suitable for macOS users, as minikube provides a straightforward way to run Kubernetes locally on macOS.
	•	May require administrative privileges during installation.

Deploying Applications

After the cluster is up and running, deploy the applications:

./tufin deploy

	•	Description:
	•	Deploys a MySQL deployment and service.
	•	Deploys a WordPress deployment and service, connected to MySQL.
	•	Prerequisites:
	•	Ensure the kubeconfig file is available at ~/.kube/config or adjust the path in the code.
	•	The cluster should be running (either k3s or minikube).

Checking Pod Status

To check the status of the deployed pods:

./tufin status

	Description:
	•	Lists all pods in the default namespace along with their current status (e.g., Running, Pending).
Example Output:
```bash
Getting pod status in default namespace...
Pod Name                       Status
mysql-xxxxxxxxxx-xxxxx         Running
wordpress-xxxxxxxxxx-xxxxx     Running
```




## Code Structure and Explanation

### Overview

Code Overview

The tufin client is written in Go and leverages the following key libraries:

	•	Cobra: For building command-line interfaces (CLI). It simplifies the creation and organization of commands and subcommands.
	•	client-go: The official Kubernetes client library for Go. It enables interaction with Kubernetes clusters through the Kubernetes API.
	•	apimachinery: Provides tools and interfaces for working with Kubernetes API machinery.

Project Structure

The project is organized into multiple files and directories for better maintainability:

	•	main.go: Entry point of the application. Initializes the root command and adds subcommands.
	•	cluster.go: Contains functions related to cluster deployment (deployK3sCluster, deployMinikubeCluster, installK3s, installMinikube).
	•	deploy.go: Contains functions for deploying applications (deployMySQL, deployWordPress).
	•	status.go: Contains functions for checking the status of deployments.
	•	kubeclient.go: Contains functions for Kubernetes client initialization (getKubernetesClient).
	•	applyyaml.go: Contains the applyYAML function for applying YAML configurations to the cluster.
	•	deployments/: Directory containing YAML configuration files for deployments:
	•	mysql-deployment.yaml
	•	mysql-service.yaml
	•	wordpress-deployment.yaml
	•	wordpress-service.yaml

Key Components

1. Cluster Deployment (cluster.go)
```bash
   •	Purpose: Deploys a local Kubernetes cluster using either k3s (default) or minikube if specified.
   •	Functions:
   •	deployK3sCluster: Checks for k3s installation and starts the k3s server.
   •	deployMinikubeCluster: Checks for minikube installation and starts minikube.
   •	installK3s: Installs k3s on Linux systems.
   •	installMinikube: Installs minikube based on the operating system.
```
2. Application Deployment (deploy.go)
```bash
   •	Purpose: Deploys MySQL and WordPress applications into the cluster using YAML configuration files.
   •	Functions:
   •	deployMySQL: Applies the MySQL deployment and service YAML files to the cluster.
   •	deployWordPress: Applies the WordPress deployment and service YAML files to the cluster.
   •	applyYAML: Reads a YAML file and applies it to the cluster using the dynamic Kubernetes client.
```
3. Kubernetes Client Initialization (kubeclient.go)
```bash
   •	Function: getKubernetesClient initializes and returns a Kubernetes clientset for interacting with the cluster.
```
4. Status Monitoring (status.go)
```bash
   •	Purpose: Prints out a status table containing the pod names and their status in the default namespace.
   •	Function: Retrieves a list of pods and displays their names and current status.
```
Deployment YAML Files (deployments/)

	•	Purpose: Defines the Kubernetes resources (Deployments and Services) for MySQL and WordPress.
	•	Files:
	•	mysql-deployment.yaml
	•	mysql-service.yaml
	•	wordpress-deployment.yaml
	•	wordpress-service.yaml

Applying YAML Configurations (applyyaml.go)

	•	Functionality:
	•	The applyYAML function reads a YAML file, decodes it into an unstructured object, and applies it to the cluster.
	•	Uses the dynamic Kubernetes client to handle resources generically.
	•	Checks if the resource already exists and updates it if necessary.

Additional Details

Handling Image Pull Issues

	•	ImagePullPolicy:
	•	By default, the image pull policy is IfNotPresent.
	•	To ensure the latest image is pulled, you can set ImagePullPolicy: corev1.PullAlways in the container specification.
	•	Image Pull Secrets:
	•	If encountering Docker Hub rate limits or authentication issues, you can create an image pull secret and attach it to the pod specification.

Adding Minikube Support

	•	Rationale:
	•	Installing and running k3s directly on macOS is complex due to the lack of native support and the need for virtualization.
	•	Minikube provides an easier alternative for macOS users, enabling them to run Kubernetes locally without extensive setup.
	•	Adjustments Made:
	•	The cluster command accepts an optional [minikube] argument to specify the use of minikube.
	•	Added functions installMinikube and deployMinikubeCluster to handle minikube installation and cluster setup.
	•	The code ensures compatibility with both k3s and minikube clusters by adjusting paths and commands accordingly.

Error Handling and Logging

	•	Command Existence Checks:
	•	Uses exec.LookPath to verify if required binaries (k3s, minikube) are installed before attempting to use them.
	•	Context Handling:
	•	Uses context.TODO() as a placeholder context for Kubernetes API calls.
	•	Exiting on Errors:
	•	The program exits with a non-zero status code if critical errors occur, ensuring that failures are communicated effectively.

