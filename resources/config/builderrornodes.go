package resources

import (
	"github.com/michaelhenkel/validator/k8s/clientset"
)

type errornode struct {
	Name  string
	Edges []string
}

func builderrorgraph(clientConfig *clientset.Client) []errornode {
	nodes, err := VirtualMachineInterfaceErrorNode(clientConfig)
	if err != nil {
		return []errornode{}
	}
	return nodes

}
