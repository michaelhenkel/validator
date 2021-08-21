package resources

import (
	"context"
	"fmt"

	"github.com/s3kim2018/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigFile struct {
	Name   string
	Config string
}

type ConfigFileNode struct {
	Resource      ConfigFile
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *ConfigFileNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*ConfigFileNode)
	if !ok {
		return fmt.Errorf("not a ConfigFileNode resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *ConfigFileNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *ConfigFileNode) Name() string {
	return r.Resource.Name
}

func (r *ConfigFileNode) Type() graph.NodeType {
	return graph.ConfigFile
}

func (r *ConfigFileNode) Plane() graph.Plane {
	return plane
}

func (r *ConfigFileNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *ConfigFileNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *ConfigFileNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
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
			for cmDataKey, cmDataValue := range configMap.Data {
				configFile := ConfigFile{
					Name:   cmDataKey,
					Config: cmDataValue,
				}
				resourceNode := &ConfigFileNode{
					Resource: configFile,
					EdgeLabels: []graph.EdgeLabel{{
						Value: map[string]string{"ConfigMap": k + "-configmap"},
					}},
				}

				graphNodeList = append(graphNodeList, resourceNode)
			}
		}
	}
	return graphNodeList, nil
}
