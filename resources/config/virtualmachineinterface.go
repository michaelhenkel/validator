package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

var originalresource sourcecoderesource

func (r *VirtualMachineInterfaceNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	resourceList, err := g.ClientConfig.ContrailCoreV1.VirtualMachineInterfaces("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	fmt.Println("start")
	originalresource.getspecvals()
	originalresource.getstatusvals()
	fmt.Println("References ******")
	for i := 0; i < len(originalresource.References); i++ {
		fmt.Println(originalresource.References[i])
	}
	fmt.Println("Reference *******")
	for i := 0; i < len(originalresource.Reference); i++ {
		fmt.Println(originalresource.References[i])
	}
	fmt.Println("end")
	for _, resource := range resourceList.Items {
		r.Resource = resource
		var edgeSelectorList []graph.EdgeSelector
		if resource.Spec.Parent.Kind == "VirtualRouter" {
			fmt.Println("Parent Found for Virtural Router! Name is: ", resource.Spec.Parent.Name)
			edgeSelector := graph.EdgeSelector{
				NodeType: graph.VirtualRouter,
				Plane:    graph.ConfigPlane,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"VirtualRouterName": resource.Spec.Parent.Name},
				}},
			}
			edgeSelectorList = append(edgeSelectorList, edgeSelector)
		}
		for _, routingInstanceReference := range resource.Status.RoutingInstanceReferences {
			fmt.Println("Parent Found for RoutingInstanceReferences! Name is: ", routingInstanceReference.Name)
			edgeSelector := graph.EdgeSelector{
				NodeType: graph.RoutingInstance,
				Plane:    graph.ConfigPlane,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"RoutingInstanceName": routingInstanceReference.Name},
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

		for _, virtualMachineReference := range resource.Spec.VirtualMachineReferences {
			edgeSelector := graph.EdgeSelector{
				NodeType: graph.VirtualMachine,
				Plane:    graph.ConfigPlane,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"VirtualMachineName": virtualMachineReference.Name},
				}},
			}
			edgeSelectorList = append(edgeSelectorList, edgeSelector)
		}
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
