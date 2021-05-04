package validate

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	dpv1alpha1 "ssd-git.juniper.net/contrail/cn2/deployer/pkg/apis/dataplane/v1alpha1"
)

type VrouterNode struct {
	Vrouter  dpv1alpha1.Vrouter
	ID       int64
	NodeType graph.NodeType
}

func NewVrouterNode(vrouter dpv1alpha1.Vrouter) VrouterNode {
	return VrouterNode{
		Vrouter:  vrouter,
		NodeType: graph.Vrouter,
	}
}

func (v *VrouterNode) GetNodeType() graph.NodeType {
	return v.NodeType
}

func (v *VrouterNode) GetID() int64 {
	return v.ID
}

func (v *VrouterNode) SetID(id int64) {
	v.ID = id
}

func (v *VrouterNode) GetName() string {
	return v.Vrouter.Name
}

func (v *VrouterNode) Shape() graph.Shape {
	return graph.DeploymentResource
}

func addVrouterNodes(validator *Validator) error {
	vrouterList, err := validator.clientConfig.Client.DeployerDataV1.Vrouters("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	var graphNode graph.Node
	for _, vrouter := range vrouterList.Items {
		node := NewVrouterNode(vrouter)
		graphNode = &node
		validator.graph.AddNode(graphNode)
	}
	return nil
}

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
