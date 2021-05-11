package cmd

import (
	"github.com/michaelhenkel/validator/builder"
	"github.com/michaelhenkel/validator/graph"
	"github.com/spf13/cobra"
)

var (
	Name     string
	NodeType string
	Plane    string
)

func init() {
	rootCmd.AddCommand(printCmd)

	printCmd.PersistentFlags().StringVarP(&Name, "name", "n", "", "resource name")
	printCmd.PersistentFlags().StringVarP(&NodeType, "type", "t", "", "resource type")
	printCmd.PersistentFlags().StringVarP(&Plane, "plane", "p", "", "resource plane")
}

type myfunc func(string, string)

func test(key string, val string) {

}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "prints a graph",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		g := builder.BuildGraph(Client)
		g.String()
		g.EdgeMatcher()
		if Name != "" && NodeType != "" {
			nodeType := graph.NodeType(NodeType)
			plane := graph.Plane(Plane)
			nodeInterface := g.GetNodeByTypePlaneName(nodeType, plane, Name)
			var sourceNodeInterfaceList []graph.NodeInterface
			sourceNodeInterfaceList = append(sourceNodeInterfaceList, nodeInterface)
			graphWalker := graph.GraphWalker{
				G:           g,
				SourceNodes: sourceNodeInterfaceList,
			}
			graphWalker.Walk(graph.NodeFilterOption{
				NodeType:  graph.VirtualMachine,
				NodePlane: graph.ConfigPlane,
				KeyValueMap: func(nodeType graph.NodeType, name string) {
					graphWalker.KeyValueMap = map[graph.NodeType]string{nodeType: name}
				},
			}).Walk(graph.NodeFilterOption{
				NodeType:  graph.VirtualMachineInterface,
				NodePlane: graph.ConfigPlane,
				KeyValueMap: func(nodeType graph.NodeType, name string) {
					graphWalker.KeyValueMap = map[graph.NodeType]string{nodeType: name}
				},
			}).Walk(graph.NodeFilterOption{
				NodeType:  graph.RoutingInstance,
				NodePlane: graph.ConfigPlane,
			}).Walk(graph.NodeFilterOption{
				NodeType:  graph.RoutingInstance,
				NodePlane: graph.ControlPlane,
			}).Walk(graph.NodeFilterOption{
				NodeType:  graph.BGPNeighbor,
				NodePlane: graph.ControlPlane,
			}).Walk(graph.NodeFilterOption{
				NodeType:  graph.VirtualRouter,
				NodePlane: graph.ConfigPlane,
			}).Walk(graph.NodeFilterOption{
				NodeType:  graph.VirtualMachine,
				NodePlane: graph.ConfigPlane,
				MapGetter: func() string {
					return graphWalker.KeyValueMap[graph.VirtualMachine]
				},
			})
		}
	},
}
