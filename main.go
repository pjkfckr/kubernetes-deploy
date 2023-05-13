package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func main() {
	var (
		client           *kubernetes.Clientset
		deploymentLabels map[string]string
		err              error
	)
	ctx := context.Background()

	if client, err = getClient(); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	if deploymentLabels, err = deploy(ctx, client); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	fmt.Printf("deploy finished Did a deploy with labels: %+v\n", deploymentLabels)
}

func getClient() (*kubernetes.Clientset, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func deploy(ctx context.Context, client *kubernetes.Clientset) (map[string]string, error) {
	var deployment *v1.Deployment

	appFile, err := os.ReadFile("app.yaml")
	if err != nil {
		return nil, fmt.Errorf("readfile error: %s", err)
	}

	obj, groupVersionKind, err := scheme.Codecs.UniversalDeserializer().Decode(appFile, nil, nil)

	switch obj.(type) {
	case *v1.Deployment:
		deployment = obj.(*v1.Deployment)
	default:
		return nil, fmt.Errorf("unrecongnized type: %s\n", groupVersionKind)
	}

	deploymentResponse, err := client.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("deploy Error: %s", err)
	}
	return deploymentResponse.Spec.Template.Labels, nil

}
