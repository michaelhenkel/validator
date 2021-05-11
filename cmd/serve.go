package cmd

import (
	"net/http"
	"strings"

	"github.com/michaelhenkel/validator/builder"
	"github.com/michaelhenkel/validator/graph"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the graph",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		serve()

	},
}

func print() map[graph.NodeInterface][]graph.NodeInterface {
	g := builder.BuildGraph(Client)
	return g.NodeEdges
}

func serve() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8090", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var nodeEdges map[graph.NodeInterface][]graph.NodeInterface
	if len(r.URL.Query()) > 0 {
		if vars, ok := r.URL.Query()["walk"]; ok {
			varList := strings.Split(vars[0], ",")
			nodeEdges = walk(graph.NodeType(varList[0]), graph.Plane(varList[1]), varList[2])
		}
	} else {
		nodeEdges = print()

	}
	page := builder.RenderPage(nodeEdges)
	if err := page.Render(w); err != nil {
		panic(err)
	}

}

func walk(nodeType graph.NodeType, plane graph.Plane, name string) map[graph.NodeInterface][]graph.NodeInterface {
	g := builder.BuildGraph(Client)
	g.EdgeMatcher()
	nodeInterface := g.GetNodeByTypePlaneName(nodeType, plane, name)
	var sourceNodeInterfaceList []graph.NodeInterface
	sourceNodeInterfaceList = append(sourceNodeInterfaceList, nodeInterface)
	graphWalker := graph.GraphWalker{
		G:           g,
		SourceNodes: sourceNodeInterfaceList,
	}
	return walkerMap[nodeType](graphWalker)
}

var walkerMap = make(map[graph.NodeType]func(graph.GraphWalker) map[graph.NodeInterface][]graph.NodeInterface)

func podWalker(graphWalker graph.GraphWalker) map[graph.NodeInterface][]graph.NodeInterface {
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
