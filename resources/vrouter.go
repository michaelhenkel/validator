package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	dpv1alpha1 "ssd-git.juniper.net/contrail/cn2/deployer/pkg/apis/dataplane/v1alpha1"
)

type VrouterNode struct {
	Resource      dpv1alpha1.Vrouter
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *VrouterNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*VrouterNode)
	if !ok {
		return fmt.Errorf("not a virtualrouter resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *VrouterNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *VrouterNode) Name() string {
	return r.Resource.Name
}

func (r *VrouterNode) Type() graph.NodeType {
	return graph.Vrouter
}

func (r *VrouterNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *VrouterNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *VrouterNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	resourceList, err := g.ClientConfig.DeployerDataV1.Vrouters("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range resourceList.Items {
		resourceNode := &VrouterNode{
			Resource: resource,
			EdgeSelectors: []graph.EdgeSelector{{
				NodeType: graph.Pod,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"app": resource.Name},
				}},
			}},
		}
		graphNodeList = append(graphNodeList, resourceNode)
	}
	return graphNodeList, nil
}
