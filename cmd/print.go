package cmd

import (
	"github.com/michaelhenkel/validator/builder"
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
		g := builder.BuildGraph(Client)
		g.String()
		g.EdgeMatcher()
	},
}
