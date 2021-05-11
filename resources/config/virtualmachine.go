package resources

import (
	"context"
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

type VirtualMachineNode struct {
	Resource      contrailcorev1alpha1.VirtualMachine
	EdgeLabels    []graph.EdgeLabel
	EdgeSelectors []graph.EdgeSelector
}

func (r *VirtualMachineNode) Convert(in graph.NodeInterface) error {
	_, ok := in.(*VirtualMachineNode)
	if !ok {
		return fmt.Errorf("not a VirtualMachineNode resource")
	}
	graph.Convert(in, r)
	return nil
}

func (r *VirtualMachineNode) AdderFunc() func(g *graph.Graph) ([]graph.NodeInterface, error) {
	return r.Adder
}

func (r *VirtualMachineNode) Name() string {
	return r.Resource.Name
}

func (r *VirtualMachineNode) Type() graph.NodeType {
	return graph.VirtualMachine
}

func (r *VirtualMachineNode) Plane() graph.Plane {
	return plane
}

func (r *VirtualMachineNode) GetEdgeLabels() []graph.EdgeLabel {
	return r.EdgeLabels
}

func (r *VirtualMachineNode) GetEdgeSelectors() []graph.EdgeSelector {
	return r.EdgeSelectors
}

func (r *VirtualMachineNode) Adder(g *graph.Graph) ([]graph.NodeInterface, error) {
	var graphNodeList []graph.NodeInterface
	resourceList, err := g.ClientConfig.ContrailCoreV1.VirtualMachines().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, resource := range resourceList.Items {
		r.Resource = resource
		resourceNode := &VirtualMachineNode{
			Resource: resource,
			EdgeLabels: []graph.EdgeLabel{{
				Value: map[string]string{"PodNamespaceName": fmt.Sprintf("%s/%s", resource.Spec.ServerNamespace, resource.Spec.ServerName)},
			}, {
				Value: map[string]string{"VirtualMachineName": resource.Name},
			}},
		}
		graphNodeList = append(graphNodeList, resourceNode)
	}
	return graphNodeList, nil
}
