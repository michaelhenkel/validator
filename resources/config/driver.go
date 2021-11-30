package resources

import (
	"fmt"
	"reflect"
	"strings"

	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

/**
Struct used to extract value of the parent field, spec field, and the status field
They are extracted as string values.
**/
type sourcecoderesource struct {
	Parents    []string
	References []string
	Reference  []string
}

/**
Instance Method to sourcecoderesource.
Populates the sourcecode pointer with Parents, References, and Reference.
@Params None
@Return None
**/
func (source *sourcecoderesource) getspecandstatusvals(resource string) {
	var vmi interface{}

	switch resource {
	case "virtualmachineinterface":
		vmi = &contrailcorev1alpha1.VirtualMachineInterface{}
	case "routinginstance":
		vmi = &contrailcorev1alpha1.RoutingInstance{}
	}

	val := reflect.Indirect(reflect.ValueOf(vmi))
	// t := val.Type().String()
	specField, ok := val.Type().FieldByName("Spec")
	if !ok {
		fmt.Println("no spec field")
	}
	spectype := specField.Type
	specFieldNum := spectype.NumField()
	for i := 0; i < specFieldNum; i++ {
		f := specField.Type.Field(i)
		referencesvisited := false
		refsFieldList := strings.Split(f.Name, "References")
		if len(refsFieldList) > 1 {
			referencesvisited = true
			// fmt.Println("References Found in driver", refsFieldList[0])
			source.References = append(source.References, refsFieldList[0])
		}
		refFieldList := strings.Split(f.Name, "Reference")
		if len(refFieldList) > 1 && !referencesvisited {
			// fmt.Println("Reference Found in driver", refFieldList[0])

			source.Reference = append(source.Reference, refFieldList[0])
		}
		parentFieldList := strings.Split(f.Name, "Parent")
		if len(parentFieldList) > 1 {
			// fmt.Println("Parent Found in driver", parentFieldList[0])

			source.Parents = append(source.Parents, "Parent")
		}
	}
	statusField, ok := val.Type().FieldByName("Status")
	statusFieldNum := statusField.Type.NumField()
	for i := 0; i < statusFieldNum; i++ {
		f := statusField.Type.Field(i)
		referencesvisited := false
		refsFieldList := strings.Split(f.Name, "References")
		if len(refsFieldList) > 1 {
			referencesvisited = true
			// fmt.Println("References Found in driver", refsFieldList[0])
			source.References = append(source.References, refsFieldList[0])
		}
		refFieldList := strings.Split(f.Name, "Reference")
		if len(refFieldList) > 1 && !referencesvisited {
			// fmt.Println("Reference Found in driver", refFieldList[0])
			source.Reference = append(source.Reference, refFieldList[0])
		}
		parentFieldList := strings.Split(f.Name, "Parent")
		if len(parentFieldList) > 1 {
			// fmt.Println("Parent Found in driver", parentFieldList[0])
			source.Parents = append(source.Parents, parentFieldList[0])
		}
	}
}
