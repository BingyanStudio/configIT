package utils

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	k8sClient *kubernetes.Clientset
	k8sConfig *rest.Config
)

func InitK8sClient() error {
	var err error
	k8sConfig, err = rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	k8sClient, err = kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}
	return nil
}

func GetNamespaceIPs(namespace string) ([]string, error) {
	pods, err := k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	var podIPs []string
	for _, pod := range pods.Items {
		if pod.Status.Phase == "Running" {
			podIPs = append(podIPs, pod.Status.PodIP)
		}
	}
	return podIPs, nil
}

func GetPodIP(namespace, podName string, fuzzmatch bool) (string, error) {
	if fuzzmatch {
		pods, err := k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to list pods: %w", err)
		}
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, podName) {
				return pod.Status.PodIP, nil
			}
		}
		return "", fmt.Errorf("pod with name %s not found in namespace %s", podName, namespace)
	} else {
		pod, err := k8sClient.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return "", fmt.Errorf("pod %s not found in namespace %s", podName, namespace)
			}
			return "", fmt.Errorf("failed to get pod: %w", err)
		}
		return pod.Status.PodIP, nil
	}
}

func GetNamespaces() ([]string, error) {
	namespaces, err := k8sClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	var nsList []string
	for _, ns := range namespaces.Items {
		nsList = append(nsList, ns.Name)
	}
	return nsList, nil
}

func GetPods(namespace string) ([]string, error) {
	pods, err := k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	var podList []string
	for _, pod := range pods.Items {
		podList = append(podList, pod.Name)
	}
	return podList, nil
}

func CreateConfigMap(namespace, name string, data map[string]string) error {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	_, err := k8sClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create configmap: %w", err)
	}
	return nil
}

func UpdateConfigMap(namespace, name string, data map[string]string) error {
	configMap, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("configmap %s not found in namespace %s", name, namespace)
		}
		return fmt.Errorf("failed to get configmap: %w", err)
	}

	configMap.Data = data
	_, err = k8sClient.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update configmap: %w", err)
	}
	return nil
}
