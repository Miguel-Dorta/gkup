package main

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/cmd/gkup/cmd"
	"github.com/Miguel-Dorta/gkup/internal"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/repo"
	"os"
	"time"
)

func init() {
	pkg.Version = internal.Version
	pkg.Log.Formatter = func(levelName, msg string) string {
		now := time.Now()
		return fmt.Sprintf(
			"[%04d-%02d-%02d %02d:%02d:%02d] %s: %s",
			now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second(),
			levelName, msg,
		)
	}
}

func main() {
	if err := cmd.Parse(); err != nil {
		pkg.Log.Criticalf("Error parsing commands: %s", err.Error())
		os.Exit(1)
	}

	if len(cmd.ArgsErrors) != 0 {
		for _, err := range cmd.ArgsErrors {
			pkg.Log.Critical(err.Error())
		}
		os.Exit(1)
	}

	pkg.BufferSize      = cmd.BufferSize
	pkg.NumberOfThreads = cmd.NumberOfThreads
	pkg.OmitErrors      = cmd.OmitErrors
	pkg.Log.Level       = cmd.VerboseLevel

	r := repo.New(cmd.RepoPath)
	switch cmd.Cmd {
	case "backup":
		if len(cmd.Args) == 0 {
			pkg.Log.Critical("No files to backup. Skipping empty backup.")
			os.Exit(1)
		}

		if err := r.LoadSettings(); err != nil {
			pkg.Log.Criticalf("Error loading repo settings: %s", err.Error())
			os.Exit(1)
		}

		if err := r.BackupPaths(cmd.Args, cmd.BackupName, cmd.OmitHidden, cmd.ReadSymLinks); err != nil {
			pkg.Log.Criticalf("Error while backing up files: %s", err.Error())
			os.Exit(1)
		}
	case "check":
		if err := r.LoadSettings(); err != nil {
			pkg.Log.Criticalf("Error loading repo settings: %s", err.Error())
			os.Exit(1)
		}

		if err := r.CheckIntegrity(); err != nil {
			pkg.Log.Criticalf("Errors found while checking repo: %s", err.Error())
			os.Exit(1)
		}
	case "init":
		if err := r.Create(cmd.Sum); err != nil {
			pkg.Log.Criticalf("Error initializing repository: %s", err.Error())
			os.Exit(1)
		}
	case "list":
		if err := r.ListBackups(); err != nil {
			pkg.Log.Criticalf("Error listing backups: %s", err.Error())
			os.Exit(1)
		}
	case "restore":
		if len(cmd.Args) == 0 {
			pkg.Log.Critical("Destination path not provided.")
			os.Exit(1)
		}
		if len(cmd.Args) > 1 {
			pkg.Log.Critical("More than one path to restore provided.")
			os.Exit(1)
		}

		if err := r.LoadSettings(); err != nil {
			pkg.Log.Criticalf("Error loading repo settings: %s", err.Error())
			os.Exit(1)
		}

		if err := r.RestoreBackup(cmd.BackupName, cmd.BackupDate, cmd.Args[0]); err != nil {
			pkg.Log.Criticalf("Error restoring backup: %s", err.Error())
			os.Exit(1)
		}
	default:
		fmt.Printf("gkup: %s: command not found\n", cmd.Cmd)
	}
}
