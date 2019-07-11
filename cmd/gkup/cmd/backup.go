package cmd

import (
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a new backup from the paths provided",
	Long: `backup will create a new backup with the current date and the name provided in
the repository.`,
	Run: parseCmd,
}

func init() {
	rootCmd.AddCommand(backupCmd)

	addFlagBackupName(backupCmd)
	addFlagBufferSize(backupCmd)
	addFlagNumberOfThreads(backupCmd)
	addFlagOmitHidden(backupCmd)
	addFlagReadSymLinks(backupCmd)
}
