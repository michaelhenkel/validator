package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8SNodeNode struct {
	Resource      corev1.Node
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *K8SNodeNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*K8SNodeNode)
	if !ok {
		return fmt.Errorf("not a k8s node resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *K8SNodeNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *K8SNodeNode) Name() string {
	return r.Resource.Name
}

func (r *K8SNodeNode) Type() graph.NodeType {
	return graph.K8SNode
}

func (r *K8SNodeNode) Plane() graph.Plane {
	return plane
}

func (r *K8SNodeNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *K8SNodeNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *K8SNodeNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface

	resourceList, err := g.ClientConfig.CoreV1.Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range resourceList.Items {
		var internalIP, hostname string
		for _, address := range resource.Status.Addresses {
			if address.Type == "InternalIP" {
				internalIP = address.Address
			}
			if address.Type == "Hostname" {
				hostname = address.Address
			}
		}
		resourceNode := &K8SNodeNode{
			Resource: resource,
			EdgeLabels: []graph.EdgeLabel{{
				Value: map[string]string{"NodeIP": internalIP},
			}, {
				Value: map[string]string{"Hostname": hostname},
			}},
		}
		graphNodeList = append(graphNodeList, resourceNode)
	}

	return graphNodeList, nil
}
