package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type BGPRouterNode struct {
	Resource      contrailcorev1alpha1.BGPRouter
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *BGPRouterNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*BGPRouterNode)
	if !ok {
		return fmt.Errorf("not a BGPRouterNode resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *BGPRouterNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *BGPRouterNode) Name() string {
	return r.Resource.Name
}

func (r *BGPRouterNode) Type() graph.NodeType {
	return graph.BGPRouter
}

func (r *BGPRouterNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *BGPRouterNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *BGPRouterNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	resourceList, err := g.ClientConfig.ContrailCoreV1.BGPRouters("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range resourceList.Items {
		r.Resource = resource
		resourceNode := &BGPRouterNode{
			Resource: resource,
			EdgeLabels: []graph.EdgeLabel{{
				Value: map[string]string{"BGPRouterIP": string(resource.Spec.BGPRouterParameters.Address)},
			}},
			EdgeSelectors: []graph.EdgeSelector{{
				NodeType: graph.BGPNeighbor,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"BGPRouterNeighborLocalIP": string(resource.Spec.BGPRouterParameters.Address)},
				}},
			}},
		}
		graphNodeList = append(graphNodeList, resourceNode)
	}
	return graphNodeList, nil
}
