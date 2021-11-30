package resources

import (
	"github.com/s3kim2018/validator/k8s/clientset"
)

type Errornode struct {
	Name  string
	Edges []string
}

func Builderrorgraph(clientConfig *clientset.Client) []Errornode {
	nodes, err := VirtualMachineInterfaceErrorNode(clientConfig)
	//node2, err2 := RoutingInstanceErrorNode(clientConfig)
	if err != nil {
		return []Errornode{}
	}
	// if err2 != nil {
	// 	return []Errornode{}
	// }
	// nodes = append(nodes, node2...)

	return nodes

}
