package resources

import (
	"context"
	"fmt"
	"reflect"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
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
	var originalresource sourcecoderesource
	//originalresource.getspecandstatusvals("routinginstance")
	for i := 0; i < len(resourceList.Items); i++ {
		resource := resourceList.Items[i]
		r.Resource = resource
		var edgeSelectorList []graph.EdgeSelector
		hashmap := buildhash(g.ClientConfig, i, "RoutingInstance")
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
					default:
						fmt.Println("Unexpected Type")
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
						}
					case *v1alpha1.ResourceReference:
						if thetype != nil {
							val := reflect.Indirect(reflect.ValueOf(thetype))
							_, ok := val.Type().FieldByName("Name")
							if ok {
								edgeSelector := graph.EdgeSelector{
									NodeType: nodetype,
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
		}
		if _, ok := hashmap["Parent"]; ok {
			fmt.Println("Parent Found for Routing Instance! Name is: ", resource.Spec.Parent.Name)
			edgeSelector := graph.EdgeSelector{
				NodeType: graph.VirtualNetwork,
				Plane:    graph.ConfigPlane,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"VirtualNetworkNamespaceName": fmt.Sprintf("%s/%s", resource.Spec.Parent.Namespace, resource.Spec.Parent.Name)},
				}},
			}
			edgeSelectorList = append(edgeSelectorList, edgeSelector)
		}
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
