package cmd

import (
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a backup in a directory.",
	Long: `Restore a backup that matches the given name (if provided) and date in the
directory specified.`,
	Run: parseCmd,
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	addFlagBackupName(restoreCmd)
	addFlagBufferSize(restoreCmd)
	addFlagBackupDate(restoreCmd)
	if err := restoreCmd.MarkFlagRequired("date"); err != nil {
		panic(err)
	}
}
