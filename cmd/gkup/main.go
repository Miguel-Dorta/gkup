package main

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(
`Usage:    ./gkup <create/check/backup/restore> <repo-path> [optional-args]
Optional args:
    create - <sha256/sha1>
    backup - <path-to-backup>
    restore - <which-backup-restore> <path-to-restore>`,
		)
		os.Exit(0)
	}

	repo := pkg.NewRepo(os.Args[2])
	switch os.Args[1] {
	case "create":
		hashAlg := "sha256"
		if len(os.Args) == 4 {
			hashAlg = os.Args[3]
		}
		err := repo.Create(hashAlg)
		if err != nil {
			fmt.Println(err.Error())
		}
		break
	case "check":
		err := repo.LoadSettings()
		if err != nil {
			fmt.Println(err.Error())
		}
		errs := repo.CheckIntegrity()
		fmt.Printf("Errors found: %d\n", errs)
		break
	case "backup":
		err := repo.LoadSettings()
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(os.Args) != 4 {
			fmt.Println("Nothing to backup - aborting!")
		}
		err = repo.BackupPaths(os.Args[3:])
		if err != nil {
			fmt.Println(err.Error())
		}
		break
	case "restore":
		err := repo.LoadSettings()
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(os.Args) != 5 {
			fmt.Println("Insufficient arguments - aborting!")
		}
		err = repo.RestoreBackup(os.Args[3], os.Args[4])
		if err != nil {
			fmt.Println(err.Error())
		}
		break
	default:
		fmt.Printf("Command not found!")
	}
}
