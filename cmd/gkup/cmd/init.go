package cmd

import (
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize repository where backups will be stored",
	Long:  `init initializes the gkup repository where backups will be stored.`,
	Run:   parseCmd,
}

func init() {
	rootCmd.AddCommand(initCmd)

	addFlagSum(initCmd)
}
