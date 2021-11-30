package walker

import (
	"fmt"
	"strings"

	"github.com/s3kim2018/validator/builder"
	"github.com/s3kim2018/validator/graph"
	"github.com/s3kim2018/validator/k8s/clientset"
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

func (gw *GraphWalker) Walk2(sourceNodes []graph.NodeInterface, filterOpts []graph.NodeFilterOption) []graph.NodeInterface {
	var nodeInterfaceList []graph.NodeInterface
	if gw.Tracker == nil {
		gw.Tracker = make(map[graph.NodeInterface]struct{})
	}
	if gw.Result == nil {
		gw.Result = make(map[graph.NodeInterface][]graph.NodeInterface)
	}
	for _, sourceNode := range sourceNodes {
		fmt.Printf("source %s:%s:%s\n", sourceNode.Plane(), sourceNode.Type(), sourceNode.Name())
		gw.Tracker[sourceNode] = struct{}{}
		for _, filterOpt := range filterOpts {
			if filterOpt.ID != "" {
				fmt.Println("ID", filterOpt.ID)
			}
			nodeInterfaceList = append(nodeInterfaceList, gw.G.GetNodeEdge(sourceNode, filterOpt)...)
			targetNodes := gw.G.GetNodeEdge(sourceNode, filterOpt)
			if len(targetNodes) == 0 {
				theplane := string(filterOpt.NodePlane)
				thetype := string(filterOpt.NodeType)
				errNode := &graph.ErrorNode{
					NodeName:   "ErrorNode, Cannot Find Node of Plane: " + theplane + " And Type: " + thetype + " From Node: " + sourceNode.Name(),
					ErrorPlane: sourceNode.Plane(),
				}
				if _, ok := gw.Result[errNode]; ok {
					gw.Result[errNode] = append(gw.Result[errNode], sourceNode)
				} else {
					gw.Result[errNode] = []graph.NodeInterface{sourceNode}
				}
				gw.Result[sourceNode] = append(gw.Result[sourceNode], errNode)
				fmt.Printf("--> %s:%s:%s\n", errNode.Plane(), errNode.Type(), errNode.Name())
				return nodeInterfaceList
			} else {
				for _, targetNode := range targetNodes {
					if filterOpt.TargetFilter != "" {
						for trackerNode := range gw.Tracker {
							if trackerNode.Type() == filterOpt.TargetFilter && trackerNode.Name() == targetNode.Name() {
								fmt.Printf("--> %s:%s:%s\n", targetNode.Plane(), targetNode.Type(), targetNode.Name())
								gw.Result[sourceNode] = append(gw.Result[sourceNode], targetNode)
								gw.Result[targetNode] = append(gw.Result[targetNode], sourceNode)
								//gw.Result[sourceNode] = append(gw.Result[sourceNode], trackerNode)
								//gw.Result[trackerNode] = append(gw.Result[trackerNode], sourceNode)
							}
						}
					} else {
						fmt.Printf("--> %s:%s:%s\n", targetNode.Plane(), targetNode.Type(), targetNode.Name())
						gw.Result[sourceNode] = append(gw.Result[sourceNode], targetNode)
						gw.Result[targetNode] = append(gw.Result[targetNode], sourceNode)
					}
				}
			}
		}
	}
	return nodeInterfaceList
}

var walkerMap = make(map[graph.NodeType]func(GraphWalker) map[graph.NodeInterface][]graph.NodeInterface)

var walker2Map = make(map[graph.NodeType]func(g *graph.Graph, nodeType graph.NodeType, plane graph.Plane, name string) map[graph.NodeInterface][]graph.NodeInterface)

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
	FilterOpts []graph.NodeFilterOption
	WalkerFunc func(sourceNodes []graph.NodeInterface, filterOpts []graph.NodeFilterOption) []graph.NodeInterface
}

type Configurations struct {
	Sourcenode string   `json:"sourcenode"`
	Next       []Config `json:"next"`
}

type Config struct {
	Thetype  string   `json:"type"`
	Theplane string   `json:"plane"`
	Next     []Config `json:"next"`
}

func Dynamicwalk(client *clientset.Client, config Configurations) map[graph.NodeInterface][]graph.NodeInterface {
	g := builder.BuildGraph(client)
	g.EdgeMatcher()
	fmt.Println(config)
	sourcelst := strings.Split(config.Sourcenode, "-")
	nodetype := graph.TypeMap[sourcelst[len(sourcelst)-1]]
	nodeplane := graph.PlaneMap[sourcelst[len(sourcelst)-2]]
	name := strings.Join(sourcelst[:len(sourcelst)-2], "-")
	nodeInterface := g.GetNodeByTypePlaneName(nodetype, nodeplane, name)
	fmt.Println(nodeInterface, "Thenodeinterface")
	var sourceNodeInterfaceList []graph.NodeInterface
	sourceNodeInterfaceList = append(sourceNodeInterfaceList, nodeInterface)
	graphWalker := GraphWalker{
		G:           g,
		SourceNodes: sourceNodeInterfaceList,
	}
	w := Walker{}
	for _, elem := range config.Next {
		pointer := &w
		pointer.helper(elem, &graphWalker)
	}
	w.runner(sourceNodeInterfaceList)
	fmt.Println("The Walker: ", w)
	fmt.Println("The Result", graphWalker.Result)
	return graphWalker.Result

}
func (w *Walker) helper(queue Config, graphwalker *GraphWalker) {
	configuration := queue
	nodetype := graph.TypeMap[configuration.Thetype]
	nodeplane := graph.PlaneMap2[configuration.Theplane]
	fmt.Println(configuration.Thetype)
	w.FilterOpts = append(w.FilterOpts, graph.NodeFilterOption{
		NodeType:  nodetype,
		NodePlane: nodeplane,
	})
	w.WalkerFunc = graphwalker.Walk2
	for _, newconfig := range configuration.Next {
		newwalker := Walker{}
		newwalkerpointer := &newwalker
		newwalkerpointer.helper(newconfig, graphwalker)
		w.Next = append(w.Next, newwalker)
	}
}
func WalkTest(client *clientset.Client, nodeType graph.NodeType, plane graph.Plane, name string) map[graph.NodeInterface][]graph.NodeInterface {
	g := builder.BuildGraph(client)
	g.EdgeMatcher()
	return podToVrouter(g, nodeType, plane, name)
	//return walker2Map[nodeType](g, nodeType, plane, name)
}

func (w Walker) runner(sourceNodes []graph.NodeInterface) {
	// fmt.Println("Runner Running!")
	// for _, elem := range w.FilterOpts {
	// 	fmt.Println(elem.NodeType)
	// }
	fmt.Println(w.FilterOpts)
	nextSourceNodes := w.WalkerFunc(sourceNodes, w.FilterOpts)
	// fmt.Println("WHAT", len(w.Next))
	// fmt.Println(w.Next)
	if len(w.Next) > 0 {
		for _, nextWalker := range w.Next {
			nextWalker.runner(nextSourceNodes)
		}
	}
}
