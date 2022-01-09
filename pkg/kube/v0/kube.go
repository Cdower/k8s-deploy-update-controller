package kube

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"deployUpdater/pkg/utils"
	// "k8s.io/apimachinery/pkg/api/errors"
	appsv1 "k8s.io/api/apps/v1"
	// apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	// in cluster
	"k8s.io/client-go/kubernetes"
	appsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	apiV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	// out of cluster
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubeClient struct {
	settings          *utils.Specification
	client            *kubernetes.Clientset
	deploymentsClient appsV1.DeploymentInterface
	podsClient        apiV1.PodInterface
}

func NewKubeClient(s *utils.Specification) *KubeClient {
	k := KubeClient{
		settings: s,
	}
	log.Printf("InCluster: %t\n", k.settings.InCluster)
	if k.settings.InCluster { // Load the in-cluster config
		log.Println("loading in cluster config")
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		k.client, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	} else { // Load kube config file
		log.Println("loading out of cluster config")
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		// create the clientset
		k.client, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}
	// Set Namespace in client
	k.deploymentsClient = k.client.AppsV1().Deployments(k.settings.Namespace)
	return &k
}

func (k *KubeClient) getContainerIndexFromName(deploy *appsv1.Deployment) int {
	for i, cntr := range deploy.Spec.Template.Spec.Containers {
		log.Printf("Index: %d, Name: %s", i, cntr.Name)
		if cntr.Name == k.settings.Container {
			log.Printf("Using: Index: %d, Name: %s", i, cntr.Name)
			return i
		}
	}
	return 0
}

func (k *KubeClient) getDeployment() *appsv1.Deployment {
	// if k.deploymentsClient == nil {
	// 	k.deploymentsClient = k.client.AppsV1().Deployments(k.settings.Namespace)
	// }
	result, getErr := k.deploymentsClient.Get(context.TODO(), k.settings.Deployment, metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("failed to get latest version of deployment: %v", getErr))
	}
	return result
}

func (k *KubeClient) GetPods() {
	if k.podsClient == nil {
		k.podsClient = k.client.CoreV1().Pods(k.settings.Namespace)
	}
	pods, err := k.podsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if k.settings.Namespace == "" {
		k.settings.Namespace = "all"
	}
	log.Printf("There are %d pods in the %s namespace\n", len(pods.Items), k.settings.Namespace)
}

func (k *KubeClient) PrintDeployments() {
	// if k.deploymentsClient == nil {
	// 	k.deploymentsClient = k.client.AppsV1().Deployments(k.settings.Namespace)
	// }
	deployList, err := k.deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if k.settings.Namespace == "" {
		k.settings.Namespace = "all"
	}
	log.Printf("There are %d Deployments in the %s namespace\n", len(deployList.Items), k.settings.Namespace)
	for _, d := range deployList.Items {
		log.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}

// result, getErr := deploymentsClient.Get(context.TODO(), "demo-deployment", metav1.GetOptions{})
func (k *KubeClient) GetDeploymentVersion() string {
	// if k.deploymentsClient == nil {
	// 	k.deploymentsClient = k.client.AppsV1().Deployments(k.settings.Namespace)
	// }
	result, getErr := k.deploymentsClient.Get(context.TODO(), k.settings.Deployment, metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("failed to get latest version of deployment: %v", getErr))
	}
	idx := k.getContainerIndexFromName(result)
	return strings.Split(result.Spec.Template.Spec.Containers[idx].Image, ":")[1]
}

func (k *KubeClient) UpdateDeploymentVersion(newVer string) error {
	// if k.deploymentsClient == nil {
	// 	k.deploymentsClient = k.client.AppsV1().Deployments(k.settings.Namespace)
	// }
	deploy := k.getDeployment()
	idx := k.getContainerIndexFromName(deploy)
	splitImage := strings.Split(deploy.Spec.Template.Spec.Containers[idx].Image, ":")
	splitImage[1] = newVer
	deploy.Spec.Template.Spec.Containers[idx].Image = strings.Join(splitImage, ":")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		updatedDeploy, updateErr := k.deploymentsClient.Update(context.TODO(), deploy, metav1.UpdateOptions{})
		log.Printf("new tag: %s\n", strings.Split(updatedDeploy.Spec.Template.Spec.Containers[idx].Image, ":"))
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	return nil
}
