package cmd

import (
	"errors"

	"github.com/michaelhenkel/validator/validate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(podCmd)
}

var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "validates a pod",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a name argument")
		}
		Name = args[0]
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ClientConfig.Name = Name
		podValidator := validate.NewPod(ClientConfig)
		podValidator.Validate()
	},
}
