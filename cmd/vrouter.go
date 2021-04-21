package cmd

import (
	"errors"

	"github.com/michaelhenkel/validator/validate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(vrouterCmd)
}

var vrouterCmd = &cobra.Command{
	Use:   "vrouter",
	Short: "validates a vrouter",
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
		vrouterValidator := validate.NewVrouter(ClientConfig)
		vrouterValidator.Validate()
	},
}
