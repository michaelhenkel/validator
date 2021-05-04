package validate

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type PodNode struct {
	Pod      corev1.Pod
	ID       int64
	NodeType graph.NodeType
	Owner    graph.NodeType
}

func NewPodNode(pod corev1.Pod, owner graph.NodeType) PodNode {
	return PodNode{
		Pod:      pod,
		NodeType: graph.Pod,
		Owner:    owner,
	}
}

func (v *PodNode) GetNodeType() graph.NodeType {
	return v.NodeType
}

func (v *PodNode) GetID() int64 {
	return v.ID
}

func (v *PodNode) SetID(id int64) {
	v.ID = id
}

func (v *PodNode) GetName() string {
	return v.Pod.Name
}

func (v *PodNode) Shape() graph.Shape {
	return graph.ConfigResource
}

func addPodNodes(validator *Validator, nodeType graph.NodeType) error {
	var graphNode graph.Node
	nodes := validator.graph.GetNodesByNodeType(nodeType)
	switch nodeType {
	case graph.Vrouter:
		for _, nodeInterface := range nodes {
			node, ok := nodeInterface.(*VrouterNode)
			if !ok {
				return fmt.Errorf("not a vrouter node")
			}
			opts := metav1.ListOptions{
				LabelSelector: fmt.Sprintf("app=%s", node.Vrouter.Name),
			}
			pl, err := validator.clientConfig.Client.CoreV1.Pods("").List(context.Background(), opts)
			if err != nil {
				return err
			}
			for _, pod := range pl.Items {
				node := NewPodNode(pod, graph.Vrouter)
				graphNode = &node
				validator.graph.AddNode(graphNode)
			}
		}
	case graph.Control:
		for _, nodeInterface := range nodes {
			node, ok := nodeInterface.(*ControlNode)
			if !ok {
				return fmt.Errorf("not a vrouter node")
			}
			opts := metav1.ListOptions{
				LabelSelector: fmt.Sprintf("app=%s", node.Control.Name),
			}
			pl, err := validator.clientConfig.Client.CoreV1.Pods("").List(context.Background(), opts)
			if err != nil {
				return err
			}
			for _, pod := range pl.Items {
				node := NewPodNode(pod, graph.Control)
				graphNode = &node
				validator.graph.AddNode(graphNode)
			}
		}
	}
	return nil
}

func addPodToVirtualRouterEdges(validator *Validator) error {
	virtualRouterNodeList := validator.graph.GetNodesByNodeType(graph.VirtualRouter)
	podNodeList := validator.graph.GetNodesByNodeType(graph.Pod)
	for _, podNodeInterface := range podNodeList {
		podNode, ok := podNodeInterface.(*PodNode)
		if !ok {
			return fmt.Errorf("not a pod node")
		}
		if podNode.Owner == graph.Vrouter {
			for _, nodeInterface := range virtualRouterNodeList {
				node, ok := nodeInterface.(*VirtualRouterNode)
				if !ok {
					return fmt.Errorf("not a virtualRouter node")
				}
				if node.VirtualRouter.Spec.IPAddress == contrailcorev1alpha1.IPAddress(podNode.Pod.Status.PodIP) {
					validator.graph.AddEdge(podNode, node, "")
				}
			}
		}
	}
	return nil
}

func addPodToBGPRouterEdges(validator *Validator) error {
	nodeList := validator.graph.GetNodesByNodeType(graph.BGPRouter)
	podNodeList := validator.graph.GetNodesByNodeType(graph.Pod)
	for _, podNodeInterface := range podNodeList {
		podNode, ok := podNodeInterface.(*PodNode)
		if !ok {
			return fmt.Errorf("not a pod node")
		}
		if podNode.Owner == graph.Control {
			for _, nodeInterface := range nodeList {
				node, ok := nodeInterface.(*BGPRouterNode)
				if !ok {
					return fmt.Errorf("not a bgpRouter node")
				}
				if node.BGPRouter.Spec.BGPRouterParameters.Address == contrailcorev1alpha1.IPAddress(podNode.Pod.Status.PodIP) {
					validator.graph.AddEdge(podNode, node, "")
				}
			}
		}
	}
	return nil
}
