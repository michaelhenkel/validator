package cmd

import (
	"net/http"
	"strings"

	"github.com/michaelhenkel/validator/builder"
	"github.com/michaelhenkel/validator/graph"
	"github.com/michaelhenkel/validator/walker"
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
			nodeEdges = walker.Walk(Client, graph.NodeType(varList[0]), graph.Plane(varList[1]), varList[2])
		}
	} else {
		nodeEdges = print()

	}
	page := builder.RenderPage(nodeEdges)
	if err := page.Render(w); err != nil {
		panic(err)
	}

}
