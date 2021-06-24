package resources

import (
	"context"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/**
Builds Hash Tabel that maps the name of the reference to the actual object instance of this type.
Since we are storing it as an interface, we have to check for the reference type when accessing an element.
@param g -> Takes in a Graph
@Return map[string]interface{}
*/
func buildhash(g *graph.Graph, instance string) map[string]interface{} {
	switch instance {
	case "VirtualMachineInterface":
		resourceList, err := g.ClientConfig.ContrailCoreV1.VirtualMachineInterfaces("").List(context.Background(), metav1.ListOptions{})
		retmap := make(map[string]interface{})

		if err != nil {
			return retmap
		}
		for _, resource := range resourceList.Items {
			if &resource.Spec.Parent != nil {
				retmap["Parent"] = resource.Spec.Parent
			}
			if &resource.Status.RoutingInstanceReferences != nil {
				retmap["RoutingInstance"] = resource.Status.RoutingInstanceReferences
			}
			if &resource.Spec.VirtualMachineReferences != nil {
				retmap["VirtualMachine"] = resource.Spec.VirtualMachineReferences
			}
			if &resource.Spec.VirtualNetworkReference != nil {
				retmap["VirturalNetwork"] = resource.Spec.VirtualNetworkReference
			}
			if &resource.Status.BGPRouterReference != nil {
				retmap["BGPRouter"] = resource.Status.BGPRouterReference

			}
		}

		return retmap
	case "RoutingInstance":
		resourceList, err := g.ClientConfig.ContrailCoreV1.RoutingInstances("").List(context.Background(), metav1.ListOptions{})
		retmap := make(map[string]interface{})
		if err != nil {
			return retmap
		}
		for _, resource := range resourceList.Items {
			if &resource.Spec.Parent != nil {
				retmap["Parent"] = resource.Spec.Parent
			}
			if &resource.Status.RouteTargetReferences != nil {
				retmap["RouteTarget"] = resource.Status.RouteTargetReferences
			}
		}
		return retmap
	default:
		return make(map[string]interface{})
	}

}
