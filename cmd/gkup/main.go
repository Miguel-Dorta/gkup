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
    restore - <which-backup-restore>`,
		)
		os.Exit(0)
	}
	pkg.RepoPath = os.Args[2]
	switch os.Args[1] {
	case "create":
		if len(os.Args) == 4 {
			pkg.HashAlgorithm = os.Args[3]
		}
		err := pkg.CreateRepo()
		if err != nil {
			fmt.Println(err.Error())
		}
		break
	case "check":
		errs := pkg.CheckIntegrity()
		if len(errs) != 0 {
			fmt.Printf("%+v\n", errs)
		}
		break
	case "backup":
		if len(os.Args) != 4 {
			fmt.Println("Nothing to backup - aborting!")
		}
		err := pkg.BackupPaths([]string{os.Args[3]})
		if err != nil {
			panic(err)
		}
		break
	case "restore":
		if len(os.Args) != 4 {
			fmt.Println("Nothing to restore - aborting!")
		}

		break
	default:
		fmt.Printf("Command not found!")
	}
}
