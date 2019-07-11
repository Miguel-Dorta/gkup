package cmd

import (
	"github.com/Miguel-Dorta/gkup/internal"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gkup",
	Short: "gkup is a backup tool that reduces file redundancy",
	Long: `gkup it's a backup tool which aim is to reduce file redundancy to a minimum.
This will save a lot of space in your backup drives, avoiding multiple copies
of the same files over and over.

gkup's backups are also designed to be human-readable, so they'll be easy to
restore even if all the copies of this program are erased of the surface of
the Earth, and they'll also easily parseable by other programs.

gkup is not aimed to provide any kind of compression, encryption or redundancy.
It's the user's' responsibility to do this if they feel they wanted.`,
	Version: internal.Version,
}

// Parse adds all child commands to the root command and sets flags appropriately.
func Parse() error {
	return rootCmd.Execute()
}

func init() {
	addPersistentFlagRepoPath(rootCmd)
	addPersistentFlagVerboseLevel(rootCmd)
	addPersistentFlagOmitErrors(rootCmd)
}
