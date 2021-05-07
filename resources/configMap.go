package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapNode struct {
	Resource      corev1.ConfigMap
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *ConfigMapNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*ConfigMapNode)
	if !ok {
		return fmt.Errorf("not a virtualrouter resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *ConfigMapNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *ConfigMapNode) Name() string {
	return r.Resource.Name
}

func (r *ConfigMapNode) Type() graph.NodeType {
	return graph.ConfigMap
}

func (r *ConfigMapNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *ConfigMapNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *ConfigMapNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	var nameNamespaceMapList []map[string]string
	controlList, err := g.ClientConfig.DeployerControlV1.Controls("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range controlList.Items {
		nameNamespaceMap := map[string]string{resource.Name: resource.Namespace}
		nameNamespaceMapList = append(nameNamespaceMapList, nameNamespaceMap)
	}
	vrouterList, err := g.ClientConfig.DeployerDataV1.Vrouters("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range vrouterList.Items {
		nameNamespaceMap := map[string]string{resource.Name: resource.Namespace}
		nameNamespaceMapList = append(nameNamespaceMapList, nameNamespaceMap)
	}

	for _, nameNamespaceMap := range nameNamespaceMapList {
		for k, v := range nameNamespaceMap {
			configMap, err := g.ClientConfig.CoreV1.ConfigMaps(v).Get(context.Background(), k+"-configmap", metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			resourceNode := &ConfigMapNode{
				Resource: *configMap,
				EdgeLabels: []graph.EdgeLabel{{
					Value: map[string]string{"ConfigMap": configMap.Name},
				}},
				EdgeSelectors: []graph.EdgeSelector{{
					NodeType: graph.ConfigFile,
					MatchValues: []graph.MatchValue{{
						Value: map[string]string{"ConfigMap": configMap.Name},
					}},
				}},
			}
			graphNodeList = append(graphNodeList, resourceNode)
		}
	}
	return graphNodeList, nil
}
