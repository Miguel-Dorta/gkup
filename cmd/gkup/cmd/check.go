package cmd

import (
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the integrity of your repository",
	Long: `Check the integrity of the files in your repository. It will detect if the files
are corrupted or have defects from bad copying or bad hardware.`,
	Run: parseCmd,
}

func init() {
	rootCmd.AddCommand(checkCmd)

	addFlagBufferSize(checkCmd)
	addFlagNumberOfThreads(checkCmd)
}
