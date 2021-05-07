package builder

import (
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/michaelhenkel/validator/graph"
	"github.com/michaelhenkel/validator/k8s/clientset"
	"github.com/michaelhenkel/validator/resources"
)

func BuildGraph(clientConfig *clientset.Client) *graph.Graph {

	g := graph.NewGraph(clientConfig)
	virtualRouter := resources.VirtualRouterNode{}
	bgpNeighbor := resources.BGPNeighborNode{}
	bgpRouter := resources.BGPRouterNode{}
	control := resources.ControlNode{}
	configFile := resources.ConfigFileNode{}
	configMap := resources.ConfigMapNode{}
	pod := resources.PodNode{}
	vRouter := resources.VrouterNode{}

	g.NodeAdder(virtualRouter.AdderFunc()).
		NodeAdder(bgpNeighbor.AdderFunc()).
		NodeAdder(bgpRouter.AdderFunc()).
		NodeAdder(control.AdderFunc()).
		NodeAdder(configFile.AdderFunc()).
		NodeAdder(configMap.AdderFunc()).
		NodeAdder(pod.AdderFunc()).
		NodeAdder(vRouter.AdderFunc()).
		EdgeMatcher()
	return g
}

func RenderPage(g *graph.Graph) *components.Page {
	return g.RenderPage()
}
