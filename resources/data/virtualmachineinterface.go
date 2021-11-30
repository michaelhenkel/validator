package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/s3kim2018/validator/graph"

	introspectdatav1alpha1 "github.com/michaelhenkel/introspect/pkg/apis/data/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VirtualMachineInterfaceNode struct {
	Resource      introspectdatav1alpha1.ItfSandeshData
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *VirtualMachineInterfaceNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*VirtualMachineInterfaceNode)
	if !ok {
		return fmt.Errorf("not a VirtualMachineInterfaceNode resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *VirtualMachineInterfaceNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *VirtualMachineInterfaceNode) Name() string {
	name := r.Resource.ConfigName.Text
	nameList := strings.Split(r.Resource.ConfigName.Text, ":")
	if len(nameList) > 1 {
		name = fmt.Sprintf("%s:%s", nameList[1], nameList[2])
	}
	return name
}

func (r *VirtualMachineInterfaceNode) Type() graph.NodeType {
	return graph.VirtualMachineInterface
}

func (r *VirtualMachineInterfaceNode) Plane() graph.Plane {
	return plane
}

func (r *VirtualMachineInterfaceNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *VirtualMachineInterfaceNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *VirtualMachineInterfaceNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
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
			virtualMachineInterfaceList, err := g.ClientConfig.IntrospectClientV1.Data(podResource.Status.PodIP + ":" + introspectPort).VirtualMachineInterfaces().List(context.Background())
			if err != nil {
				return nil, err
			}
			for _, resource := range virtualMachineInterfaceList.ItfResp.ItfList.List.ItfSandeshData {
				r.Resource = resource
				virtualMachineInterfaceNameList := strings.Split(resource.ConfigName.Text, ":")
				virtualNetworkNameList := strings.Split(resource.VnName.Text, ":")
				var edgeLabelList []graph.EdgeLabel
				if len(virtualMachineInterfaceNameList) > 1 {
					edgeLabel := graph.EdgeLabel{
						Value: map[string]string{"VirtualMachineInterfaceName": fmt.Sprintf("%s/%s", virtualMachineInterfaceNameList[1], virtualMachineInterfaceNameList[2])},
					}
					edgeLabelList = append(edgeLabelList, edgeLabel)
				}
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
				resourceNode := &VirtualMachineInterfaceNode{
					Resource:      resource,
					EdgeSelectors: edgeSelectorList,
					EdgeLabels:    edgeLabelList,
				}
				graphNodeList = append(graphNodeList, resourceNode)
			}
		}
	}

	return graphNodeList, nil
}
