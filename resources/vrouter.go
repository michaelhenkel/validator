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
	Edges         []graph.NodeEdge
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

func (r *VrouterNode) GetNodeEdges() []graph.NodeEdge {
	return r.Edges
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
		}
		graphNodeList = append(graphNodeList, resourceNode)
	}
	return graphNodeList, nil
}

/*

func addVrouterToPodEdges(validator *Validator) error {
	nodeList := validator.graph.GetNodesByNodeType(graph.Vrouter)
	podNodeList := validator.graph.GetNodesByNodeType(graph.Pod)
	for _, nodeInterface := range nodeList {
		node, ok := nodeInterface.(*VrouterNode)
		if !ok {
			return fmt.Errorf("not a vrouter node")
		}
		for _, podNodeInterface := range podNodeList {
			podNode, ok := podNodeInterface.(*PodNode)
			if !ok {
				return fmt.Errorf("not a pod node")
			}
			if appName, ok := podNode.Pod.Labels["app"]; ok && appName == node.Vrouter.Name {
				validator.graph.AddEdge(node, podNode, "")
			}

		}
	}
	return nil
}
*/
