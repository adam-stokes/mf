package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Upstream repo management",
	Long:  `Manage upstream git repos, mainly keeping forks updated and code synced.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("repo called")
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
