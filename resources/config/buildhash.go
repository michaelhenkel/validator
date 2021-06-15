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
func buildhash(g *graph.Graph) map[string]interface{} {
	resourceList, err := g.ClientConfig.ContrailCoreV1.VirtualMachineInterfaces("").List(context.Background(), metav1.ListOptions{})
	resource := resourceList.Items[0]
	retmap := make(map[string]interface{})

	if err != nil {
		return retmap
	}
	retmap["RoutingInstance"] = resource.Status.RoutingInstanceReferences
	retmap["VirtualMachine"] = resource.Spec.VirtualMachineReferences
	retmap["VirturalNetwork"] = resource.Spec.VirtualNetworkReference
	retmap["BGPRouter"] = resource.Status.BGPRouterReference
	return retmap

}
