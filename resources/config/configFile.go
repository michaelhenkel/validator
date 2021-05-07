package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
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

/*

func addConfigFileToVirtualRouterBGPRouterEdges(validator *Validator) error {
	configFileNodes := validator.graph.GetNodesByNodeType(graph.ConfigFile)
	virtualRouterNodes := validator.graph.GetNodesByNodeType(graph.VirtualRouter)
	bgpRouterNodes := validator.graph.GetNodesByNodeType(graph.BGPRouter)
	for _, configFileNodeInterface := range configFileNodes {
		configFileNode, ok := configFileNodeInterface.(*ConfigFileNode)
		if !ok {
			return fmt.Errorf("not a configFile node")
		}
		hostnameArray := strings.Split(configFileNode.ConfigFile.Name, "-")
		hostnameConf := hostnameArray[len(hostnameArray)-1]
		hostname := strings.Split(hostnameConf, ".conf")
		if len(hostnameArray) > 1 {
			if hostnameArray[1] == "vrouter" {
				for _, virtualRouterNodeInterface := range virtualRouterNodes {
					virtualRouterNode, ok := virtualRouterNodeInterface.(*VirtualRouterNode)
					if !ok {
						return fmt.Errorf("not a virtualRouter node")
					}
					cfg, err := ini.Load([]byte(configFileNode.ConfigFile.Config))
					if err != nil {
						return err
					}
					controlNetworkIP := cfg.Section("NETWORKS").Key("control_network_ip").String()
					cidr := cfg.Section("VIRTUAL-HOST-INTERFACE").Key("ip").String()
					cidrArray := strings.Split(cidr, "/")
					ip := cidrArray[0]
					if virtualRouterNode.VirtualRouter.Name == hostname[0] && virtualRouterNode.VirtualRouter.Spec.IPAddress == contrailcorev1alpha1.IPAddress(controlNetworkIP) && virtualRouterNode.VirtualRouter.Spec.IPAddress == contrailcorev1alpha1.IPAddress(ip) {
						validator.graph.AddEdge(configFileNode, virtualRouterNode, "")
						break
					}
				}
			} else if hostnameArray[1] == "control" {
				for _, bgpRouterNodeInterface := range bgpRouterNodes {
					bgpRouterNode, ok := bgpRouterNodeInterface.(*BGPRouterNode)
					if !ok {
						return fmt.Errorf("not a virtualRouter node")
					}
					if bgpRouterNode.BGPRouter.Name == hostname[0] {
						validator.graph.AddEdge(configFileNode, bgpRouterNode, "")
						break
					}
				}
			}
		}

	}
	return nil
}

func addConfigFileVirtualRouterToBGPRouterEdges(validator *Validator) error {
	configFileNodes := validator.graph.GetNodesByNodeType(graph.ConfigFile)
	bgpRouterNodes := validator.graph.GetNodesByNodeType(graph.BGPRouter)
	for _, configFileNodeInterface := range configFileNodes {
		configFileNode, ok := configFileNodeInterface.(*ConfigFileNode)
		if !ok {
			return fmt.Errorf("not a configFile node")
		}
		hostnameArray := strings.Split(configFileNode.ConfigFile.Name, "-")
		if len(hostnameArray) > 1 {
			if hostnameArray[1] == "vrouter" && configFileNode.ConfigFile.Config != "" {
				cfg, err := ini.Load([]byte(configFileNode.ConfigFile.Config))
				if err != nil {
					return err
				}
				controlIPPortString := cfg.Section("CONTROL-NODE").Key("servers").String()
				controlIPPortList := strings.Split(controlIPPortString, ",")
				for _, controlIPPort := range controlIPPortList {
					controlIPPortArray := strings.Split(controlIPPort, ":")
					controlIP := controlIPPortArray[0]
					//controlPort, err := strconv.Atoi(controlIPPortArray[1])
					if err != nil {
						return err
					}
					for _, bgpRouterNodeInterface := range bgpRouterNodes {
						bgpRouterNode, ok := bgpRouterNodeInterface.(*BGPRouterNode)
						if !ok {
							return fmt.Errorf("not a control node")
						}
						if bgpRouterNode.BGPRouter.Spec.BGPRouterParameters.Address == contrailcorev1alpha1.IPAddress(controlIP) {
							validator.graph.AddEdge(configFileNode, bgpRouterNode, "")
						}
					}

				}
			}
		}

	}
	return nil
}
*/
