package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var charmCmd = &cobra.Command{
	Use:   "charm",
	Short: "Charm  management",
	Long:  `Builds and publishes charms`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("charm called")
	},
}

func init() {
	rootCmd.AddCommand(charmCmd)
}
