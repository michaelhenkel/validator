package render

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/michaelhenkel/validator/graph"
)

func graphNodes(nodeEdges map[graph.NodeInterface][]graph.NodeInterface) []opts.GraphNode {
	var graphNodes []opts.GraphNode
	for k := range nodeEdges {
		graphNode := opts.GraphNode{
			Name:       fmt.Sprintf("%s:%s:%s", k.Plane(), k.Type(), k.Name()),
			SymbolSize: 40,
		}
		graphNode.Symbol = graph.PlaneSymbolMap[k.Plane()]
		graphNode.ItemStyle = graph.CategoryColorMap[graph.CategoryMap[k.Type()]]
		graphNodes = append(graphNodes, graphNode)
	}
	return graphNodes
}

func graphBase(nodeEdges map[graph.NodeInterface][]graph.NodeInterface) *charts.Graph {
	cgraph := charts.NewGraph()

	cgraph.Initialization.Width = "100%"
	cgraph.Initialization.Height = "2000px"
	cgraph.TextStyle = &opts.TextStyle{
		FontSize: 120,
	}
	cgraph.SubtitleStyle = &opts.TextStyle{
		FontSize: 120,
	}
	cgraph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "contrail graph"}),
	)
	var edges []opts.GraphLink
	for sourceNode, targetNodes := range nodeEdges {
		for _, targetNode := range targetNodes {
			edge := opts.GraphLink{
				Source: fmt.Sprintf("%s:%s:%s", sourceNode.Plane(), sourceNode.Type(), sourceNode.Name()),
				Target: fmt.Sprintf("%s:%s:%s", targetNode.Plane(), targetNode.Type(), targetNode.Name()),
				Value:  10,
			}
			edges = append(edges, edge)
		}
	}
	cgraph.AddSeries("graph", graphNodes(nodeEdges), edges,
		charts.WithGraphChartOpts(
			opts.GraphChart{
				FocusNodeAdjacency: true,
				Roam:               true,
				Force: &opts.GraphForce{
					Repulsion:  1000,
					EdgeLength: 60,
				},
			},
		),
	)
	return cgraph
}

func RenderPage(nodeEdges map[graph.NodeInterface][]graph.NodeInterface) *components.Page {
	page := components.NewPage()
	page.SetLayout(components.PageNoneLayout)
	return page.AddCharts(
		graphBase(nodeEdges),
	)
}
