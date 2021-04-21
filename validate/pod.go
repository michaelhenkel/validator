package validate

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodValidator struct {
	clientConfig *ClientConfig
}

func NewPod(clientConfig *ClientConfig) *PodValidator {
	v := &PodValidator{
		clientConfig: clientConfig,
	}
	return v
}

func (v *PodValidator) Validate() {
	coreV1 := v.clientConfig.Client.CoreV1
	name := v.clientConfig.Name
	namespace := v.clientConfig.Namespace
	pod, err := coreV1.Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("pod %s/%s\n %s\n", namespace, name, pod)
}
