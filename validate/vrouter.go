package validate

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VrouterValidator struct {
	clientConfig *ClientConfig
}

func NewVrouter(clientConfig *ClientConfig) *VrouterValidator {
	v := &VrouterValidator{
		clientConfig: clientConfig,
	}
	return v
}

func (v *VrouterValidator) Validate() {
	contrailV1 := v.clientConfig.Client.ContrailV1
	name := v.clientConfig.Name
	vrouter, err := contrailV1.VirtualRouters().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("pod %s\n %s\n", name, vrouter)
}
