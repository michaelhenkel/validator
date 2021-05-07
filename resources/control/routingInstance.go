package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	introspectcontrolv1alpha1 "github.com/michaelhenkel/introspect/pkg/apis/control/v1alpha1"
)

type RoutingInstanceNode struct {
	Resource      introspectcontrolv1alpha1.ShowRoutingInstance
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *RoutingInstanceNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*RoutingInstanceNode)
	if !ok {
		return fmt.Errorf("not a RoutingInstanceNode resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *RoutingInstanceNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *RoutingInstanceNode) Name() string {
	return r.Resource.Name.Text
}

func (r *RoutingInstanceNode) Type() graph.NodeType {
	return graph.RoutingInstance
}

func (r *RoutingInstanceNode) Plane() graph.Plane {
	return plane
}

func (r *RoutingInstanceNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *RoutingInstanceNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *RoutingInstanceNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
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
			routingInstanceList, err := g.ClientConfig.IntrospectClientV1.Control(podResource.Status.PodIP + ":" + introspectPort).RoutingInstances().List(context.Background())
			if err != nil {
				return nil, err
			}
			for _, resource := range routingInstanceList.Instances.List.ShowRoutingInstance {
				r.Resource = resource
				routingInstanceNameList := strings.Split(resource.Name.Text, ":")
				resourceNode := &RoutingInstanceNode{
					Resource: resource,
					EdgeSelectors: []graph.EdgeSelector{{
						NodeType: graph.RoutingInstance,
						Plane:    graph.ConfigPlane,
						MatchValues: []graph.MatchValue{{
							Value: map[string]string{"RoutingInstanceName": fmt.Sprintf("%s:%s", routingInstanceNameList[1], routingInstanceNameList[3])},
						}},
					}},
					EdgeLabels: []graph.EdgeLabel{{
						Value: map[string]string{"RoutingInstanceName": fmt.Sprintf("%s:%s", routingInstanceNameList[1], routingInstanceNameList[3])},
					}},
				}
				graphNodeList = append(graphNodeList, resourceNode)
			}
		}
	}

	return graphNodeList, nil
}
