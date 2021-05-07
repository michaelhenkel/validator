package cmd

import (
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the graph",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			Name = args[0]
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		serve()

	},
}

func printGraph(w http.ResponseWriter, req *http.Request) {

	/*
		validator := validate.NewValidator(ClientConfig)
		if err := validator.Validate(); err != nil {
			panic(err)
		}
		page := validator.RenderPage()
		if err := page.Render(w); err != nil {
			panic(err)
		}
	*/

}

func serve() {
	http.HandleFunc("/", printGraph)
	http.ListenAndServe(":8090", nil)
}
