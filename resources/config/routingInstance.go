package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type RoutingInstanceNode struct {
	Resource      contrailcorev1alpha1.RoutingInstance
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
	return fmt.Sprintf("%s:%s", r.Resource.Namespace, r.Resource.Name)
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
	resourceList, err := g.ClientConfig.ContrailCoreV1.RoutingInstances("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// var originalresource sourcecoderesource
	// originalresource.getspecvals("RoutingInstance")
	// originalresource.getstatusvals("RoutingInstance")
	// hashmap := buildhash(g, "RoutingInstance")
	// combinedlist := append(originalresource.References, originalresource.Reference...)
	var edgeSelectorList []graph.EdgeSelector
	// for i := 0; i < len(combinedlist); i++ {
	// 	if references, ok := hashmap[combinedlist[i]]; ok {
	// 		primaryname := combinedlist[i]
	// 		name := primaryname + "Name"
	// 		switch thetype := references.(type) {
	// 		default:
	// 			fmt.Println("Unexpected Type")
	// 		case []v1alpha1.RoutingInstanceReference:
	// 			if thetype != nil {
	// 				for _, routingInstanceReference := range thetype {
	// 					val := reflect.Indirect(reflect.ValueOf(routingInstanceReference))
	// 					_, ok := val.Type().FieldByName("Name")
	// 					if ok {
	// 						//fmt.Println("References Found for "+primaryname+" Name is: ", routingInstanceReference.Name)
	// 						edgeSelector := graph.EdgeSelector{
	// 							NodeType: graph.RoutingInstance,
	// 							Plane:    graph.ConfigPlane,
	// 							MatchValues: []graph.MatchValue{{
	// 								Value: map[string]string{name: routingInstanceReference.Name},
	// 							}},
	// 						}
	// 						edgeSelectorList = append(edgeSelectorList, edgeSelector)
	// 					}
	// 				}
	// 			}
	// 		case []v1alpha1.ResourceReference:
	// 			if thetype != nil {
	// 				for _, resourcereference := range thetype {
	// 					val := reflect.Indirect(reflect.ValueOf(resourcereference))
	// 					_, ok := val.Type().FieldByName("Name")
	// 					if ok {
	// 						//fmt.Println("References Found for "+primaryname+" Name is: ", resourcereference.Name)
	// 						edgeSelector := graph.EdgeSelector{
	// 							NodeType: graph.RoutingInstance,
	// 							Plane:    graph.ConfigPlane,
	// 							MatchValues: []graph.MatchValue{{
	// 								Value: map[string]string{name: resourcereference.Name},
	// 							}},
	// 						}
	// 						edgeSelectorList = append(edgeSelectorList, edgeSelector)
	// 					}
	// 				}
	// 			}

	// 		case *v1alpha1.ResourceReference:
	// 			if thetype != nil {
	// 				val := reflect.Indirect(reflect.ValueOf(thetype))
	// 				_, ok := val.Type().FieldByName("Name")
	// 				if ok {
	// 					//fmt.Println("References Found for "+primaryname+" Name is: ", thetype.Name)
	// 					edgeSelector := graph.EdgeSelector{
	// 						NodeType: graph.RoutingInstance,
	// 						Plane:    graph.ConfigPlane,
	// 						MatchValues: []graph.MatchValue{{
	// 							Value: map[string]string{name: thetype.Name},
	// 						}},
	// 					}
	// 					edgeSelectorList = append(edgeSelectorList, edgeSelector)
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	for _, resource := range resourceList.Items {
		r.Resource = resource
		if resource.Spec.Parent.Kind == "VirtualNetwork" {
			fmt.Println("Found Virtual Network", resource.Spec.Parent.Namespace, resource.Spec.Parent.Name)
			edgeSelector := graph.EdgeSelector{
				NodeType: graph.VirtualNetwork,
				Plane:    graph.ConfigPlane,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"VirtualNetworkNamespaceName": fmt.Sprintf("%s/%s", resource.Spec.Parent.Namespace, resource.Spec.Parent.Name)},
				}},
			}
			edgeSelectorList = append(edgeSelectorList, edgeSelector)
		}
		fmt.Println("Found Routing Instance", resource.Namespace, resource.Name)

		edgeSelector := graph.EdgeSelector{
			NodeType: graph.RoutingInstance,
			Plane:    graph.ControlPlane,
			MatchValues: []graph.MatchValue{{
				Value: map[string]string{"RoutingInstanceName": fmt.Sprintf("%s/%s", resource.Namespace, resource.Name)},
			}},
		}
		edgeSelectorList = append(edgeSelectorList, edgeSelector)
		resourceNode := &RoutingInstanceNode{
			Resource: resource,
			EdgeLabels: []graph.EdgeLabel{{
				Value: map[string]string{"RoutingInstanceName": resource.Name},
			}},
			EdgeSelectors: edgeSelectorList,
		}
		graphNodeList = append(graphNodeList, resourceNode)
	}
	return graphNodeList, nil
}
