package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/s3kim2018/validator/graph"

	introspectdatav1alpha1 "github.com/michaelhenkel/introspect/pkg/apis/data/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VirtualNetworkNode struct {
	Resource      introspectdatav1alpha1.VnSandeshData
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *VirtualNetworkNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*VirtualNetworkNode)
	if !ok {
		return fmt.Errorf("not a VirtualNetworkNode resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *VirtualNetworkNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *VirtualNetworkNode) Name() string {
	name := r.Resource.Name.Text
	nameList := strings.Split(r.Resource.Name.Text, ":")
	if len(nameList) > 1 {
		name = fmt.Sprintf("%s:%s", nameList[1], nameList[2])
	}
	return name
}

func (r *VirtualNetworkNode) Type() graph.NodeType {
	return graph.VirtualNetwork
}

func (r *VirtualNetworkNode) Plane() graph.Plane {
	return plane
}

func (r *VirtualNetworkNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *VirtualNetworkNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *VirtualNetworkNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
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
			virtualNetworkList, err := g.ClientConfig.IntrospectClientV1.Data(podResource.Status.PodIP + ":" + introspectPort).VirtualNetworks().List(context.Background())
			if err != nil {
				return nil, err
			}
			for _, resource := range virtualNetworkList.VnListResp.VnList.List.VnSandeshData {
				r.Resource = resource
				virtualNetworkNameList := strings.Split(resource.Name.Text, ":")
				resourceNode := &VirtualNetworkNode{
					Resource: resource,
					EdgeSelectors: []graph.EdgeSelector{{
						NodeType: graph.VirtualNetwork,
						Plane:    graph.ConfigPlane,
						MatchValues: []graph.MatchValue{{
							Value: map[string]string{"VirtualNetworkNamespaceName": fmt.Sprintf("%s/%s", virtualNetworkNameList[1], virtualNetworkNameList[2])},
						}},
					}},
					EdgeLabels: []graph.EdgeLabel{{
						Value: map[string]string{"VirtualNetworkNamespaceName": fmt.Sprintf("%s/%s", virtualNetworkNameList[1], virtualNetworkNameList[2])},
					}},
				}
				graphNodeList = append(graphNodeList, resourceNode)
			}
		}
	}

	return graphNodeList, nil
}
