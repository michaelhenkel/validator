package validate

import (
	"context"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type BGPRouterNode struct {
	BGPRouter contrailcorev1alpha1.BGPRouter
	ID        int64
	NodeType  graph.NodeType
}

func NewBGPRouterNode(virtualRouter contrailcorev1alpha1.BGPRouter) BGPRouterNode {
	return BGPRouterNode{
		BGPRouter: virtualRouter,
		NodeType:  graph.BGPRouter,
	}
}

func (v *BGPRouterNode) GetNodeType() graph.NodeType {
	return v.NodeType
}

func (v *BGPRouterNode) GetID() int64 {
	return v.ID
}

func (v *BGPRouterNode) SetID(id int64) {
	v.ID = id
}

func (v *BGPRouterNode) GetName() string {
	return v.BGPRouter.Name
}

func (v *BGPRouterNode) Shape() graph.Shape {
	return graph.ConfigResource
}

func addBGPRouterNodes(validator *Validator) error {
	bgpRouterList, err := validator.clientConfig.Client.ContrailCoreV1.BGPRouters("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	var graphNode graph.Node
	for _, bgpRouter := range bgpRouterList.Items {
		node := NewBGPRouterNode(bgpRouter)
		graphNode = &node
		validator.graph.AddNode(graphNode)
	}
	return nil
}
