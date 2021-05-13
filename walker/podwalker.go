package walker

import "github.com/michaelhenkel/validator/graph"

func podToVrouter(g *graph.Graph, nodeType graph.NodeType, plane graph.Plane, name string) map[graph.NodeInterface][]graph.NodeInterface {
	nodeInterface := g.GetNodeByTypePlaneName(nodeType, plane, name)
	var sourceNodeInterfaceList []graph.NodeInterface
	sourceNodeInterfaceList = append(sourceNodeInterfaceList, nodeInterface)
	graphWalker := GraphWalker{
		G:           g,
		SourceNodes: sourceNodeInterfaceList,
	}
	w := Walker{
		FilterOpts: []graph.NodeFilterOption{{
			NodeType:  graph.VirtualMachine,
			NodePlane: graph.ConfigPlane,
		}},
		WalkerFunc: graphWalker.Walk2,
		Next: []Walker{{
			FilterOpts: []graph.NodeFilterOption{{
				NodeType:  graph.VirtualMachineInterface,
				NodePlane: graph.ConfigPlane,
			}},
			WalkerFunc: graphWalker.Walk2,
			Next: []Walker{{
				FilterOpts: []graph.NodeFilterOption{{
					NodeType:  graph.RoutingInstance,
					NodePlane: graph.ConfigPlane,
				}},
				WalkerFunc: graphWalker.Walk2,
				Next: []Walker{{
					FilterOpts: []graph.NodeFilterOption{{
						NodeType:  graph.VirtualNetwork,
						NodePlane: graph.ConfigPlane,
					}},
					WalkerFunc: graphWalker.Walk2,
					Next: []Walker{{
						FilterOpts: []graph.NodeFilterOption{{
							NodeType:     graph.Pod,
							NodePlane:    graph.ConfigPlane,
							TargetFilter: graph.Pod,
						}},
						WalkerFunc: graphWalker.Walk2,
					}, {
						FilterOpts: []graph.NodeFilterOption{{
							NodeType:     graph.VirtualMachineInterface,
							NodePlane:    graph.ConfigPlane,
							TargetFilter: graph.VirtualMachineInterface,
						}},
						WalkerFunc: graphWalker.Walk2,
					}},
				}, {
					FilterOpts: []graph.NodeFilterOption{{
						NodeType:  graph.RoutingInstance,
						NodePlane: graph.ControlPlane,
					}},
					WalkerFunc: graphWalker.Walk2,
					Next: []Walker{{
						FilterOpts: []graph.NodeFilterOption{{
							NodeType:  graph.BGPNeighbor,
							NodePlane: graph.ControlPlane,
							ErrorMsg:  "no bgp neighbor for routing instance in control, check xmpp",
						}},
						WalkerFunc: graphWalker.Walk2,
						Next: []Walker{{
							FilterOpts: []graph.NodeFilterOption{{
								NodeType:     graph.RoutingInstance,
								NodePlane:    graph.DataPlane,
								TargetFilter: graph.RoutingInstance,
							}},
							WalkerFunc: graphWalker.Walk2,
							Next: []Walker{{
								FilterOpts: []graph.NodeFilterOption{{
									NodeType:     graph.VirtualNetwork,
									NodePlane:    graph.DataPlane,
									TargetFilter: graph.VirtualNetwork,
								}},
								WalkerFunc: graphWalker.Walk2,
								Next: []Walker{{
									FilterOpts: []graph.NodeFilterOption{{
										NodeType:     graph.VirtualMachineInterface,
										NodePlane:    graph.DataPlane,
										TargetFilter: graph.VirtualMachineInterface,
									}},
									WalkerFunc: graphWalker.Walk2,
								}},
							}},
						}, {
							FilterOpts: []graph.NodeFilterOption{{
								NodeType:  graph.VirtualRouter,
								NodePlane: graph.ConfigPlane,
							}},
							WalkerFunc: graphWalker.Walk2,
							Next: []Walker{{
								FilterOpts: []graph.NodeFilterOption{{
									NodeType:     graph.VirtualMachine,
									NodePlane:    graph.ConfigPlane,
									TargetFilter: graph.VirtualMachine,
								}},
								WalkerFunc: graphWalker.Walk2,
							}},
						}},
					}},
				}},
			}},
		}},
	}
	w.runner(sourceNodeInterfaceList)
	return graphWalker.Result
}

func init() {
	walker2Map[graph.Pod] = podToVrouter
}
