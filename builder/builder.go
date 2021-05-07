package builder

import (
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/michaelhenkel/validator/graph"
	"github.com/michaelhenkel/validator/k8s/clientset"
	configresources "github.com/michaelhenkel/validator/resources/config"
	controlresources "github.com/michaelhenkel/validator/resources/control"
)

func BuildGraph(clientConfig *clientset.Client) *graph.Graph {

	g := graph.NewGraph(clientConfig)
	virtualRouter := configresources.VirtualRouterNode{}
	bgpNeighbor := controlresources.BGPNeighborNode{}
	bgpRouter := configresources.BGPRouterNode{}
	control := configresources.ControlNode{}
	configFile := configresources.ConfigFileNode{}
	configMap := configresources.ConfigMapNode{}
	pod := configresources.PodNode{}
	vRouter := configresources.VrouterNode{}
	routingInstanceConfig := configresources.RoutingInstanceNode{}
	routingInstanceControl := controlresources.RoutingInstanceNode{}
	virtualMachineInterfaceConfig := configresources.VirtualMachineInterfaceNode{}

	g.NodeAdder(virtualRouter.AdderFunc()).
		NodeAdder(bgpNeighbor.AdderFunc()).
		NodeAdder(bgpRouter.AdderFunc()).
		NodeAdder(control.AdderFunc()).
		NodeAdder(configFile.AdderFunc()).
		NodeAdder(configMap.AdderFunc()).
		NodeAdder(pod.AdderFunc()).
		NodeAdder(vRouter.AdderFunc()).
		NodeAdder(routingInstanceConfig.AdderFunc()).
		NodeAdder(virtualMachineInterfaceConfig.AdderFunc()).
		NodeAdder(routingInstanceControl.AdderFunc()).
		EdgeMatcher()
	return g
}

func RenderPage(g *graph.Graph) *components.Page {
	return g.RenderPage()
}
