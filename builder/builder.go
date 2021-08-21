package builder

import (
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/s3kim2018/validator/graph"
	"github.com/s3kim2018/validator/k8s/clientset"
	"github.com/s3kim2018/validator/render"
	configresources "github.com/s3kim2018/validator/resources/config"
	controlresources "github.com/s3kim2018/validator/resources/control"
	dataresources "github.com/s3kim2018/validator/resources/data"
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
	k8snode := configresources.K8SNodeNode{}
	vRouter := configresources.VrouterNode{}
	routingInstanceConfig := configresources.RoutingInstanceNode{}
	routingInstanceControl := controlresources.RoutingInstanceNode{}
	kubemanager := configresources.KubemanagerNode{}
	virtualMachineInterfaceConfig := configresources.VirtualMachineInterfaceNode{}
	virtualMachineConfig := configresources.VirtualMachineNode{}
	virtualNetworkConfig := configresources.VirtualNetworkNode{}
	virtualNetworkData := dataresources.VirtualNetworkNode{}
	routingInstanceData := dataresources.RoutingInstanceNode{}
	virtualMachineInterfaceData := dataresources.VirtualMachineInterfaceNode{}

	g.NodeAdder(virtualRouter.AdderFunc()).
		NodeAdder(bgpNeighbor.AdderFunc()).
		NodeAdder(bgpRouter.AdderFunc()).
		NodeAdder(control.AdderFunc()).
		NodeAdder(configFile.AdderFunc()).
		NodeAdder(configMap.AdderFunc()).
		NodeAdder(pod.AdderFunc()).
		NodeAdder(k8snode.AdderFunc()).
		NodeAdder(vRouter.AdderFunc()).
		NodeAdder(routingInstanceConfig.AdderFunc()).
		NodeAdder(virtualMachineInterfaceConfig.AdderFunc()).
		NodeAdder(virtualMachineConfig.AdderFunc()).
		NodeAdder(virtualNetworkConfig.AdderFunc()).
		NodeAdder(routingInstanceControl.AdderFunc()).
		NodeAdder(kubemanager.AdderFunc()).
		NodeAdder(virtualNetworkData.AdderFunc()).
		NodeAdder(routingInstanceData.AdderFunc()).
		NodeAdder(virtualMachineInterfaceData.AdderFunc()).
		EdgeMatcher()
	return g
}

func RenderPage(nodeEdges map[graph.NodeInterface][]graph.NodeInterface) *components.Page {
	return render.RenderPage(nodeEdges)
}
