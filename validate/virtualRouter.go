package validate

import (
	"context"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type VirtualRouterNode struct {
	VirtualRouter contrailcorev1alpha1.VirtualRouter
	ID            int64
	NodeType      graph.NodeType
}

func NewVirtualRouterNode(virtualRouter contrailcorev1alpha1.VirtualRouter) VirtualRouterNode {
	return VirtualRouterNode{
		VirtualRouter: virtualRouter,
		NodeType:      graph.VirtualRouter,
	}
}

func (v *VirtualRouterNode) GetNodeType() graph.NodeType {
	return v.NodeType
}

func (v *VirtualRouterNode) GetID() int64 {
	return v.ID
}

func (v *VirtualRouterNode) SetID(id int64) {
	v.ID = id
}

func (v *VirtualRouterNode) GetName() string {
	return v.VirtualRouter.Name
}

func (v *VirtualRouterNode) Shape() graph.Shape {
	return graph.ConfigResource
}

func addVirtualRouterNodes(validator *Validator) error {
	virtualRouterList, err := validator.clientConfig.Client.ContrailCoreV1.VirtualRouters().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	var graphNode graph.Node
	for _, virtualRouter := range virtualRouterList.Items {
		node := NewVirtualRouterNode(virtualRouter)
		graphNode = &node
		validator.graph.AddNode(graphNode)
	}
	return nil
}
