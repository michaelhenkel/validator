package cmd

import (
	"fmt"

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
			//nodeEdges := walker.Walk(Client, graph.NodeType(NodeType), graph.Plane(Plane), Name)
			//fmt.Println(nodeEdges)

			nodeEdges := walker.WalkTest(Client, graph.NodeType(NodeType), graph.Plane(Plane), Name)
			fmt.Printf("%+v\n", nodeEdges)
		}
	},
}
