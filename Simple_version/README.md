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
   cd tufin_hometask/Simple_version
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
   go build -o tufin main.go
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

Structure

The tufin client is written in Go and leverages the following key libraries:

	•	Cobra: For building command-line interfaces (CLI). It simplifies the creation and organization of commands and subcommands.
	•	client-go: The official Kubernetes client library for Go. It enables interaction with Kubernetes clusters through the Kubernetes API.

The main components of the code are:

	•	main.go: The primary Go file containing the implementation of the commands and the core logic of the client.

Main Components

1. Import Statements

At the beginning of the main.go file, necessary packages are imported, including standard Go packages and external libraries:
```bash
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
```
2. Helper Functions
```bash
int32Ptr
```
A helper function to get a pointer to an int32 value, which is required by certain Kubernetes API fields.
```bash
func int32Ptr(i int32) *int32 { return &i }
```
3. Main Function

The main function initializes the root command and adds subcommands for cluster management, deployment, and status checking.
```bash
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
```
Commands

1. clusterCmd Command

Purpose

Deploys a local Kubernetes cluster using either k3s (default) or minikube if specified.

Usage

tufin cluster [minikube]

Implementation Details

	•	Argument Parsing: Checks if the optional argument minikube is provided.
	•	Cluster Deployment:
	•	If minikube is specified, calls deployMinikubeCluster().
	•	Otherwise, defaults to deploying a k3s cluster via deployK3sCluster().

Functions

deployK3sCluster

	•	Checks for k3s Installation:
	•	Uses exec.LookPath("k3s") to check if k3s is installed.
	•	If not installed, calls installK3s().
	•	Installs k3s:
	•	The installK3s function runs the k3s installation script for Linux.
	•	Starts k3s Server:
	•	Runs sudo k3s server with flags to write the kubeconfig file to ~/.kube/config.
	•	Uses time.Sleep to wait for the server to be ready.

deployMinikubeCluster

	•	Checks for minikube Installation:
	•	Uses exec.LookPath("minikube") to check if minikube is installed.
	•	If not installed, calls installMinikube().
	•	Installs minikube:
	•	The installMinikube function handles installation for different operating systems:
	•	Linux: Downloads the minikube binary and installs it to /usr/local/bin/.
	•	macOS: Installs minikube via Homebrew.
	•	Starts minikube:
	•	Runs minikube start to initialize the cluster.

installK3s

	•	Functionality:
	•	Downloads and runs the official k3s installation script for Linux.
	•	Operating System Support:
	•	Only supports Linux.
	•	For other operating systems, prompts the user to install k3s manually.

installMinikube

	•	Functionality:
	•	Installs minikube based on the operating system.
	•	Operating System Support:
	•	Linux: Downloads the minikube binary.
	•	macOS: Uses Homebrew to install minikube.
	•	Windows: Prompts the user to install minikube manually.

2. deployCmd Command

Purpose

Deploys MySQL and WordPress pods into the cluster, setting up WordPress to connect to MySQL.

Usage

tufin deploy

Implementation Details

	•	Kubeconfig Handling:
	•	Determines the path to the kubeconfig file, defaulting to ~/.kube/config.
	•	Clientset Creation:
	•	Builds a Kubernetes client configuration and creates a clientset for interacting with the cluster.
	•	Deployment of Resources:
	•	MySQL Deployment:
	•	Creates a Deployment for MySQL with one replica.
	•	Uses the image mysql:5.6.
	•	Sets the MYSQL_ROOT_PASSWORD environment variable.
	•	Exposes port 3306.
	•	MySQL Service:
	•	Creates a Service to expose the MySQL deployment internally.
	•	WordPress Deployment:
	•	Creates a Deployment for WordPress with one replica.
	•	Uses the image wordpress:4.8-apache.
	•	Sets environment variables to connect to MySQL (WORDPRESS_DB_HOST, WORDPRESS_DB_PASSWORD).
	•	Exposes port 80.
	•	WordPress Service:
	•	Creates a Service of type NodePort to expose WordPress externally.
	•	Error Handling:
	•	Checks for errors during resource creation and prints informative messages.

3. statusCmd Command

Purpose

Prints out a status table containing the pod names and their status in the default namespace.

Usage

tufin status

Implementation Details

	•	Kubeconfig Handling:
	•	Same as in deployCmd.
	•	Clientset Creation:
	•	Same as in deployCmd.
	•	Listing Pods:
	•	Retrieves a list of pods in the default namespace.
	•	Iterates over the pods and prints their names and status phases in a formatted table.

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

