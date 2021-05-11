package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	configplanev1alpha1 "ssd-git.juniper.net/contrail/cn2/deployer/pkg/apis/configplane/v1alpha1"
)

type KubemanagerNode struct {
	Resource      configplanev1alpha1.Kubemanager
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *KubemanagerNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*KubemanagerNode)
	if !ok {
		return fmt.Errorf("not a KubemanagerNode resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *KubemanagerNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *KubemanagerNode) Name() string {
	return r.Resource.Name
}

func (r *KubemanagerNode) Type() graph.NodeType {
	return graph.Kubemanager
}

func (r *KubemanagerNode) Plane() graph.Plane {
	return plane
}

func (r *KubemanagerNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *KubemanagerNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *KubemanagerNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	resourceList, err := g.ClientConfig.DeployerConfigV1.Kubemanagers("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range resourceList.Items {
		r.Resource = resource
		resourceNode := &KubemanagerNode{
			Resource: resource,
			EdgeLabels: []graph.EdgeLabel{{
				Value: map[string]string{"DefaultPodNetworkNamespaceName": fmt.Sprintf("%s/%s-default-podnetwork", resource.Spec.Namespace, resource.Spec.ClusterName)},
			}, {
				Value: map[string]string{"DefaultServiceNetworkNamespaceName": fmt.Sprintf("%s/%s-default-servicenetwork", resource.Spec.Namespace, resource.Spec.ClusterName)},
			}},
		}
		graphNodeList = append(graphNodeList, resourceNode)
	}
	return graphNodeList, nil
}
