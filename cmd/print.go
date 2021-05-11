package cmd

import (
	"github.com/michaelhenkel/validator/builder"
	"github.com/michaelhenkel/validator/graph"
	"github.com/michaelhenkel/validator/walker"
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
			graphWalker := walker.GraphWalker{
				G:           g,
				SourceNodes: sourceNodeInterfaceList,
			}
			graphWalker.Walk(graph.NodeFilterOption{
				NodeType:  graph.VirtualMachine,
				NodePlane: graph.ConfigPlane,
			}).Walk(graph.NodeFilterOption{
				NodeType:  graph.VirtualMachineInterface,
				NodePlane: graph.ConfigPlane,
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
				NodeType:     graph.VirtualMachine,
				NodePlane:    graph.ConfigPlane,
				TargetFilter: graph.VirtualMachine,
			})
		}
	},
}
