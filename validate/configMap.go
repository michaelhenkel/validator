package validate

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapNode struct {
	ConfigMap corev1.ConfigMap
	ID        int64
	NodeType  graph.NodeType
}

func NewConfigMapNode(configMap corev1.ConfigMap) ConfigMapNode {
	return ConfigMapNode{
		ConfigMap: configMap,
		NodeType:  graph.ConfigMap,
	}
}

func (v *ConfigMapNode) GetNodeType() graph.NodeType {
	return v.NodeType
}

func (v *ConfigMapNode) GetID() int64 {
	return v.ID
}

func (v *ConfigMapNode) SetID(id int64) {
	v.ID = id
}

func (v *ConfigMapNode) GetName() string {
	return v.ConfigMap.Name
}

func (v *ConfigMapNode) Shape() graph.Shape {
	return graph.DeploymentResource
}

func addConfigMapNodes(validator *Validator, nodeType graph.NodeType) error {
	var graphNode graph.Node
	nodes := validator.graph.GetNodesByNodeType(nodeType)
	var configMap *corev1.ConfigMap
	var err error
	for _, nodeInterface := range nodes {
		switch nodeType {
		case graph.Vrouter:
			node, ok := nodeInterface.(*VrouterNode)
			if !ok {
				return fmt.Errorf("not a vrouter node")
			}
			configMap, err = validator.clientConfig.Client.CoreV1.ConfigMaps(node.Vrouter.Namespace).Get(context.Background(), node.Vrouter.Name+"-configmap", metav1.GetOptions{})
			if err != nil {
				return err
			}
		case graph.Control:
			node, ok := nodeInterface.(*ControlNode)
			if !ok {
				return fmt.Errorf("not a control node")
			}
			configMap, err = validator.clientConfig.Client.CoreV1.ConfigMaps(node.Control.Namespace).Get(context.Background(), node.Control.Name+"-configmap", metav1.GetOptions{})
			if err != nil {
				return err
			}
		}
		configMapNode := NewConfigMapNode(*configMap)
		graphNode = &configMapNode
		validator.graph.AddNode(graphNode)
	}
	return nil
}

func addConfigMapToPodEdges(validator *Validator, owner graph.NodeType) error {
	nodeList := validator.graph.GetNodesByNodeType(graph.ConfigMap)
	podNodeList := validator.graph.GetNodesByNodeType(graph.Pod)
	for _, nodeInterface := range nodeList {
		node, ok := nodeInterface.(*ConfigMapNode)
		if !ok {
			return fmt.Errorf("not a configMap node")
		}
		for _, podNodeInterface := range podNodeList {
			podNode, ok := podNodeInterface.(*PodNode)
			if !ok {
				return fmt.Errorf("not a pod node")
			}
			if podNode.Owner == owner {
				for _, volume := range podNode.Pod.Spec.Volumes {
					if volume.Name == "config-volume" {
						if volume.VolumeSource.ConfigMap != nil && volume.VolumeSource.ConfigMap.Name == node.ConfigMap.Name {
							validator.graph.AddEdge(podNode, node, "")
						}
					}
				}

			}
		}
	}
	return nil
}
