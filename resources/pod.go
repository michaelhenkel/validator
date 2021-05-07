package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodNode struct {
	Resource      corev1.Pod
	Edges         []graph.NodeEdge
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *PodNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*PodNode)
	if !ok {
		return fmt.Errorf("not a virtualrouter resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *PodNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *PodNode) Name() string {
	return r.Resource.Name
}

func (r *PodNode) Type() graph.NodeType {
	return graph.Pod
}

func (r *PodNode) GetNodeEdges() []graph.NodeEdge {
	return r.Edges
}

func (r *PodNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *PodNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *PodNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var nodeInterfaceList []graph.NodeInterface
	vrouterList, err := g.ClientConfig.DeployerDataV1.Vrouters("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range vrouterList.Items {
		vrouterNodeList, err := r.buildNodeList(resource.Name, g, graph.Vrouter)
		if err != nil {
			return nil, err
		}
		nodeInterfaceList = append(nodeInterfaceList, vrouterNodeList...)
	}
	controlList, err := g.ClientConfig.DeployerControlV1.Controls("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range controlList.Items {
		controlList, err := r.buildNodeList(resource.Name, g, graph.Control)
		if err != nil {
			return nil, err
		}
		nodeInterfaceList = append(nodeInterfaceList, controlList...)
	}
	return nodeInterfaceList, nil
}

func (r *PodNode) buildNodeList(filter string, g *graph.Graph, nodeType graph.NodeType) ([]graph.NodeInterface, error) {
	var nodeInterfaceList []graph.NodeInterface
	opts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", filter),
	}
	resourceList, err := g.ClientConfig.CoreV1.Pods("").List(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	for _, resource := range resourceList.Items {
		var edgeList []graph.NodeEdge
		var edgeSelectorList []graph.EdgeSelector
		for _, volume := range resource.Spec.Volumes {
			if volume.Name == "config-volume" {
				if volume.VolumeSource.ConfigMap != nil {
					nodeEdge := graph.NodeEdge{
						To:          graph.ConfigMap,
						MatchValues: []map[string]string{{"ConfigMap": volume.ConfigMap.Name}},
					}
					edgeList = append(edgeList, nodeEdge)
					edgeSelector := graph.EdgeSelector{
						NodeType: graph.ConfigMap,
						MatchValues: []graph.MatchValue{{
							Value: map[string]string{"ConfigMap": volume.ConfigMap.Name},
						}},
					}
					edgeSelectorList = append(edgeSelectorList, edgeSelector)
				}
			}
		}
		virtualRouterEdge := graph.NodeEdge{
			To: graph.VirtualRouter,
			MatchValues: []map[string]string{{
				"PodIP": resource.Status.PodIP,
			}, {
				"NodeType": string(nodeType),
			}},
		}
		edgeList = append(edgeList, virtualRouterEdge)

		bgpRouterEdge := graph.NodeEdge{
			To:          graph.BGPRouter,
			MatchValues: []map[string]string{{"PodIP": resource.Status.PodIP}},
		}
		edgeList = append(edgeList, bgpRouterEdge)

		r.Resource = resource
		resourceNode := &PodNode{
			Resource:      resource,
			Edges:         edgeList,
			EdgeSelectors: edgeSelectorList,
		}
		nodeInterfaceList = append(nodeInterfaceList, resourceNode)
	}
	return nodeInterfaceList, nil
}

/*

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

*/
