package main

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/logger"
	"github.com/Miguel-Dorta/gkup/pkg/repo"
	"github.com/Miguel-Dorta/gkup/pkg/version"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(
`Usage:    ./gkup <create/check/backup/restore/version> <repo-path> [optional-args]
Optional args:
    create - <sha256/sha1>
    backup - <path-to-backup>
    restore - <which-backup-restore> <path-to-restore>`,
		)
		os.Exit(0)
	}

	logger.OmitErrors = true
	repotiar := repo.New(os.Args[2])
	switch os.Args[1] {
	case "create":
		hashAlg := "sha256"
		if len(os.Args) == 4 {
			hashAlg = os.Args[3]
		}
		err := repotiar.Create(hashAlg)
		if err != nil {
			fmt.Println(err.Error())
		}
		break
	case "check":
		err := repotiar.LoadSettings()
		if err != nil {
			fmt.Println(err.Error())
		}
		errs := repotiar.CheckIntegrity(4*1024*1024)
		fmt.Printf("Errors found: %d\n", errs)
		break
	case "backup":
		err := repotiar.LoadSettings()
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(os.Args) != 4 {
			fmt.Println("Nothing to backup - aborting!")
		}
		err = repotiar.BackupPaths(os.Args[3:], 4*1024*1024)
		if err != nil {
			fmt.Println(err.Error())
		}
		break
	case "restore":
		err := repotiar.LoadSettings()
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(os.Args) != 5 {
			fmt.Println("Insufficient arguments - aborting!")
		}
		err = repotiar.RestoreBackup(os.Args[3], os.Args[4], 4*1024*1024)
		if err != nil {
			fmt.Println(err.Error())
		}
		break
	case "version":
		fmt.Printf("Gkup version: %s\n", version.String(version.GkupVersion))
	default:
		fmt.Printf("Command not found!")
	}
}
