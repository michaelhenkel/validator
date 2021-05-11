package walker

import "github.com/michaelhenkel/validator/graph"

func podWalker(graphWalker GraphWalker) map[graph.NodeInterface][]graph.NodeInterface {
	graphWalker.Walk(graph.NodeFilterOption{
		NodeType:  graph.VirtualMachine,
		NodePlane: graph.ConfigPlane,
	}).Walk(graph.NodeFilterOption{
		NodeType:  graph.VirtualMachineInterface,
		NodePlane: graph.ConfigPlane,
	}).Walk(graph.NodeFilterOption{
		NodeType:  graph.RoutingInstance,
		NodePlane: graph.ConfigPlane,
		ErrorMsg:  "no routing instance in config",
	}).Walk(graph.NodeFilterOption{
		NodeType:  graph.RoutingInstance,
		NodePlane: graph.ControlPlane,
		ErrorMsg:  "no routing instance in control, check xmpp",
	}).Walk(graph.NodeFilterOption{
		NodeType:  graph.BGPNeighbor,
		NodePlane: graph.ControlPlane,
		ErrorMsg:  "no bgp neighbor for routing instance in control, check xmpp",
	}).Walk(graph.NodeFilterOption{
		NodeType:  graph.VirtualRouter,
		NodePlane: graph.ConfigPlane,
	}).Walk(graph.NodeFilterOption{
		NodeType:     graph.VirtualMachine,
		NodePlane:    graph.ConfigPlane,
		TargetFilter: graph.VirtualMachine,
	})
	return graphWalker.Result
}

func init() {
	walkerMap[graph.Pod] = podWalker
}
