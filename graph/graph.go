package graph

import (
	"fmt"
)

type NodeType string
type Shape string

const (
	Vrouter            NodeType = "vrouter"
	VirtualRouter      NodeType = "virtualRouter"
	Pod                NodeType = "pod"
	Control            NodeType = "control"
	BGPRouter          NodeType = "bgpRouter"
	ConfigMap          NodeType = "configMap"
	ConfigFile         NodeType = "configFile"
	DeploymentResource Shape    = "box"
	ConfigResource     Shape    = "oval"
	ConfigFileResource Shape    = "hexagon"
)

type edge struct {
	node  Node
	label string
}

type Node interface {
	GetID() int64
	SetID(int64)
	GetNodeType() NodeType
	GetName() string
	Shape() Shape
}

type Graph struct {
	nodes map[Node][]edge
}

func NewGraph() *Graph {
	return &Graph{nodes: make(map[Node][]edge)}
}

func (g *Graph) GetNodesByNodeType(nodeType NodeType) []Node {
	var nodeList []Node
	for node, _ := range g.nodes {
		if node.GetNodeType() == nodeType {
			nodeList = append(nodeList, node)
		}
	}
	return nodeList
}

func (g *Graph) AddNode(node Node) {
	var counter int64 = 0
	numNodes := len(g.nodes)
	if numNodes == 0 {
		node.SetID(counter)
		g.nodes[node] = []edge{}
	} else {
		node.SetID(g.nextNodeID())
		g.nodes[node] = []edge{}
	}
}

func (g *Graph) getNodeIDs() map[int64]struct{} {
	var nodeIDMap = make(map[int64]struct{})
	for node := range g.nodes {
		nodeIDMap[node.GetID()] = struct{}{}
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

func (g *Graph) AddEdge(from, to Node, label string) {
	g.nodes[from] = append(g.nodes[from], edge{node: to, label: label})
}

func (g *Graph) GetEdges(node Node) []edge {
	return g.nodes[node]
}

func (e *edge) String() string {
	return fmt.Sprintf("%v", e.node.GetID())
}

func (g *Graph) String() string {
	out := `digraph finite_state_machine {
		rankdir=LR;
		size="4,5"`
	out = fmt.Sprintf("%s\n", out)
	for k := range g.nodes {
		out += fmt.Sprintf("\tnode [shape = %s]; \t\"%d:%s:%s\" \t;\n", k.Shape(), k.GetID(), k.GetNodeType(), k.GetName())

	}
	for k := range g.nodes {
		for _, v := range g.GetEdges(k) {
			out += fmt.Sprintf("\t\"%d:%s:%s\" -> \"%d:%s:%s\"\t[ label = \"%s\" ];\n", k.GetID(), k.GetNodeType(), k.GetName(), v.node.GetID(), v.node.GetNodeType(), v.node.GetName(), v.label)
		}
	}

	out += "}"
	return out
}
