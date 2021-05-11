package cmd

import (
	"fmt"
	"net/http"

	"github.com/michaelhenkel/validator/builder"
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

func printGraph(w http.ResponseWriter, req *http.Request) {
	fmt.Println("bla")
	g := builder.BuildGraph(Client)
	page := builder.RenderPage(g)
	if err := page.Render(w); err != nil {
		panic(err)
	}
}

func serve() {
	http.HandleFunc("/", printGraph)
	http.ListenAndServe(":8090", nil)
}
