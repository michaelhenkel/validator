package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type VirtualRouterNode struct {
	Resource      contrailcorev1alpha1.VirtualRouter
	Edges         []graph.NodeEdge
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *VirtualRouterNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*VirtualRouterNode)
	if !ok {
		return fmt.Errorf("not a virtualrouter resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *VirtualRouterNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *VirtualRouterNode) Name() string {
	return r.Resource.Name
}

func (r *VirtualRouterNode) Type() graph.NodeType {
	return graph.VirtualRouter
}

func (r *VirtualRouterNode) GetNodeEdges() []graph.NodeEdge {
	return r.Edges
}

func (r *VirtualRouterNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *VirtualRouterNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *VirtualRouterNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	virtualRouterList, err := g.ClientConfig.ContrailCoreV1.VirtualRouters().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, virtualRouter := range virtualRouterList.Items {
		r.Resource = virtualRouter
		virtualRouterNode := &VirtualRouterNode{
			Resource: virtualRouter,
			Edges: []graph.NodeEdge{{
				To: graph.Pod,
				MatchValues: []map[string]string{{
					"PodIP": string(virtualRouter.Spec.IPAddress),
				}, {
					"NodeType": string(graph.Vrouter),
				}},
			}},
		}
		graphNodeList = append(graphNodeList, virtualRouterNode)
	}
	return graphNodeList, nil
}