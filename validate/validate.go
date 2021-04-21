package validate

import (
	"github.com/michaelhenkel/validator/k8s/clientset"
)

type Validator struct {
	Name      string
	Namespace string
}

type ClientConfig struct {
	Name      string
	Namespace string
	Client    *clientset.Client
}

type ValidatorInterface interface {
	Validate()
}
