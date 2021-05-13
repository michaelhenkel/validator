package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type RoutingInstanceNode struct {
	Resource      contrailcorev1alpha1.RoutingInstance
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *RoutingInstanceNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*RoutingInstanceNode)
	if !ok {
		return fmt.Errorf("not a RoutingInstanceNode resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *RoutingInstanceNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *RoutingInstanceNode) Name() string {
	return r.Resource.Name
}

func (r *RoutingInstanceNode) Type() graph.NodeType {
	return graph.RoutingInstance
}

func (r *RoutingInstanceNode) Plane() graph.Plane {
	return plane
}

func (r *RoutingInstanceNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *RoutingInstanceNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *RoutingInstanceNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	resourceList, err := g.ClientConfig.ContrailCoreV1.RoutingInstances("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range resourceList.Items {
		r.Resource = resource
		var edgeSelectorList []graph.EdgeSelector
		if resource.Spec.Parent.Kind == "VirtualNetwork" {
			edgeSelector := graph.EdgeSelector{
				NodeType: graph.VirtualNetwork,
				Plane:    graph.ConfigPlane,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"VirtualNetworkNamespaceName": fmt.Sprintf("%s/%s", resource.Spec.Parent.Namespace, resource.Spec.Parent.Name)},
				}},
			}
			edgeSelectorList = append(edgeSelectorList, edgeSelector)
		}
		edgeSelector := graph.EdgeSelector{
			NodeType: graph.RoutingInstance,
			Plane:    graph.ControlPlane,
			MatchValues: []graph.MatchValue{{
				Value: map[string]string{"RoutingInstanceName": fmt.Sprintf("%s/%s", resource.Namespace, resource.Name)},
			}},
		}
		edgeSelectorList = append(edgeSelectorList, edgeSelector)
		resourceNode := &RoutingInstanceNode{
			Resource: resource,
			EdgeLabels: []graph.EdgeLabel{{
				Value: map[string]string{"RoutingInstanceName": resource.Name},
			}},
			EdgeSelectors: edgeSelectorList,
		}
		graphNodeList = append(graphNodeList, resourceNode)
	}
	return graphNodeList, nil
}
