Tufin Go Client - README.md

Introduction

Tufin is a Go-based command-line client designed to simplify the deployment and management of a local Kubernetes cluster and applications. It provides the following functionalities:

	•	Cluster Deployment: Deploy a local Kubernetes cluster using either k3s or minikube.
	•	Application Deployment: Deploy MySQL and WordPress pods connected together within the cluster.
	•	Status Monitoring: Check the status of the deployed pods in the default namespace.

Note: Adding minikube support was essential for testing on macOS, as installing and running k3s directly on macOS is not straightforward and often requires additional tools like virtual machines or Docker Desktop.

Table of Contents

	•	Prerequisites
	•	Installation
	•	Commands Overview
	•	tufin cluster
	•	tufin deploy
	•	tufin status
	•	Usage Instructions
	•	Deploying a Kubernetes Cluster
	•	Deploying Applications
	•	Checking Pod Status
	•	Accessing WordPress Application
	•	Code Structure and Explanation
	•	Troubleshooting
	•	Conclusion
	•	References

Prerequisites

Before using the tufin client, ensure you have the following installed on your system:

	•	Go Programming Language: Install Go
	•	Git: Install Git
	•	kubectl (optional): For verifying cluster status and debugging.
	•	Docker: Required for pulling container images and running containers.

Installation

	1.	Clone the Repository:

git clone https://github.com/yourusername/tufin.git
cd tufin


	2.	Install Dependencies:
Ensure you have the necessary Go packages:

go get github.com/spf13/cobra@v1.6.1
go get k8s.io/client-go@v0.27.4


	3.	Build the Executable:
Compile the code to create the tufin executable:

go build -o tufin main.go



Commands Overview

The tufin client provides three main commands:

tufin cluster

Usage: tufin cluster [minikube]

	•	Description: Deploys a local Kubernetes cluster.
	•	Options:
	•	No arguments: Deploys a k3s cluster (default).
	•	minikube: Deploys a minikube cluster.

tufin deploy

Usage: tufin deploy

	•	Description: Deploys MySQL and WordPress pods into the cluster, connecting WordPress to MySQL.

tufin status

Usage: tufin status

	•	Description: Prints a status table containing the pod names and their current status in the default namespace.

Usage Instructions

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

	•	Description:
	•	Lists all pods in the default namespace along with their current status (e.g., Running, Pending).
	•	Example Output:

Getting pod status in default namespace...
Pod Name                       Status
mysql-xxxxxxxxxx-xxxxx         Running
wordpress-xxxxxxxxxx-xxxxx     Running



Accessing WordPress Application

Since the WordPress service is exposed via a NodePort, you can access it using the node’s IP address and the assigned port.

	1.	Find the NodePort:

kubectl get svc wordpress

Example Output:

NAME        TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
wordpress   NodePort   10.96.122.88   <none>        80:30080/TCP   5m

	•	NodePort: In this example, 30080.

	2.	Access WordPress:
	•	Open a web browser and navigate to http://<NODE_IP>:<NODE_PORT>.
	•	For a local cluster:
	•	k3s: Use localhost or 127.0.0.1.
	•	minikube: Run minikube ip to get the IP address.
	•	Example: http://localhost:30080 or http://<minikube-ip>:30080.

Code Structure and Explanation

Overview

The tufin client is structured using the Cobra library for command-line interfaces in Go. It uses the Kubernetes client-go library to interact with the Kubernetes cluster.

Main Components

	•	main.go: The primary Go file containing the implementation of the commands.

Commands

	1.	Cluster Command (tufin cluster [minikube]):
	•	Functionality:
	•	Checks if the chosen Kubernetes distribution (k3s or minikube) is installed.
	•	If not installed, attempts to install it automatically.
	•	Starts the cluster.
	•	Implementation Details:
	•	k3s:
	•	Installation: Uses the official installation script for Linux.
	•	Startup: Runs sudo k3s server with appropriate flags.
	•	minikube:
	•	Installation:
	•	Linux: Downloads the minikube binary and installs it.
	•	macOS: Installs via Homebrew.
	•	Startup: Runs minikube start.
	2.	Deploy Command (tufin deploy):
	•	Functionality:
	•	Connects to the cluster using the kubeconfig file.
	•	Deploys a MySQL deployment and service.
	•	Deploys a WordPress deployment and service, with environment variables set to connect to MySQL.
	•	Implementation Details:
	•	Uses client-go to create Kubernetes resources programmatically.
	•	Sets up necessary environment variables and ports.
	3.	Status Command (tufin status):
	•	Functionality:
	•	Lists all pods in the default namespace and displays their status.
	•	Implementation Details:
	•	Uses client-go to list pods and formats the output.

Adding Minikube for macOS Testing

	•	Reason:
	•	k3s installation and operation on macOS is not straightforward and often requires additional tools like virtual machines.
	•	Minikube provides an easy-to-use solution for running Kubernetes locally on macOS.
	•	Implementation:
	•	Modified the cluster command to accept an optional [minikube] argument.
	•	Added functions to handle the installation and startup of minikube.
	•	Adjusted the script to ensure compatibility with both k3s and minikube clusters.

Troubleshooting

Common Issues

	1.	Permission Denied Errors:
	•	Cause: Lack of administrative privileges during installation or startup.
	•	Solution:
	•	Run the commands with sudo if necessary.
	•	Ensure your user account has the required permissions.
	2.	kubeconfig Not Found:
	•	Cause: The kubeconfig file is not located at ~/.kube/config.
	•	Solution:
	•	Verify the location of your kubeconfig file.
	•	Set the KUBECONFIG environment variable if it’s in a different location.
	•	Modify the code to point to the correct path.
	3.	Cluster Not Starting:
	•	Cause: Conflicts with existing clusters or issues with the Kubernetes distribution.
	•	Solution:
	•	For k3s:
	•	Run sudo k3s-uninstall.sh to remove any existing k3s installations.
	•	For minikube:
	•	Run minikube delete to remove any existing clusters.
	•	Retry starting the cluster.
	4.	Pods Not Running:
	•	Cause: Issues with the deployment or insufficient resources.
	•	Solution:
	•	Check pod logs using kubectl logs.
	•	Ensure your system has enough resources (CPU, RAM).
	•	Verify that the images are correctly pulled from the registry.
	5.	Cannot Access WordPress Application:
	•	Cause: Networking issues or incorrect NodePort.
	•	Solution:
	•	Ensure the service is correctly exposed and note the correct NodePort.
	•	If using minikube, ensure you are using the correct minikube IP (minikube ip).
	6.	Failed to Pull Docker Images:
	•	Error Message:

stream logs failed container "mysql" in pod "mysql-xxxxxxxxx-xxxxx" is waiting to start: trying and failing to pull image for default/mysql-xxxxxxxxx-xxxxx (mysql)


	•	Cause: Kubernetes is unable to pull the required Docker images.
	•	Solutions:
	•	Verify Image Reference:
	•	Ensure the image name is correct. The MySQL image used is mysql:5.6.
	•	Check Network Connectivity:
	•	Ensure your nodes have internet access to pull images from Docker Hub.
	•	Set Image Pull Policy:
	•	Modify the deployment to set imagePullPolicy: Always.
	•	Use Image Pull Secrets:
	•	Create a Docker Hub account and create a Kubernetes secret:

kubectl create secret docker-registry regcred \
--docker-server=https://index.docker.io/v1/ \
--docker-username=your-username \
--docker-password=your-password \
--docker-email=your-email


	•	Update the deployment to use the secret:

Spec: corev1.PodSpec{
ImagePullSecrets: []corev1.LocalObjectReference{
{
Name: "regcred",
},
},
Containers: []corev1.Container{
// ...
},
},


	•	Alternative Registries:
	•	Use images from alternative registries if Docker Hub is inaccessible.

Conclusion

The tufin Go client provides a simple and efficient way to:

	•	Deploy a local Kubernetes cluster using k3s or minikube.
	•	Deploy MySQL and WordPress applications within the cluster.
	•	Monitor the status of your deployments.

Adding minikube support was essential for testing on macOS, given the complexities involved with installing and running k3s directly on macOS systems. Minikube offers a more accessible and user-friendly approach for macOS users, ensuring that the tufin client remains versatile and adaptable across different operating systems.

Feel free to customize and extend the client to suit your specific needs. Contributions and improvements are welcome!

References

	•	Go Programming Language: https://golang.org
	•	Cobra Library: https://github.com/spf13/cobra
	•	Kubernetes Client-Go: https://github.com/kubernetes/client-go
	•	k3s Documentation: https://docs.k3s.io
	•	Minikube Documentation: https://minikube.sigs.k8s.io/docs/
	•	kubectl Tool: https://kubernetes.io/docs/tasks/tools/
	•	Docker Hub Rate Limits: https://www.docker.com/increase-rate-limits
	•	Go exec Package: https://pkg.go.dev/os/exec
	•	Go runtime Package: https://pkg.go.dev/runtime

Feel free to reach out if you have any questions or need further assistance!