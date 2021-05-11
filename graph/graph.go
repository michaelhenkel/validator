package graph

import (
	"fmt"
	"reflect"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/michaelhenkel/validator/k8s/clientset"
)

type NodeType string
type Plane string
type Category string
type Shape string

const (
	Vrouter                 NodeType = "vrouter"
	VirtualRouter           NodeType = "virtualRouter"
	Pod                     NodeType = "pod"
	Control                 NodeType = "control"
	BGPRouter               NodeType = "bgpRouter"
	BGPNeighbor             NodeType = "bgpNeighbor"
	ConfigMap               NodeType = "configMap"
	ConfigFile              NodeType = "configFile"
	RoutingInstance         NodeType = "routingInstance"
	VirtualMachineInterface NodeType = "virtualMachineInterface"
	VirtualMachine          NodeType = "virtualMachine"
	VirtualNetwork          NodeType = "virtualNetwork"
	Kubemanager             NodeType = "kubemanager"
	K8SNode                 NodeType = "k8snode"
	DeploymentResource      Shape    = "box"
	ConfigResource          Shape    = "oval"
	ConfigFileResource      Shape    = "hexagon"
	ConfigPlane             Plane    = "configPlane"
	ControlPlane            Plane    = "controlPlane"
	DataPlane               Plane    = "dataPlane"
	ControlCategory         Category = "controlCategory"
	DataCategory            Category = "dataCategory"
	DeploymentCategory      Category = "deploymentCategory"
)

var categoryColorMap = map[Category]*opts.ItemStyle{
	ControlCategory: {
		Color: "violet",
	},
	DataCategory: {
		Color: "green",
	},
	DeploymentCategory: {
		Color: "blue",
	},
}

var planeSymbolMap = map[Plane]string{
	ConfigPlane:  "roundRect",
	ControlPlane: "circle",
	DataPlane:    "diamond",
}

var categoryMap = map[NodeType]Category{
	Vrouter:                 DeploymentCategory,
	VirtualRouter:           DataCategory,
	Pod:                     DeploymentCategory,
	Control:                 DeploymentCategory,
	Kubemanager:             DeploymentCategory,
	BGPRouter:               ControlCategory,
	BGPNeighbor:             ControlCategory,
	RoutingInstance:         ControlCategory,
	VirtualNetwork:          ControlCategory,
	VirtualMachineInterface: DataCategory,
	VirtualMachine:          DataCategory,
	ConfigMap:               DeploymentCategory,
	ConfigFile:              DeploymentCategory,
}

func Convert(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		destin = x.Interface()
	}
}

type GraphNode struct {
	Node NodeInterface
	ID   int64
}

type Graph struct {
	nodes        map[GraphNode]struct{}
	ClientConfig *clientset.Client
	NodeEdges    map[NodeInterface][]NodeInterface
}

func (gn *GraphNode) setID(id int64) {
	gn.ID = id
}

func (gn *GraphNode) getID() int64 {
	return gn.ID
}

func (g *Graph) getNodeIDs() map[int64]struct{} {
	var nodeIDMap = make(map[int64]struct{})
	for node := range g.nodes {
		nodeIDMap[node.getID()] = struct{}{}
	}
	return nodeIDMap
}

func (g *Graph) nextNodeID() int64 {
	nodeIDMap := g.getNodeIDs()
	var nodeIDCounter int64 = 1
	for nodeID, _ := range nodeIDMap {
		if _, ok := nodeIDMap[nodeID+nodeIDCounter]; !ok {
			return nodeID + nodeIDCounter
		} else {
			nodeIDCounter++
		}
	}
	return nodeIDCounter
}

func (g *Graph) GetNodesByTypePlane(nodeType NodeType, plane Plane) (nodeInterfaceList []NodeInterface) {
	for graphNode := range g.nodes {
		if graphNode.Node.Type() == nodeType && graphNode.Node.Plane() == plane {
			nodeInterfaceList = append(nodeInterfaceList, graphNode.Node)
		}
	}
	return nodeInterfaceList
}

func (g *Graph) GetNodeByTypePlaneName(nodeType NodeType, plane Plane, name string) (nodeInterface NodeInterface) {
	for graphNode := range g.nodes {
		if graphNode.Node.Type() == nodeType && graphNode.Node.Plane() == plane && graphNode.Node.Name() == name {
			return graphNode.Node
		}
	}
	return nodeInterface
}

func NewGraph(clientConfig *clientset.Client) *Graph {
	return &Graph{
		nodes:        make(map[GraphNode]struct{}),
		ClientConfig: clientConfig,
	}
}

func (g *Graph) EdgeMatcher() {
	g.NodeEdges = make(map[NodeInterface][]NodeInterface)
	for graphNode := range g.nodes {
		for _, nodeEdgeSelector := range graphNode.Node.GetEdgeSelectors() {
			for _, dstNodeInterface := range g.GetNodesByTypePlane(nodeEdgeSelector.NodeType, nodeEdgeSelector.Plane) {
				for _, dstEdgeLabel := range dstNodeInterface.GetEdgeLabels() {
					match := true
					for _, matchValue := range nodeEdgeSelector.MatchValues {
						for matchValueK, matchValueV := range matchValue.Value {
							if dstMatchValue, ok := dstEdgeLabel.Value[matchValueK]; !ok || dstMatchValue != matchValueV {
								match = false
							} else {
								match = true
								continue
							}
						}
					}
					if match {
						fmt.Printf("Found edge from %s:%s to %s:%s\n", graphNode.Node.Type(), graphNode.Node.Name(), dstNodeInterface.Type(), dstNodeInterface.Name())
						g.NodeEdges[graphNode.Node] = append(g.NodeEdges[graphNode.Node], dstNodeInterface)
						g.NodeEdges[dstNodeInterface] = append(g.NodeEdges[dstNodeInterface], graphNode.Node)
					}
				}
			}
		}
	}
}

func (g *Graph) NodeAdder(adder func(g *Graph) ([]NodeInterface, error)) *Graph {
	graphNodeResourceList, err := adder(g)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	for _, graphNodeResource := range graphNodeResourceList {
		graphNode := GraphNode{Node: graphNodeResource}
		numNodes := len(g.nodes)
		if numNodes == 0 {
			graphNode.setID(0)
			g.nodes[graphNode] = struct{}{}
		} else {
			graphNode.setID(g.nextNodeID())
			g.nodes[graphNode] = struct{}{}
		}
	}
	return g
}

type MatchValue struct {
	Value     map[string]string
	MustMatch bool
}

type EdgeSelector struct {
	NodeType    NodeType
	Plane       Plane
	MatchValues []MatchValue
}

type EdgeLabel struct {
	Value map[string]string
}

type NodeInterface interface {
	Convert(NodeInterface) error
	Type() NodeType
	Plane() Plane
	Name() string
	GetEdgeSelectors() []EdgeSelector
	GetEdgeLabels() []EdgeLabel
}

func (g *Graph) String() {
	for k := range g.nodes {
		fmt.Printf("%d:%s:%s\n", k.getID(), k.Node.Type(), k.Node.Name())
	}
}

type Node interface {
	GetID() int64
	SetID(int64)
	GetNodeType() NodeType
	GetName() string
	Shape() Shape
}

type Edge struct {
	Source NodeInterface
	Target NodeInterface
}

type NewNode struct {
	NodeType NodeType
	Resource interface{}
}

type NodeFilterOption struct {
	NodeType    NodeType
	NodeName    string
	NodePlane   Plane
	MapSetter   func(string)
	MapGetter   func() string
	KeyValueMap func(NodeType, string)
}

func (g *Graph) GetNodeEdge(node NodeInterface, filterOpts NodeFilterOption) []NodeInterface {
	var nodeEdgeList []NodeInterface
	if sourceNode, ok := g.NodeEdges[node]; ok {
		for _, targetNode := range sourceNode {
			if targetNode.Plane() == filterOpts.NodePlane &&
				targetNode.Type() == filterOpts.NodeType {
				nodeEdgeList = append(nodeEdgeList, targetNode)
			}

		}

	}
	return nodeEdgeList
}

func (g *Graph) graphNodes() []opts.GraphNode {
	var graphNodes []opts.GraphNode
	for k := range g.nodes {
		graphNode := opts.GraphNode{
			Name:       fmt.Sprintf("%s:%s", k.Node.Type(), k.Node.Name()),
			SymbolSize: 40,
		}
		graphNode.Symbol = planeSymbolMap[k.Node.Plane()]
		graphNode.ItemStyle = categoryColorMap[categoryMap[k.Node.Type()]]
		//graphNode.Category = 0
		graphNodes = append(graphNodes, graphNode)
	}
	return graphNodes
}

type GraphWalker struct {
	SourceNodes []NodeInterface
	G           *Graph
	KeyValueMap map[NodeType]string
}

type NodeTypePlane struct {
	NodeType  NodeType
	NodePlane Plane
	NodeName  string
}

func (gw *GraphWalker) Walk(filterOpts NodeFilterOption) *GraphWalker {
	var nodeInterfaceList []NodeInterface
	for _, sourceNode := range gw.SourceNodes {
		if filterOpts.KeyValueMap != nil {
			filterOpts.KeyValueMap(sourceNode.Type(), sourceNode.Name())
		}
		if filterOpts.MapGetter != nil {
			value := filterOpts.MapGetter()
			fmt.Println("Value", value)
		}
		fmt.Printf("source %s:%s:%s\n", sourceNode.Plane(), sourceNode.Type(), sourceNode.Name())
		nodeInterfaceList = append(nodeInterfaceList, gw.G.GetNodeEdge(sourceNode, filterOpts)...)
		for _, node := range gw.G.GetNodeEdge(sourceNode, filterOpts) {
			fmt.Printf("--> %s:%s:%s\n", node.Plane(), node.Type(), node.Name())
		}
	}
	gw.SourceNodes = nodeInterfaceList
	return gw
}

func (g *Graph) graphBase() *charts.Graph {
	cgraph := charts.NewGraph()

	cgraph.Initialization.Width = "100%"
	cgraph.Initialization.Height = "2000px"
	cgraph.TextStyle = &opts.TextStyle{
		FontSize: 120,
	}
	cgraph.SubtitleStyle = &opts.TextStyle{
		FontSize: 120,
	}
	cgraph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "contrail graph"}),
	)
	var edges []opts.GraphLink
	for sourceNode, targetNodes := range g.NodeEdges {
		for _, targetNode := range targetNodes {
			edge := opts.GraphLink{
				Source: fmt.Sprintf("%s:%s", sourceNode.Type(), sourceNode.Name()),
				Target: fmt.Sprintf("%s:%s", targetNode.Type(), targetNode.Name()),
				Value:  10,
			}
			edges = append(edges, edge)
		}
	}
	cgraph.AddSeries("graph", g.graphNodes(), edges,
		charts.WithGraphChartOpts(
			opts.GraphChart{
				FocusNodeAdjacency: true,
				Roam:               true,
				Force: &opts.GraphForce{
					Repulsion:  1000,
					EdgeLength: 60,
				},
				/*
					Categories: []*opts.GraphCategory{{
						Name: "cat1",
						Label: &opts.Label{
							Show:      true,
							Color:     "red",
							Position:  "top",
							Formatter: "",
						},
					}},
				*/
			},
		),
	)
	return cgraph
}

func (g *Graph) RenderPage() *components.Page {
	page := components.NewPage()
	page.SetLayout(components.PageNoneLayout)
	return page.AddCharts(
		g.graphBase(),
	)
}
