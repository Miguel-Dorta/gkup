package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the backups saved in the repo",
	Long: `Print a list of the backups saved in the repository provided.`,
	Run: parseCmd,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
