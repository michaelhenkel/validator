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
	BGPRouter:               ControlCategory,
	BGPNeighbor:             ControlCategory,
	RoutingInstance:         ControlCategory,
	VirtualMachineInterface: DataCategory,
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
	nodes        map[GraphNode][]edge
	ClientConfig *clientset.Client
	edges        []opts.GraphLink
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

func NewGraph(clientConfig *clientset.Client) *Graph {
	return &Graph{
		nodes:        make(map[GraphNode][]edge),
		ClientConfig: clientConfig,
	}
}

func (g *Graph) EdgeMatcher() {
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
						edge := opts.GraphLink{
							Source: fmt.Sprintf("%s:%s", graphNode.Node.Type(), graphNode.Node.Name()),
							Target: fmt.Sprintf("%s:%s", dstNodeInterface.Type(), dstNodeInterface.Name()),
							Value:  10,
						}
						g.edges = append(g.edges, edge)
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
			g.nodes[graphNode] = []edge{}
		} else {
			graphNode.setID(g.nextNodeID())
			g.nodes[graphNode] = []edge{}
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

type edge struct {
	node  GraphNode
	label string
}

type Node interface {
	GetID() int64
	SetID(int64)
	GetNodeType() NodeType
	GetName() string
	Shape() Shape
}

type NewNode struct {
	NodeType NodeType
	Resource interface{}
}

func (g *Graph) AddEdge(from, to GraphNode, label string) {
	g.nodes[from] = append(g.nodes[from], edge{node: to, label: label})
}

func (g *Graph) GetEdges(graphNode GraphNode) []edge {
	return g.nodes[graphNode]
}

func (e *edge) String() string {
	return fmt.Sprintf("%v", e.node.getID())
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
		graphNodes = append(graphNodes, graphNode)
	}
	return graphNodes
}

func (g *Graph) graphBase() *charts.Graph {
	cgraph := charts.NewGraph()
	cgraph.Initialization.Width = "100%"
	cgraph.Initialization.Height = "2000px"
	cgraph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "contrail graph"}),
		/*
			charts.WithDataZoomOpts(
				opts.DataZoom{
					Type:  "slider",
					Start: 50,
					End:   100,
				},
			),
		*/
	)
	cgraph.AddSeries("graph", g.graphNodes(), g.edges,
		charts.WithGraphChartOpts(
			opts.GraphChart{
				FocusNodeAdjacency: true,
				Roam:               true,
				Force: &opts.GraphForce{
					Repulsion:  1000,
					EdgeLength: 60,
				},
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
