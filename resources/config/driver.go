package resources

import (
	"fmt"
	"reflect"
	"strings"

	contrailcorev1alpha1 "ssd-git.juniper.net/contrail/cn2/contrail/pkg/apis/core/v1alpha1"
)

func getspecvals() map[string][]string {
	ret := make(map[string][]string)

	vmi := &contrailcorev1alpha1.VirtualMachineInterface{}
	val := reflect.Indirect(reflect.ValueOf(vmi))
	t := val.Type().String()
	specField, ok := val.Type().FieldByName("Spec")
	if !ok {
		fmt.Println("no spec field")
	}
	specVal := reflect.Indirect(reflect.ValueOf(specField))
	specFieldNum := reflect.TypeOf(specVal).NumField()
	for i := 0; i < specFieldNum; i++ {
		f := specField.Type.Field(i)
		referencesvisited := false
		refsFieldList := strings.Split(f.Name, "References")
		if len(refsFieldList) > 1 {
			referencesvisited = true
			ret["references"] = append(ret["references"], f.Name)
			fmt.Println("found references to ", refsFieldList[0])
			val := reflect.Indirect(reflect.ValueOf(f.Type))
			fmt.Println("Name is:", val)
		}
		refFieldList := strings.Split(f.Name, "Reference")
		if len(refFieldList) > 1 && !referencesvisited {
			ret["reference"] = append(ret["references"], f.Name)
			fmt.Println("found reference to ", refFieldList[0])
			fmt.Println("Name is:", reflect.Indirect(reflect.ValueOf(f)))
		}
		parentFieldList := strings.Split(f.Name, "Parent")
		if len(parentFieldList) > 1 {
			ret["parent"] = append(ret["references"], f.Name)
			fmt.Println("found Parent to ", parentFieldList[0])
			val := f.Tag
			fmt.Println("Name is:", val)
		}
	}
	fmt.Println("SpecField", specField)
	fmt.Println("Type", t)
	return ret
}

func getstatusvals() map[string][]string {
	ret := make(map[string][]string)
	vmi := &contrailcorev1alpha1.VirtualMachineInterface{}
	val := reflect.Indirect(reflect.ValueOf(vmi))
	t := val.Type().String()
	statusField, ok := val.Type().FieldByName("Status")
	if !ok {
		fmt.Println("no status field")
	}
	statusval := reflect.Indirect(reflect.ValueOf(statusField))
	statusFieldNum := reflect.TypeOf(statusval).NumField()

	for i := 0; i < statusFieldNum; i++ {
		f := statusField.Type.Field(i)
		referencesvisited := false
		refsFieldList := strings.Split(f.Name, "References")
		if len(refsFieldList) > 1 {
			referencesvisited = true

			ret["references"] = append(ret["references"], f.Name)
			fmt.Println("found references to ", refsFieldList[0])
			fmt.Println("Name is:", reflect.Indirect(reflect.ValueOf(f)))
		}
		refFieldList := strings.Split(f.Name, "Reference")
		if len(refFieldList) > 1 && !referencesvisited {
			ret["reference"] = append(ret["references"], f.Name)
			fmt.Println("found reference to ", refFieldList[0])
			fmt.Println("Name is:", reflect.Indirect(reflect.ValueOf(f)))
		}
		parentFieldList := strings.Split(f.Name, "Parent")
		if len(parentFieldList) > 1 {
			ret["Parent"] = append(ret["references"], f.Name)
			fmt.Println("found Parent to ", parentFieldList[0])
			fmt.Println("Name is:", reflect.Indirect(reflect.ValueOf(f)))
		}
	}
	fmt.Println("SpecField", statusField)
	fmt.Println("Type", t)
	return ret
}
