package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/michaelhenkel/validator/graph"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const PPP graph.NodeType = "dada"

type PodNode struct {
	Resource      corev1.Pod
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *PodNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*PodNode)
	if !ok {
		return fmt.Errorf("not a virtualrouter resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *PodNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *PodNode) Name() string {
	return r.Resource.Name
}

func (r *PodNode) Type() graph.NodeType {
	return graph.Pod
}

func (r *PodNode) Plane() graph.Plane {
	return plane
}

func (r *PodNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *PodNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

type virtualNetworkAnnotation struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (r *PodNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var nodeInterfaceList []graph.NodeInterface
	vrouterList, err := g.ClientConfig.DeployerDataV1.Vrouters("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range vrouterList.Items {
		vrouterNodeList, err := r.buildNodeList(resource.Name, g, graph.Vrouter)
		if err != nil {
			return nil, err
		}
		nodeInterfaceList = append(nodeInterfaceList, vrouterNodeList...)
	}
	controlList, err := g.ClientConfig.DeployerControlV1.Controls("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range controlList.Items {
		controlList, err := r.buildNodeList(resource.Name, g, graph.Control)
		if err != nil {
			return nil, err
		}
		nodeInterfaceList = append(nodeInterfaceList, controlList...)
	}

	overlayPodList, err := g.ClientConfig.CoreV1.Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range overlayPodList.Items {
		if !resource.Spec.HostNetwork {
			var edgeSelectorList []graph.EdgeSelector
			kubemanagerList, err := g.ClientConfig.DeployerConfigV1.Kubemanagers("").List(context.Background(), metav1.ListOptions{})
			if err != nil {
				fmt.Println(err)
			}
			for _, kubemanager := range kubemanagerList.Items {
				defaultPodNetworkNamespaceName := fmt.Sprintf("%s/%s-default-podnetwork", kubemanager.Spec.Namespace, kubemanager.Spec.ClusterName)
				virtualNetworkEdgeSelector := graph.EdgeSelector{
					NodeType: graph.VirtualNetwork,
					Plane:    graph.ConfigPlane,
					MatchValues: []graph.MatchValue{{
						Value: map[string]string{"VirtualNetworkNamespaceName": defaultPodNetworkNamespaceName},
					}},
				}
				edgeSelectorList = append(edgeSelectorList, virtualNetworkEdgeSelector)
			}

			virtualMachineEdgeSelector := graph.EdgeSelector{
				NodeType: graph.VirtualMachine,
				Plane:    graph.ConfigPlane,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"PodNamespaceName": fmt.Sprintf("%s/%s", resource.Namespace, resource.Name)},
				}},
			}
			edgeSelectorList = append(edgeSelectorList, virtualMachineEdgeSelector)
			if networkAnnotation, ok := resource.Annotations["k8s.v1.cni.cncf.io/networks"]; ok {
				networkAnnotationsList := strings.Split(networkAnnotation, "/")
				var networkNamespaceName string
				if len(networkAnnotationsList) == 2 {
					networkNamespaceNameList := strings.Split(networkAnnotation, "@")
					networkNamespaceName = networkNamespaceNameList[0]
				} else {
					//vnAnnotation := &virtualNetworkAnnotation{}
					var vnAnnotationArray []virtualNetworkAnnotation
					if err := json.Unmarshal([]byte(networkAnnotation), &vnAnnotationArray); err != nil {
						fmt.Println(err)
					}
					vnName := strings.Split(vnAnnotationArray[0].Name, "@")
					networkNamespaceName = fmt.Sprintf("%s/%s", vnAnnotationArray[0].Namespace, vnName[0])
				}
				virtualNetworkEdgeSelector := graph.EdgeSelector{
					NodeType: graph.VirtualNetwork,
					Plane:    graph.ConfigPlane,
					MatchValues: []graph.MatchValue{{
						Value: map[string]string{"VirtualNetworkNamespaceName": networkNamespaceName},
					}},
				}
				edgeSelectorList = append(edgeSelectorList, virtualNetworkEdgeSelector)
			}
			nodeEdgeSelector := graph.EdgeSelector{
				NodeType: graph.K8SNode,
				Plane:    graph.ConfigPlane,
				MatchValues: []graph.MatchValue{{
					Value: map[string]string{"NodeIP": resource.Status.HostIP},
				}},
			}
			edgeSelectorList = append(edgeSelectorList, nodeEdgeSelector)
			resourceNode := &PodNode{
				Resource:      resource,
				EdgeSelectors: edgeSelectorList,
			}
			nodeInterfaceList = append(nodeInterfaceList, resourceNode)
		}
	}
	return nodeInterfaceList, nil
}

func (r *PodNode) buildNodeList(filter string, g *graph.Graph, nodeType graph.NodeType) ([]graph.NodeInterface, error) {
	var nodeInterfaceList []graph.NodeInterface
	opts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", filter),
	}
	resourceList, err := g.ClientConfig.CoreV1.Pods("").List(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	for _, resource := range resourceList.Items {
		var edgeSelectorList []graph.EdgeSelector
		for _, volume := range resource.Spec.Volumes {
			if volume.Name == "config-volume" {
				if volume.VolumeSource.ConfigMap != nil {
					edgeSelector := graph.EdgeSelector{
						NodeType: graph.ConfigMap,
						Plane:    graph.ConfigPlane,
						MatchValues: []graph.MatchValue{{
							Value: map[string]string{"ConfigMap": volume.ConfigMap.Name},
						}},
					}
					edgeSelectorList = append(edgeSelectorList, edgeSelector)
				}
			}
		}
		var value = make(map[string]string)
		var ntype graph.NodeType
		switch nodeType {
		case graph.Control:
			value = map[string]string{"BGPRouterIP": resource.Status.PodIP}
			ntype = graph.BGPRouter
		case graph.Vrouter:
			value = map[string]string{"VirtualRouterIP": resource.Status.PodIP}
			ntype = graph.VirtualRouter
		}

		edgeSelector := graph.EdgeSelector{
			NodeType: ntype,
			Plane:    graph.ConfigPlane,
			MatchValues: []graph.MatchValue{{
				Value: value,
			}},
		}
		edgeSelectorList = append(edgeSelectorList, edgeSelector)

		nodeEdgeSelector := graph.EdgeSelector{
			NodeType: graph.K8SNode,
			Plane:    graph.ConfigPlane,
			MatchValues: []graph.MatchValue{{
				Value: map[string]string{"NodeIP": resource.Status.HostIP},
			}},
		}
		edgeSelectorList = append(edgeSelectorList, nodeEdgeSelector)

		r.Resource = resource
		resourceNode := &PodNode{
			Resource:      resource,
			EdgeSelectors: edgeSelectorList,
			EdgeLabels: []graph.EdgeLabel{{
				Value: map[string]string{"app": filter},
			}},
		}
		nodeInterfaceList = append(nodeInterfaceList, resourceNode)
	}
	return nodeInterfaceList, nil
}
