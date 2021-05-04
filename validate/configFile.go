package validate

import (
	"fmt"
	"strings"

	"github.com/michaelhenkel/validator/graph"
	"gopkg.in/ini.v1"

	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type ConfigFile struct {
	Name   string
	Config string
}

type ConfigFileNode struct {
	ConfigFile ConfigFile
	ID         int64
	NodeType   graph.NodeType
	Owner      graph.NodeType
}

func NewConfigFileNode(configFile ConfigFile) ConfigFileNode {
	return ConfigFileNode{
		ConfigFile: configFile,
		NodeType:   graph.ConfigFile,
	}
}

func (v *ConfigFileNode) GetNodeType() graph.NodeType {
	return v.NodeType
}

func (v *ConfigFileNode) GetID() int64 {
	return v.ID
}

func (v *ConfigFileNode) SetID(id int64) {
	v.ID = id
}

func (v *ConfigFileNode) GetName() string {
	return v.ConfigFile.Name
}

func (v *ConfigFileNode) Shape() graph.Shape {
	return graph.DeploymentResource
}

func addConfigFileNodes(validator *Validator) error {
	var graphNode graph.Node
	configMapNodes := validator.graph.GetNodesByNodeType(graph.ConfigMap)

	for _, nodeInterface := range configMapNodes {
		node, ok := nodeInterface.(*ConfigMapNode)
		if !ok {
			return fmt.Errorf("not a configMap node")
		}
		for k, v := range node.ConfigMap.Data {
			configFile := ConfigFile{
				Name:   k,
				Config: v,
			}
			configFileNode := NewConfigFileNode(configFile)
			graphNode = &configFileNode
			validator.graph.AddNode(graphNode)
		}
	}
	return nil
}

func addConfigFileToConfigMapEdges(validator *Validator) error {
	configFileNodes := validator.graph.GetNodesByNodeType(graph.ConfigFile)
	configMapNodes := validator.graph.GetNodesByNodeType(graph.ConfigMap)
	for _, configMapNodeInterface := range configMapNodes {
		configMapNode, ok := configMapNodeInterface.(*ConfigMapNode)
		if !ok {
			return fmt.Errorf("not a configMap node")
		}
		for _, configFileNodeInterface := range configFileNodes {
			configFileNode, ok := configFileNodeInterface.(*ConfigFileNode)
			if !ok {
				return fmt.Errorf("not a configFile node")
			}
			for k := range configMapNode.ConfigMap.Data {
				if k == configFileNode.ConfigFile.Name {
					validator.graph.AddEdge(configMapNode, configFileNode, "")
				}
			}
		}
	}
	return nil
}

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
