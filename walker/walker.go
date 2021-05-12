package walker

import (
	"fmt"

	"github.com/michaelhenkel/validator/builder"
	"github.com/michaelhenkel/validator/graph"
	"github.com/michaelhenkel/validator/k8s/clientset"
)

type GraphWalker struct {
	SourceNodes []graph.NodeInterface
	G           *graph.Graph
	Tracker     map[graph.NodeInterface]struct{}
	Result      map[graph.NodeInterface][]graph.NodeInterface
}

func (gw *GraphWalker) Walk(filterOpts graph.NodeFilterOption) *GraphWalker {
	var nodeInterfaceList []graph.NodeInterface
	if gw.Tracker == nil {
		gw.Tracker = make(map[graph.NodeInterface]struct{})
	}
	if gw.Result == nil {
		gw.Result = make(map[graph.NodeInterface][]graph.NodeInterface)
	}
	for _, sourceNode := range gw.SourceNodes {
		fmt.Printf("source %s:%s:%s\n", sourceNode.Plane(), sourceNode.Type(), sourceNode.Name())
		gw.Tracker[sourceNode] = struct{}{}
		nodeInterfaceList = append(nodeInterfaceList, gw.G.GetNodeEdge(sourceNode, filterOpts)...)
		targetNodes := gw.G.GetNodeEdge(sourceNode, filterOpts)
		if len(targetNodes) == 0 && filterOpts.ErrorMsg != "" {
			errNode := &graph.ErrorNode{
				NodeName:   fmt.Sprintf("%s\n%s", filterOpts.ErrorMsg, sourceNode.Name()),
				ErrorPlane: sourceNode.Plane(),
			}
			if _, ok := gw.Result[errNode]; ok {
				gw.Result[errNode] = append(gw.Result[errNode], sourceNode)
			} else {
				gw.Result[errNode] = []graph.NodeInterface{sourceNode}
			}
			gw.Result[sourceNode] = append(gw.Result[sourceNode], errNode)
			fmt.Printf("--> %s:%s:%s\n", errNode.Plane(), errNode.Type(), errNode.Name())
		} else {
			for _, node := range targetNodes {
				if filterOpts.TargetFilter != "" {
					for targetNode := range gw.Tracker {
						if targetNode.Type() == filterOpts.TargetFilter && targetNode.Name() == node.Name() {
							fmt.Printf("--> %s:%s:%s\n", node.Plane(), node.Type(), node.Name())
							gw.Result[sourceNode] = append(gw.Result[sourceNode], targetNode)
						}
					}
				} else {
					fmt.Printf("--> %s:%s:%s\n", node.Plane(), node.Type(), node.Name())
					gw.Result[sourceNode] = append(gw.Result[sourceNode], node)
				}
			}
		}
	}
	gw.SourceNodes = nodeInterfaceList
	return gw
}

var walkerMap = make(map[graph.NodeType]func(GraphWalker) map[graph.NodeInterface][]graph.NodeInterface)

func Walk(client *clientset.Client, nodeType graph.NodeType, plane graph.Plane, name string) map[graph.NodeInterface][]graph.NodeInterface {
	g := builder.BuildGraph(client)
	g.EdgeMatcher()
	nodeInterface := g.GetNodeByTypePlaneName(nodeType, plane, name)
	var sourceNodeInterfaceList []graph.NodeInterface
	sourceNodeInterfaceList = append(sourceNodeInterfaceList, nodeInterface)
	graphWalker := GraphWalker{
		G:           g,
		SourceNodes: sourceNodeInterfaceList,
	}
	return walkerMap[nodeType](graphWalker)
}

type Walker struct {
	Next       []Walker
	FilterOpts graph.NodeFilterOption
	WalkerFunc func(filterOpts graph.NodeFilterOption) *GraphWalker
}

func WalkTest(client *clientset.Client, nodeType graph.NodeType, plane graph.Plane, name string) map[graph.NodeInterface][]graph.NodeInterface {
	g := builder.BuildGraph(client)
	g.EdgeMatcher()
	nodeInterface := g.GetNodeByTypePlaneName(nodeType, plane, name)
	var sourceNodeInterfaceList []graph.NodeInterface
	sourceNodeInterfaceList = append(sourceNodeInterfaceList, nodeInterface)
	graphWalker := GraphWalker{
		G:           g,
		SourceNodes: sourceNodeInterfaceList,
	}
	w := Walker{
		FilterOpts: graph.NodeFilterOption{
			NodeType:  graph.VirtualMachine,
			NodePlane: graph.ConfigPlane,
		},
		WalkerFunc: graphWalker.Walk,
		Next: []Walker{{
			FilterOpts: graph.NodeFilterOption{
				NodeType:  graph.VirtualMachineInterface,
				NodePlane: graph.ConfigPlane,
			},
			WalkerFunc: graphWalker.Walk,
			Next: []Walker{{
				FilterOpts: graph.NodeFilterOption{
					NodeType:  graph.RoutingInstance,
					NodePlane: graph.ConfigPlane,
				},
				WalkerFunc: graphWalker.Walk,
				Next: []Walker{{
					FilterOpts: graph.NodeFilterOption{
						NodeType:  graph.RoutingInstance,
						NodePlane: graph.ControlPlane,
					},
					WalkerFunc: graphWalker.Walk,
					Next: []Walker{{
						FilterOpts: graph.NodeFilterOption{
							NodeType:  graph.BGPNeighbor,
							NodePlane: graph.ControlPlane,
							ErrorMsg:  "no bgp neighbor for routing instance in control, check xmpp",
						},
						WalkerFunc: graphWalker.Walk,
						Next: []Walker{{
							FilterOpts: graph.NodeFilterOption{
								NodeType:  graph.VirtualRouter,
								NodePlane: graph.ConfigPlane,
							},
							WalkerFunc: graphWalker.Walk,
							Next: []Walker{{
								FilterOpts: graph.NodeFilterOption{
									NodeType:     graph.VirtualMachine,
									NodePlane:    graph.ConfigPlane,
									TargetFilter: graph.VirtualMachine,
								},
								WalkerFunc: graphWalker.Walk,
							}},
						}},
					}},
				}},
			}},
		}, {
			FilterOpts: graph.NodeFilterOption{
				NodeType:  graph.VirtualNetwork,
				NodePlane: graph.ConfigPlane,
			},
			WalkerFunc: graphWalker.Walk,
		}},
	}
	w.runner()
	return graphWalker.Result
}

func (w Walker) runner() {
	w.WalkerFunc(w.FilterOpts)
	if len(w.Next) > 0 {
		if len(w.Next) == 2 {
			fmt.Println("2")
		}
		for _, nextWalker := range w.Next {
			fmt.Println(w.FilterOpts.NodeType)
			if w.FilterOpts.NodeType == graph.VirtualNetwork {
				fmt.Println("vb")
			}
			nextWalker.runner()
		}
	}
}
