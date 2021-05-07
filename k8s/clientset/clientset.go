package clientset

import (
	"os"

	introspectv1alpha1clientset "github.com/michaelhenkel/introspect/pkg/client/clientset"
	introspectrest "github.com/michaelhenkel/introspect/pkg/rest"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	contrailv1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/client/clientset_generated/clientset"
	contrailcorev1alpha1client "ssd-git.juniper.net/contrail/cn2/contrail/pkg/client/clientset_generated/clientset/typed/core/v1alpha1"
	deployerv1alpha1 "ssd-git.juniper.net/contrail/cn2/deployer/pkg/client/clientset_generated/clientset"
	deployerconfigplanev1alpha1client "ssd-git.juniper.net/contrail/cn2/deployer/pkg/client/clientset_generated/clientset/typed/configplane/v1alpha1"
	deployercontrolplanev1alpha1client "ssd-git.juniper.net/contrail/cn2/deployer/pkg/client/clientset_generated/clientset/typed/controlplane/v1alpha1"
	deployerdataplanev1alpha1client "ssd-git.juniper.net/contrail/cn2/deployer/pkg/client/clientset_generated/clientset/typed/dataplane/v1alpha1"
)

type Client struct {
	CoreV1             corev1client.CoreV1Interface
	ContrailCoreV1     contrailcorev1alpha1client.CoreV1alpha1Interface
	DeployerConfigV1   deployerconfigplanev1alpha1client.ConfigplaneV1alpha1Interface
	DeployerControlV1  deployercontrolplanev1alpha1client.ControlplaneV1alpha1Interface
	DeployerDataV1     deployerdataplanev1alpha1client.DataplaneV1alpha1Interface
	IntrospectClientV1 introspectv1alpha1clientset.Clientset
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
	deployerClientSet, err := deployerv1alpha1.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		CoreV1:             k8sClientset.CoreV1(),
		ContrailCoreV1:     contrailClientSet.CoreV1alpha1(),
		DeployerConfigV1:   deployerClientSet.ConfigplaneV1alpha1(),
		DeployerControlV1:  deployerClientSet.ControlplaneV1alpha1(),
		DeployerDataV1:     deployerClientSet.DataplaneV1alpha1(),
		IntrospectClientV1: introspectv1alpha1clientset.NewClientSet(introspectrest.NewRESTClient("http")),
	}, nil

}
