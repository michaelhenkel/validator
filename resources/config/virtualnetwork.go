package resources

import (
	"context"
	"fmt"

	"github.com/s3kim2018/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type VirtualNetworkNode struct {
	Resource      contrailcorev1alpha1.VirtualNetwork
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *VirtualNetworkNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*VirtualNetworkNode)
	if !ok {
		return fmt.Errorf("not a VirtualNetworkNode resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *VirtualNetworkNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *VirtualNetworkNode) Name() string {
	return fmt.Sprintf("%s:%s", r.Resource.Namespace, r.Resource.Name)
}

func (r *VirtualNetworkNode) Type() graph.NodeType {
	return graph.VirtualNetwork
}

func (r *VirtualNetworkNode) Plane() graph.Plane {
	return plane
}

func (r *VirtualNetworkNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *VirtualNetworkNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *VirtualNetworkNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	resourceList, err := g.ClientConfig.ContrailCoreV1.VirtualNetworks("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range resourceList.Items {
		r.Resource = resource

		resourceNode := &VirtualNetworkNode{
			Resource: resource,
			EdgeLabels: []graph.EdgeLabel{{
				Value: map[string]string{"VirtualNetworkNamespaceName": fmt.Sprintf("%s/%s", resource.Namespace, resource.Name)},
			}},
		}
		graphNodeList = append(graphNodeList, resourceNode)
	}
	return graphNodeList, nil
}
