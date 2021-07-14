package cmd

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
	config b "/Users/briank/Desktop/github.com/michaelhenkel/validator/resources/config"
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

func serve() {
	http.HandleFunc("/template", loadtemplate)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/oldgraph", handler)
	http.ListenAndServe(":8090", nil)
}

type templatestruct struct {
	Title     string
	Names     []string
	Planes    map[string]string
	Types     map[string]string
	UniqTypes []string
	Edges     map[string][]string
}

func loadtemplate(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("Could not parse files")
		return
	}
	var names []string
	edges := make(map[string][]string)
	planes := make(map[string]string)
	types := make(map[string]string)
	var UniqTypes []string
	g := builder.BuildGraph(Client)
	for key, value := range g.NodeEdges {
		names = append(names, key.Name()+"-"+string(key.Plane())+"-"+string(key.Type()))
		for _, elem := range value {
			edges[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = append(edges[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())], elem.Name()+"-"+string(elem.Plane())+"-"+string(elem.Type()))
		}
		planes[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = string(key.Plane())
		types[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = string(key.Type())
	}
	nodes := config.builderrorgraph(Client)
	for node := range nodes {
		names = append(names, node.Name)ßß
	}
	for _, val := range types {
		if !contains(UniqTypes, val) {
			UniqTypes = append(UniqTypes, val)
		}
	}
	structval := templatestruct{Title: "bob", Names: names, Edges: edges, Planes: planes, Types: types, UniqTypes: UniqTypes}
	template.Execute(w, structval)

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var nodeEdges map[graph.NodeInterface][]graph.NodeInterface
	if len(r.URL.Query()) > 0 {
		if vars, ok := r.URL.Query()["walk"]; ok {
			varList := strings.Split(vars[0], ",")
			//nodeEdges = walker.Walk(Client, graph.NodeType(varList[0]), graph.Plane(varList[1]), varList[2])
			nodeEdges = walker.WalkTest(Client, graph.NodeType(varList[0]), graph.Plane(varList[1]), varList[2])
		}
	} else {
		g := builder.BuildGraph(Client)
		nodeEdges = g.NodeEdges
	}
	page := builder.RenderPage(nodeEdges)
	if err := page.Render(w); err != nil {
		panic(err)
	}
	duration := time.Since(start)
	fmt.Println("Execution time: ", duration)

}
