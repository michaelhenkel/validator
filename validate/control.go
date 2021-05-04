package validate

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cpv1alpha1 "ssd-git.juniper.net/contrail/cn2/deployer/pkg/apis/controlplane/v1alpha1"
)

type ControlNode struct {
	Control  cpv1alpha1.Control
	ID       int64
	NodeType graph.NodeType
}

func NewControlNode(vrouter cpv1alpha1.Control) ControlNode {
	return ControlNode{
		Control:  vrouter,
		NodeType: graph.Control,
	}
}

func (v *ControlNode) GetNodeType() graph.NodeType {
	return v.NodeType
}

func (v *ControlNode) GetID() int64 {
	return v.ID
}

func (v *ControlNode) SetID(id int64) {
	v.ID = id
}

func (v *ControlNode) GetName() string {
	return v.Control.Name
}

func (v *ControlNode) Shape() graph.Shape {
	return graph.DeploymentResource
}

func addControlNodes(validator *Validator) error {
	controlList, err := validator.clientConfig.Client.DeployerControlV1.Controls("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	var graphNode graph.Node
	for _, control := range controlList.Items {
		node := NewControlNode(control)
		graphNode = &node
		validator.graph.AddNode(graphNode)
	}
	return nil
}

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
