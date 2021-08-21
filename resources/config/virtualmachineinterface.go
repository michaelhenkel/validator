package resources

import (
	"context"
	"fmt"
	"reflect"

	"github.com/s3kim2018/validator/graph"
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
	originalresource.getspecandstatusvals("virtualmachineinterface")
	for i := 0; i < len(resourceList.Items); i++ {
		resource := resourceList.Items[i]
		r.Resource = resource
		var edgeSelectorList []graph.EdgeSelector
		hashmap := buildhash(g.ClientConfig, i, "VirtualMachineInterface")
		combinedlist := append(originalresource.References, originalresource.Reference...)
		combinedlist = append(combinedlist, originalresource.Parents...)
		for i := 0; i < len(combinedlist); i++ {
			if references, ok := hashmap[combinedlist[i]]; ok {
				for j := 0; j < len(references); j++ {
					reference := references[j][0]
					var nodetype graph.NodeType
					nodetype = references[j][1].(graph.NodeType)
					primaryname := combinedlist[i]
					name := primaryname + "Name"
					switch thetype := reference.(type) {
					case []v1alpha1.RoutingInstanceReference:
						if thetype != nil {
							for _, routingInstanceReference := range thetype {
								val := reflect.Indirect(reflect.ValueOf(routingInstanceReference))
								_, ok := val.Type().FieldByName("Name")
								if ok {
									edgeSelector := graph.EdgeSelector{
										NodeType: nodetype,
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
							fmt.Println("Thetype resourcereferences is not nil! ")
							for _, resourcereference := range thetype {
								val := reflect.Indirect(reflect.ValueOf(resourcereference))
								_, ok := val.Type().FieldByName("Name")
								if ok {
									edgeSelector := graph.EdgeSelector{
										NodeType: nodetype,
										Plane:    graph.ConfigPlane,
										MatchValues: []graph.MatchValue{{
											Value: map[string]string{name: resourcereference.Name},
										}},
									}
									edgeSelectorList = append(edgeSelectorList, edgeSelector)
								}
							}
						} else {
							fmt.Println("Thetype resourcereferences is nil :( ")
						}
					case *v1alpha1.ResourceReference:
						if thetype != nil {
							val := reflect.Indirect(reflect.ValueOf(thetype))
							_, ok := val.Type().FieldByName("Name")
							if ok {
								fmt.Println("Shouldn't be okay...")
								edgeSelector := graph.EdgeSelector{
									NodeType: nodetype,
									Plane:    graph.ConfigPlane,
									MatchValues: []graph.MatchValue{{
										Value: map[string]string{name: thetype.Name},
									}},
								}
								edgeSelectorList = append(edgeSelectorList, edgeSelector)
							}
						} else {
							fmt.Println("Thetype is nil :( ")
						}
					}

				}

			}
		}
		if _, ok := hashmap["Parent"]; ok {
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
	}
	return graphNodeList, nil
}
