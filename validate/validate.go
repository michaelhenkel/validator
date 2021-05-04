package validate

import (
	"fmt"

	"github.com/michaelhenkel/validator/graph"
	"github.com/michaelhenkel/validator/k8s/clientset"
)

type Validator struct {
	clientConfig *ClientConfig
	graph        *graph.Graph
}

type ClientConfig struct {
	Name   string
	Client *clientset.Client
}

func NewValidator(clientConfig *ClientConfig) *Validator {
	v := &Validator{
		clientConfig: clientConfig,
		graph:        graph.NewGraph(),
	}
	return v
}

func (v *Validator) Validate() error {

	if err := addVirtualRouterNodes(v); err != nil {
		return err
	}

	if err := addVrouterNodes(v); err != nil {
		return err
	}

	if err := addControlNodes(v); err != nil {
		return err
	}

	if err := addBGPRouterNodes(v); err != nil {
		return err
	}

	if err := addPodNodes(v, graph.Vrouter); err != nil {
		return err
	}

	if err := addPodNodes(v, graph.Control); err != nil {
		return err
	}

	if err := addConfigMapNodes(v, graph.Vrouter); err != nil {
		return err
	}

	if err := addConfigMapNodes(v, graph.Control); err != nil {
		return err
	}

	if err := addConfigFileNodes(v); err != nil {
		return err
	}

	if err := addVrouterToPodEdges(v); err != nil {
		return err
	}

	if err := addControlToPodEdges(v); err != nil {
		return err
	}

	if err := addPodToVirtualRouterEdges(v); err != nil {
		return err
	}

	if err := addPodToBGPRouterEdges(v); err != nil {
		return err
	}

	if err := addConfigMapToPodEdges(v, graph.Vrouter); err != nil {
		return err
	}

	if err := addConfigMapToPodEdges(v, graph.Control); err != nil {
		return err
	}

	if err := addConfigFileToConfigMapEdges(v); err != nil {
		return err
	}

	if err := addConfigFileToVirtualRouterBGPRouterEdges(v); err != nil {
		return err
	}

	if err := addConfigFileVirtualRouterToBGPRouterEdges(v); err != nil {
		return err
	}

	return nil
}

func (v *Validator) Print() {
	fmt.Printf("%s\n", v.graph.String())
}
