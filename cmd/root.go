package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd This represents the binary itself
var RootCmd = &cobra.Command{
	Use:   "phileas",
	Short: "Phileas creates your bucket list",
	Long:  "Phileas creates your bucket list based on your instagram & reddit likes & upvotes, respectively.",
	Run: func(cmd *cobra.Command, args []string) {
		// Empty root command
	},
}
