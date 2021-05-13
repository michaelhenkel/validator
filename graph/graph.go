package graph

import (
	"fmt"
	"reflect"

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
	ErrNode                 NodeType = "errorNode"
	DeploymentResource      Shape    = "box"
	ConfigResource          Shape    = "oval"
	ConfigFileResource      Shape    = "hexagon"
	ConfigPlane             Plane    = "configPlane"
	ControlPlane            Plane    = "controlPlane"
	DataPlane               Plane    = "dataPlane"
	ControlCategory         Category = "controlCategory"
	DataCategory            Category = "dataCategory"
	DeploymentCategory      Category = "deploymentCategory"
	ErrorCategory           Category = "errorCategory"
)

var CategoryColorMap = map[Category]*opts.ItemStyle{
	ControlCategory: {
		Color: "violet",
	},
	DataCategory: {
		Color: "green",
	},
	DeploymentCategory: {
		Color: "blue",
	},
	ErrorCategory: {
		Color: "red",
	},
}

var PlaneSymbolMap = map[Plane]string{
	ConfigPlane:  "roundRect",
	ControlPlane: "circle",
	DataPlane:    "diamond",
}

var CategoryMap = map[NodeType]Category{
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
	ErrNode:                 ErrorCategory,
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
	NodeType     NodeType
	NodeName     string
	NodePlane    Plane
	TargetFilter NodeType
	ErrorMsg     string
	ID           string
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

type GraphWalker struct {
	SourceNodes []NodeInterface
	G           *Graph
	Tracker     map[NodeInterface]struct{}
	Result      map[NodeInterface][]NodeInterface
}

type NodeTypePlane struct {
	NodeType  NodeType
	NodePlane Plane
	NodeName  string
}

type ErrorNode struct {
	NodeName   string
	ErrorPlane Plane
}

func (e *ErrorNode) Convert(node NodeInterface) error {
	return nil
}

func (e *ErrorNode) Type() NodeType {
	return ErrNode
}

func (e *ErrorNode) Plane() Plane {
	return e.ErrorPlane
}

func (e *ErrorNode) Name() string {
	return e.NodeName
}

func (e *ErrorNode) GetEdgeSelectors() []EdgeSelector {
	return []EdgeSelector{}
}

func (e *ErrorNode) GetEdgeLabels() []EdgeLabel {
	return []EdgeLabel{}
}
