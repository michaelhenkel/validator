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
	Edges         []graph.NodeEdge
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

func (r *ControlNode) GetNodeEdges() []graph.NodeEdge {
	return r.Edges
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
		}
		graphNodeList = append(graphNodeList, nodeResource)
	}
	return graphNodeList, nil
}

/*

func addControlToPodEdges(validator *Validator) error {
	nodeList := validator.graph.GetNodesByNodeType(graph.Control)
	podNodeList := validator.graph.GetNodesByNodeType(graph.Pod)
	for _, nodeInterface := range nodeList {
		node, ok := nodeInterface.(*ControlNode)
		if !ok {
			return fmt.Errorf("not a vrouter node")
		}
		for _, podNodeInterface := range podNodeList {
			podNode, ok := podNodeInterface.(*PodNode)
			if !ok {
				return fmt.Errorf("not a pod node")
			}
			if appName, ok := podNode.Pod.Labels["app"]; ok && appName == node.Control.Name {
				validator.graph.AddEdge(node, podNode, "")
			}

		}
	}
	return nil
}
*/
