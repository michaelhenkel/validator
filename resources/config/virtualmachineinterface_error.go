package resources

import (
	"context"
	"fmt"
	"reflect"

	"github.com/s3kim2018/validator/k8s/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

func VirtualMachineInterfaceErrorNode(clientConfig *clientset.Client) ([]Errornode, error) {
	var graphNodeList []Errornode
	resourceList, err := clientConfig.ContrailCoreV1.VirtualMachineInterfaces("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var originalresource sourcecoderesource
	originalresource.getspecandstatusvals("virtualmachineinterface")
	for i := 0; i < len(resourceList.Items); i++ {
		resource := resourceList.Items[i]
		hashmap := buildhash(clientConfig, i, "VirtualMachineInterface")
		combinedlist := append(originalresource.References, originalresource.Reference...)
		combinedlist = append(combinedlist, originalresource.Parents...)
		for i := 0; i < len(combinedlist); i++ {
			if references, ok := hashmap[combinedlist[i]]; ok {
				for j := 0; j < len(references); j++ {
					reference := references[j][0]
					switch thetype := reference.(type) {
					default:
						fmt.Println("Unexpected Type")
					case []v1alpha1.RoutingInstanceReference:
						if thetype != nil {
							for _, routingInstanceReference := range thetype {
								val := reflect.Indirect(reflect.ValueOf(routingInstanceReference))
								_, ok := val.Type().FieldByName("Name")
								if !ok {
									errnode := Errornode{
										Name:  "missing " + combinedlist[i],
										Edges: []string{resource.Name},
									}
									graphNodeList = append(graphNodeList, errnode)
								}
							}
						} else {
							errnode := Errornode{
								Name:  "missing " + combinedlist[i],
								Edges: []string{resource.Namespace + ":" + resource.Name},
							}
							fmt.Println("parent name is", resource.Namespace+":"+resource.Name)
							graphNodeList = append(graphNodeList, errnode)
						}
					case []v1alpha1.ResourceReference:
						if thetype != nil {
							for _, resourcereference := range thetype {
								val := reflect.Indirect(reflect.ValueOf(resourcereference))
								_, ok := val.Type().FieldByName("Name")
								if !ok {
									errnode := Errornode{
										Name:  "missing " + combinedlist[i],
										Edges: []string{resource.Name},
									}
									graphNodeList = append(graphNodeList, errnode)
								}
							}
						} else {
							errnode := Errornode{
								Name:  "missing " + combinedlist[i],
								Edges: []string{resource.Namespace + ":" + resource.Name},
							}
							fmt.Println("parent name is", resource.Namespace+":"+resource.Name)
							graphNodeList = append(graphNodeList, errnode)
						}
					case *v1alpha1.ResourceReference:
						if thetype != nil {
							val := reflect.Indirect(reflect.ValueOf(thetype))
							_, ok := val.Type().FieldByName("Name")
							if !ok {
								errnode := Errornode{
									Name:  "missing " + combinedlist[i],
									Edges: []string{resource.Name},
								}
								graphNodeList = append(graphNodeList, errnode)
							}
						} else {
							errnode := Errornode{
								Name:  "missing " + combinedlist[i],
								Edges: []string{resource.Namespace + ":" + resource.Name},
							}
							fmt.Println("parent name is", resource.Namespace+":"+resource.Name)
							graphNodeList = append(graphNodeList, errnode)
						}
					}

				}

			}
		}
	}
	return graphNodeList, nil
}
