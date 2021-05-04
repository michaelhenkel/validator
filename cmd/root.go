package cmd

import (
	"github.com/michaelhenkel/validator/k8s/clientset"
	"github.com/michaelhenkel/validator/validate"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "validator",
		Short: "validates contrail resources",
	}
	Namespace    string
	Kubeconfig   string
	Name         string
	ClientConfig *validate.ClientConfig
)

func initConfig() {
	if Namespace == "" {
		Namespace = "default"
	}
	client, err := clientset.NewClient(Kubeconfig)
	if err != nil {
		panic(err)
	}
	ClientConfig = &validate.ClientConfig{
		Name:   Name,
		Client: client,
	}
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "", "resource namespace")
	rootCmd.PersistentFlags().StringVarP(&Kubeconfig, "kubeconfig", "k", "", "path to kubeconfig")
}
