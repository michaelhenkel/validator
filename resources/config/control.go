package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cpv1alpha1 "ssd-git.juniper.net/contrail/cn2/deployer/pkg/apis/controlplane/v1alpha1"
)

type ControlNode struct {
	Resource      cpv1alpha1.Control
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *ControlNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*ControlNode)
	if !ok {
		return fmt.Errorf("not a virtualrouter resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *ControlNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *ControlNode) Name() string {
	return r.Resource.Name
}

func (r *ControlNode) Type() graph.NodeType {
	return graph.Control
}

func (r *ControlNode) Plane() graph.Plane {
	return plane
}

func (r *ControlNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *ControlNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *ControlNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	resourceList, err := g.ClientConfig.DeployerControlV1.Controls("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range resourceList.Items {
		r.Resource = resource
		nodeResource := &ControlNode{
			Resource: resource,
			EdgeSelectors: []graph.EdgeSelector{{
				NodeType: graph.BGPRouter,
				Plane:    graph.ConfigPlane,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"BGPRouterIP": resource.Spec.HostIP},
				}},
			}, {
				NodeType: graph.Pod,
				Plane:    graph.ConfigPlane,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"app": resource.Name},
				}},
			}},
		}
		graphNodeList = append(graphNodeList, nodeResource)
	}
	return graphNodeList, nil
}
