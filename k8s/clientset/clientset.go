package clientset

import (
	"os"

	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	contrailv1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/client/clientset_generated/clientset"
	contrailv1alpha1client "ssd-git.juniper.net/contrail/cn2/contrail/pkg/client/clientset_generated/clientset/typed/core/v1alpha1"
)

type Client struct {
	CoreV1     corev1client.CoreV1Interface
	ContrailV1 contrailv1alpha1client.CoreV1alpha1Interface
}

func NewClient(kubeconfigPath string) (*Client, error) {
	var kubeconfig string
	if kubeconfigPath == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}
	// Create a Config (k8s.io/client-go/rest)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	k8sClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	contrailClientSet, err := contrailv1alpha1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		CoreV1:     k8sClientset.CoreV1(),
		ContrailV1: contrailClientSet.CoreV1alpha1(),
	}, nil

}
