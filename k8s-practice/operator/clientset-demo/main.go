package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// v0.20.4
func main() {
	//getPods()
	getDeployment()
}
// 查询pod个数
func getPods() {
	config := getConfig()
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	}
}

// 创建deploy
func createDeployment(ctx context.Context, ns string) error {
	// creates the in-cluster config
	config := getConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	deployment := &appsv1.Deployment{}
	deployment.Name = "example"
	// edit deployment spec

	client := clientset.AppsV1().Deployments(ns)
	_, err = client.Create(ctx, deployment, metav1.CreateOptions{})
	return err
}

func getConfig() *rest.Config {
	var err error
	var config *rest.Config
	var kubeconfig *string

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "kind-config-my-cluster"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// 使用 ServiceAccount 创建集群配置（InCluster模式）
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用 KubeConfig 文件创建集群配置
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}

	return config
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// 查询deploy个数
func getDeployment() {

	config := getConfig()
	// 创建 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 使用 clientset 获取 Deployments
	deployments, err := clientset.AppsV1().Deployments("kube-system").List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for idx, deploy := range deployments.Items {
		fmt.Printf("%d -> %s\n", idx+1, deploy.Name)
	}

}
