package cmd

import (
	"github.com/michaelhenkel/validator/validate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(printCmd)
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "prints a graph",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			Name = args[0]
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ClientConfig.Name = Name
		validator := validate.NewValidator(ClientConfig)
		validator.Validate()
		validator.Print()
	},
}
