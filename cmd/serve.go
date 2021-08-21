package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"

	// "log"
	"net/http"
	"strings"
	"time"

	"github.com/s3kim2018/validator/builder"
	"github.com/s3kim2018/validator/graph"
	config "github.com/s3kim2018/validator/resources/config"
	"github.com/s3kim2018/validator/walker"
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
	http.HandleFunc("/", loadtemplate)
	http.HandleFunc("/walker", walkerfunc)
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

type Configurations struct {
	Sourcenode string   `json:"sourcenode"`
	Next       []Config `json:"next"`
}

type Config struct {
	Thetype  string   `json:"type"`
	Theplane string   `json:"plane"`
	Next     []Config `json:"next"`
}

func walkerfunc(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("myFile")
	var theconfig walker.Configurations
	if err != nil {
		content := r.PostFormValue("thesub")
		fmt.Println("content: ", content)
		data := []byte(content)
		err2 := json.Unmarshal(data, &theconfig)
		if err2 != nil {
			panic(err)
		}
	} else {
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)

		if err != nil {
			fmt.Println(err)
		}
		err2 := json.Unmarshal(fileBytes, &theconfig)
		if err2 != nil {
			panic(err)
		}
	}

	template, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("Could not parse files")
		return
	}
	var names []string
	edges := make(map[string][]string)
	planes := make(map[string]string)
	types := make(map[string]string)
	errnode := make(map[string][]string)
	nodes := config.Builderrorgraph(Client)
	for _, node := range nodes {
		for _, edge := range node.Edges {
			errnode[edge] = append(errnode[edge], node.Name)
		}
	}
	var UniqTypes []string
	var nodeEdges map[graph.NodeInterface][]graph.NodeInterface
	nodeEdges = walker.Dynamicwalk(Client, theconfig)
	//nodeEdges = walker.WalkTest(Client, "pod", "configPlane", "coredns-74ff55c5b-fsw2l")
	for key, value := range nodeEdges {
		names = append(names, key.Name()+"-"+string(key.Plane())+"-"+string(key.Type()))
		if errnodes, ok := errnode[key.Name()]; ok {
			for _, errnode := range errnodes {
				names = append(names, errnode+"->"+key.Name())
				edges[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = append(edges[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())], errnode+"->"+key.Name())
				planes[errnode+"->"+key.Name()] = "errorNode"
				types[errnode+"->"+key.Name()] = "errorNode"
			}
		}
		for _, elem := range value {
			edges[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = append(edges[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())], elem.Name()+"-"+string(elem.Plane())+"-"+string(elem.Type()))
		}
		planes[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = string(key.Plane())
		types[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = string(key.Type())
	}
	for _, val := range types {
		if !contains(UniqTypes, val) {
			UniqTypes = append(UniqTypes, val)
		}
	}
	structval := templatestruct{Title: "bob", Names: names, Edges: edges, Planes: planes, Types: types, UniqTypes: UniqTypes}
	template.Execute(w, structval)

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

	errnode := make(map[string][]string)
	nodes := config.Builderrorgraph(Client)
	for _, node := range nodes {
		for _, edge := range node.Edges {
			errnode[edge] = append(errnode[edge], node.Name)
		}
	}
	var UniqTypes []string
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
	for key, value := range nodeEdges {
		names = append(names, key.Name()+"-"+string(key.Plane())+"-"+string(key.Type()))
		if errnodes, ok := errnode[key.Name()]; ok {
			for _, errnode := range errnodes {
				names = append(names, errnode+"->"+key.Name())
				edges[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = append(edges[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())], errnode+"->"+key.Name())
				planes[errnode+"->"+key.Name()] = "errorNode"
				types[errnode+"->"+key.Name()] = "errorNode"
			}
		}
		for _, elem := range value {
			edges[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = append(edges[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())], elem.Name()+"-"+string(elem.Plane())+"-"+string(elem.Type()))
		}
		planes[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = string(key.Plane())
		types[key.Name()+"-"+string(key.Plane())+"-"+string(key.Type())] = string(key.Type())
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
