package resources

import (
	"context"

	"github.com/michaelhenkel/validator/graph"
	// "k8s.io/api/apiserverinternal/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/**
Builds Hash Tabel that maps the name of the reference to the actual object instance of this type.
Since we are storing it as an interface, we have to check for the reference type when accessing an element.
@param g -> Takes in a Graph
@Return map[string][]interface{}
*/
func buildhash(g *graph.Graph, resourceindex int, instance string) map[string][][]interface{} {
	switch instance {
	case "VirtualMachineInterface":
		resourceList, err := g.ClientConfig.ContrailCoreV1.VirtualMachineInterfaces("").List(context.Background(), metav1.ListOptions{})
		retmap := make(map[string][][]interface{})

		if err != nil {
			return retmap
		}
		resource := resourceList.Items[resourceindex]
		if &resource.Spec.Parent != nil {
			nested := []interface{}{resource.Spec.Parent, graph.VirtualRouter}
			retmap["Parent"] = append(retmap["parent"], nested)
		}
		if &resource.Status.RoutingInstanceReferences != nil {
			nested := []interface{}{resource.Status.RoutingInstanceReferences, graph.RoutingInstance}
			retmap["RoutingInstance"] = append(retmap["RoutingInstance"], nested)
		}
		if &resource.Spec.VirtualMachineReferences != nil {
			nested := []interface{}{resource.Spec.VirtualMachineReferences, graph.VirtualMachine}
			retmap["VirtualMachine"] = append(retmap["VirtualMachine"], nested)
		}
		if &resource.Spec.VirtualNetworkReference != nil {
			nested := []interface{}{resource.Spec.VirtualNetworkReference, graph.VirtualNetwork}
			retmap["VirturalNetwork"] = append(retmap["VirturalNetwork"], nested)

		}
		if &resource.Status.BGPRouterReference != nil {
			nested := []interface{}{resource.Status.BGPRouterReference, graph.BGPRouter}
			retmap["BGPRouter"] = append(retmap["BGPRouter"], nested)
		}

		return retmap
	case "RoutingInstance":
		resourceList, err := g.ClientConfig.ContrailCoreV1.RoutingInstances("").List(context.Background(), metav1.ListOptions{})
		retmap := make(map[string][][]interface{})
		if err != nil {
			return retmap
		}
		resource := resourceList.Items[resourceindex]
		if &resource.Spec.Parent != nil {
			nested := []interface{}{resource.Spec.Parent, graph.VirtualRouter}
			retmap["Parent"] = append(retmap["parent"], nested)
		}
		if &resource.Status.RouteTargetReferences != nil {
			nested := []interface{}{resource.Status.RouteTargetReferences, graph.RoutingInstance}
			retmap["RouteTarget"] = append(retmap["RouteTarget"], nested)
		}

		return retmap
	default:
		return make(map[string][][]interface{})
	}

}
