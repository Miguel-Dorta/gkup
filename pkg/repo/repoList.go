package repo

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"os"
	"path/filepath"
	"sort"
)

func (r *Repo) ListBackups() error {
	if r.sett == nil {
		return errors.New("settings not loaded")
	}

	pkg.Log.Debug("Listing backup directory")
	dirs, files, err := ListDirSorted(r.backupFolder)
	if err != nil {
		return fmt.Errorf("cannot list backup directory: %s", err.Error())
	}

	if len(files) == 0 && len(dirs) == 0 {
		fmt.Println("Not backup found")
		return nil
	}

	if len(files) != 0 {
		fmt.Println("Unnamed backups:")
		for _, file := range files {
			fmt.Println(" ->", file)
		}
		fmt.Print("\n")
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(r.backupFolder, dir)
		backupsInDir, err := ListDirNamesSorted(dirPath)
		if err != nil {
			return fmt.Errorf("cannot list directory \"%s\": %s", dirPath, err.Error())
		}

		fmt.Println(dir + ":")
		for _, b := range backupsInDir {
			fmt.Println(" ->", b)
		}
		fmt.Print("\n")
	}

	return nil
}

func ListDirNamesSorted(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	sort.Strings(names)

	if err = f.Close(); err != nil {
		return nil, err
	}
	return names, nil
}

func ListDirSorted(path string) (dirs, files []string, err error) {
	dirs = make([]string, 0, 10)
	files = make([]string, 0, 10)

	childs, err := utils.ListDir(path)
	if err != nil {
		return nil, nil, err
	}

	for _, child := range childs {
		if child.Mode().IsDir() {
			dirs = append(dirs, child.Name())
		} else if child.Mode().IsRegular() {
			files = append(files, child.Name())
		}
	}

	sort.Strings(dirs)
	sort.Strings(files)
	return dirs, files, nil
}
