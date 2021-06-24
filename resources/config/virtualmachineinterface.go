package resources

import (
	"context"
	"fmt"
	"reflect"

	"github.com/michaelhenkel/validator/graph"
	"ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"

	//intserver "k8s.io/api/apiserverinternal/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type VirtualMachineInterfaceNode struct {
	Resource      contrailcorev1alpha1.VirtualMachineInterface
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
	return fmt.Sprintf("%s:%s", r.Resource.Namespace, r.Resource.Name)
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
	resourceList, err := g.ClientConfig.ContrailCoreV1.VirtualMachineInterfaces("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var originalresource sourcecoderesource
	originalresource.getspecvals("VirtualMachineInterface")
	originalresource.getstatusvals("VirtualMachineInterface")
	fmt.Println("Length of resourcelist of virtualmachineinterface", len(resourceList.Items))

	resource := resourceList.Items[0]
	var edgeSelectorList []graph.EdgeSelector

	r.Resource = resource
	hashmap := buildhash(g, "VirtualMachineInterface")
	combinedlist := append(originalresource.References, originalresource.Reference...)
	for i := 0; i < len(combinedlist); i++ {
		if references, ok := hashmap[combinedlist[i]]; ok {
			primaryname := combinedlist[i]
			name := primaryname + "Name"
			switch thetype := references.(type) {
			default:
				fmt.Println("Unexpected Type")
			case []v1alpha1.RoutingInstanceReference:
				if thetype != nil {
					for _, routingInstanceReference := range thetype {
						val := reflect.Indirect(reflect.ValueOf(routingInstanceReference))
						_, ok := val.Type().FieldByName("Name")
						if ok {
							//fmt.Println("References Found for "+primaryname+" Name is: ", routingInstanceReference.Name)
							edgeSelector := graph.EdgeSelector{
								NodeType: graph.RoutingInstance,
								Plane:    graph.ConfigPlane,
								MatchValues: []graph.MatchValue{{
									Value: map[string]string{name: routingInstanceReference.Name},
								}},
							}
							edgeSelectorList = append(edgeSelectorList, edgeSelector)
						}
					}
				}
			case []v1alpha1.ResourceReference:
				if thetype != nil {
					for _, resourcereference := range thetype {
						val := reflect.Indirect(reflect.ValueOf(resourcereference))
						_, ok := val.Type().FieldByName("Name")
						if ok {
							//fmt.Println("References Found for "+primaryname+" Name is: ", resourcereference.Name)
							edgeSelector := graph.EdgeSelector{
								NodeType: graph.RoutingInstance,
								Plane:    graph.ConfigPlane,
								MatchValues: []graph.MatchValue{{
									Value: map[string]string{name: resourcereference.Name},
								}},
							}
							edgeSelectorList = append(edgeSelectorList, edgeSelector)
						}
					}
				}
			case *v1alpha1.ResourceReference:
				if thetype != nil {
					val := reflect.Indirect(reflect.ValueOf(thetype))
					_, ok := val.Type().FieldByName("Name")
					if ok {
						//fmt.Println("References Found for "+primaryname+" Name is: ", thetype.Name)
						edgeSelector := graph.EdgeSelector{
							NodeType: graph.RoutingInstance,
							Plane:    graph.ConfigPlane,
							MatchValues: []graph.MatchValue{{
								Value: map[string]string{name: thetype.Name},
							}},
						}
						edgeSelectorList = append(edgeSelectorList, edgeSelector)
					}
				}
			}
		}
	}
	if len(originalresource.Parents) > 0 {
		edgeSelector := graph.EdgeSelector{
			NodeType: graph.VirtualRouter,
			Plane:    graph.ConfigPlane,
			MatchValues: []graph.MatchValue{{
				Value: map[string]string{"VirtualRouterName": resource.Spec.Parent.Name},
			}},
		}
		edgeSelectorList = append(edgeSelectorList, edgeSelector)
	}

	edgeSelector := graph.EdgeSelector{
		NodeType: graph.VirtualNetwork,
		Plane:    graph.ConfigPlane,
		MatchValues: []graph.MatchValue{{
			Value: map[string]string{"VirtualNetworkNamespaceName": fmt.Sprintf("%s/%s", resource.Spec.VirtualNetworkReference.Namespace, resource.Spec.VirtualNetworkReference.Name)},
		}},
	}
	edgeSelectorList = append(edgeSelectorList, edgeSelector)
	resourceNode := &VirtualMachineInterfaceNode{
		Resource: resource,
		EdgeLabels: []graph.EdgeLabel{{
			Value: map[string]string{"VirtualMachineInterfaceName": resource.Name},
		}},
		EdgeSelectors: edgeSelectorList,
	}
	graphNodeList = append(graphNodeList, resourceNode)

	return graphNodeList, nil
}
