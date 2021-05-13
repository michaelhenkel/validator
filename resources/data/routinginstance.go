package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaelhenkel/validator/graph"

	introspectdatav1alpha1 "github.com/michaelhenkel/introspect/pkg/apis/data/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RoutingInstanceNode struct {
	Resource      introspectdatav1alpha1.VrfSandeshData
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
	name := r.Resource.Name.Text
	nameList := strings.Split(r.Resource.Name.Text, ":")
	if len(nameList) > 2 {
		name = fmt.Sprintf("%s:%s", nameList[1], nameList[3])
	}
	return name
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

	vrouterList, err := g.ClientConfig.DeployerDataV1.Vrouters("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, vrouterResource := range vrouterList.Items {
		filterMap[vrouterResource.Name] = struct{}{}
	}

	for vrouterFilter := range filterMap {
		opts := metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", vrouterFilter),
		}
		podList, err := g.ClientConfig.CoreV1.Pods("").List(context.Background(), opts)
		if err != nil {
			return nil, err
		}
		for _, podResource := range podList.Items {
			routingInstanceList, err := g.ClientConfig.IntrospectClientV1.Data(podResource.Status.PodIP + ":" + introspectPort).RoutingInstances().List(context.Background())
			if err != nil {
				return nil, err
			}
			for _, resource := range routingInstanceList.VrfListResp.VrfList.List.VrfSandeshData {
				r.Resource = resource
				routingInstanceNameList := strings.Split(resource.Name.Text, ":")
				virtualNetworkNameList := strings.Split(resource.Vn.Text, ":")
				var edgeSelectorList []graph.EdgeSelector
				if len(virtualNetworkNameList) > 1 {
					edgeSelector := graph.EdgeSelector{
						NodeType: graph.VirtualNetwork,
						Plane:    graph.DataPlane,
						MatchValues: []graph.MatchValue{{
							Value: map[string]string{"VirtualNetworkNamespaceName": fmt.Sprintf("%s/%s", virtualNetworkNameList[1], virtualNetworkNameList[2])},
						}},
					}
					edgeSelectorList = append(edgeSelectorList, edgeSelector)
				}
				resourceNode := &RoutingInstanceNode{
					Resource: resource,
					EdgeLabels: []graph.EdgeLabel{{
						Value: map[string]string{"RoutingInstanceName": fmt.Sprintf("%s/%s", routingInstanceNameList[1], routingInstanceNameList[3])},
					}},
					EdgeSelectors: edgeSelectorList,
				}
				graphNodeList = append(graphNodeList, resourceNode)
			}
		}
	}

	return graphNodeList, nil
}
