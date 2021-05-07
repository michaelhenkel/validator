package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	introspectcontrolv1alpha1 "github.com/michaelhenkel/introspect/pkg/apis/control/v1alpha1"
)

const (
	introspectPort string = "8083"
)

type BGPNeighborNode struct {
	Resource      introspectcontrolv1alpha1.BgpNeighborResp
	Edges         []graph.NodeEdge
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *BGPNeighborNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*BGPNeighborNode)
	if !ok {
		return fmt.Errorf("not a BGPNeighborNode resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *BGPNeighborNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *BGPNeighborNode) Name() string {
	return r.Resource.Peer.Text
}

func (r *BGPNeighborNode) Type() graph.NodeType {
	return graph.BGPNeighbor
}

func (r *BGPNeighborNode) GetNodeEdges() []graph.NodeEdge {
	return r.Edges
}

func (r *BGPNeighborNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *BGPNeighborNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *BGPNeighborNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	var filterMap = make(map[string]struct{})

	controlList, err := g.ClientConfig.DeployerControlV1.Controls("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, controlResource := range controlList.Items {
		filterMap[controlResource.Name] = struct{}{}
	}

	for controlFilter := range filterMap {
		opts := metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", controlFilter),
		}
		podList, err := g.ClientConfig.CoreV1.Pods("").List(context.Background(), opts)
		if err != nil {
			return nil, err
		}
		for _, podResource := range podList.Items {
			bgpNeighborList, err := g.ClientConfig.IntrospectClientV1.Control(podResource.Status.PodIP + ":" + introspectPort).BgpNeighbors().List(context.Background())
			if err != nil {
				return nil, err
			}
			for _, resource := range bgpNeighborList.Neighbors.List.BgpNeighborResp {
				r.Resource = resource
				resourceNode := &BGPNeighborNode{
					Resource: resource,
				}
				graphNodeList = append(graphNodeList, resourceNode)
			}
		}
	}

	return graphNodeList, nil
}

/*

func addBGPNeighborNodes(validator *Validator) error {
	var graphNode graph.Node
	podNodeList := validator.graph.GetNodesByNodeType(graph.Pod)
	for _, podNodeInterface := range podNodeList {
		podNode, ok := podNodeInterface.(*PodNode)
		if !ok {
			return fmt.Errorf("not a podNode")
		}
		if podNode.Owner == graph.Control {
			bgpNeighborList, err := validator.clientConfig.Client.IntrospectClientV1.Control(podNode.Pod.Status.PodIP + ":" + introspectPort).BgpNeighbors().List(context.Background())
			if err != nil {
				return err
			}
			for _, bgpNeighbor := range bgpNeighborList.Neighbors.List.BgpNeighborResp {
				node := NewBGPNeighborNode(bgpNeighbor)
				graphNode = &node
				validator.graph.AddNode(graphNode)
			}
		}
	}
	return nil
}


func addBGPNeighborToBgpRouterEdge(validator *Validator) error {
	bgpRouterNodeList := validator.graph.GetNodesByNodeType(graph.BGPRouter)
	bgpNighborNodeList := validator.graph.GetNodesByNodeType(graph.BGPNeighbor)
	for _, bgpRouterNodeInterface := range bgpRouterNodeList {
		bgpRouterNode, ok := bgpRouterNodeInterface.(*BGPRouterNode)
		if !ok {
			return fmt.Errorf("not a bgpRouterNode")
		}
		for _, bgpNeighborNodeInterface := range bgpNighborNodeList {
			bgpNeighborNode, ok := bgpNeighborNodeInterface.(*BGPNeighborNode)
			if !ok {
				return fmt.Errorf("not a bgpNeighbor node")
			}
			if bgpRouterNode.BGPRouter.Spec.BGPRouterParameters.Address == contrailcorev1alpha1.IPAddress(bgpNeighborNode.BGPNeighbor.LocalAddress.Text) {
				validator.graph.AddEdge(bgpRouterNode, bgpNeighborNode, "")
			}
		}
	}
	return nil
}

func addBGPNeighborToVirtualRouterEdge(validator *Validator) error {
	virtualRouterNodeList := validator.graph.GetNodesByNodeType(graph.VirtualRouter)
	bgpNighborNodeList := validator.graph.GetNodesByNodeType(graph.BGPNeighbor)
	for _, virtualRouterNodeInterface := range virtualRouterNodeList {
		virtualRouterNode, ok := virtualRouterNodeInterface.(*VirtualRouterNode)
		if !ok {
			return fmt.Errorf("not a virtualRouterNode")
		}
		for _, bgpNeighborNodeInterface := range bgpNighborNodeList {
			bgpNeighborNode, ok := bgpNeighborNodeInterface.(*BGPNeighborNode)
			if !ok {
				return fmt.Errorf("not a bgpNeighbor node")
			}
			if virtualRouterNode.VirtualRouter.Spec.IPAddress == contrailcorev1alpha1.IPAddress(bgpNeighborNode.BGPNeighbor.PeerAddress.Text) {
				validator.graph.AddEdge(bgpNeighborNode, virtualRouterNode, "")
			}
		}
	}
	return nil
}
*/
