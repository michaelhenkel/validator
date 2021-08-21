package cmd

import (
	"github.com/s3kim2018/validator/k8s/clientset"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "validator",
		Short: "validates contrail resources",
	}
	Namespace  string
	Kubeconfig string
	Client     *clientset.Client
)

func initConfig() {
	if Namespace == "" {
		Namespace = "default"
	}
	client, err := clientset.NewClient(Kubeconfig)
	if err != nil {
		panic(err)
	}
	Client = client
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&Kubeconfig, "kubeconfig", "k", "", "path to kubeconfig")
}
