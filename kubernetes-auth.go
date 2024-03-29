package main 

import (
	"path/filepath"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Gets the application client set. This can be used to initialize the various Kubernetes clients.
// Uses in cliuster configuration when app is running inside the cluster, or kubeconfig file from
// home directory when running in development mode.
func GetClientSet() (*kubernetes.Clientset, error) {
	debugMode := os.Getenv("DEBUG_LOCAL")
	if debugMode != "" {
		return getOutOfClusterClientSet()
	} else {
		return getInClusterClientSet()
	}
}

func getInClusterClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getOutOfClusterClientSet() (*kubernetes.Clientset, error) {
	kubeconfigPath := filepath.Join(homeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	// Windows
	return os.Getenv("USERPROFILE")
}