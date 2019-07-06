package repo

import (
	"errors"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

// RestoreBackup restores the backup made in the date provided in the path provided.
func (r *Repo) RestoreBackup(backupName, backupDate, destination string) error {
	if r.sett == nil {
		return errors.New("settings not loaded")
	}

	var b backup
	{
		backupFolder := r.backupFolder
		if backupName != "" {
			backupFolder = filepath.Join(backupFolder, backupName)
		}

		backupsList, err := utils.ListDir(r.backupFolder)
		if err != nil {
			return err
		}

		pkg.Log.Debug("Finding backup")
		var backupPath string
		for _, bac := range backupsList {
			if strings.HasPrefix(bac.Name(), backupDate) {
				backupPath = filepath.Join(r.backupFolder, bac.Name())
				break
			}
		}

		if backupPath == "" {
			return errors.New("backupPath not found")
		}

		pkg.Log.Debug("Reading backup")
		if b, err = readBackup(backupPath); err != nil {
			return err
		}
	}

	pkg.Log.Infof("Restoring backup in %s", destination)
	if err := r.restoreDir(files.Dir{Files: b.Files, Dirs: b.Dirs}, destination, make([]byte, pkg.BufferSize)); err != nil {
		return err
	}
	return nil
}

// restoreDir restores a specific files.Dir in the path provided.
func (r *Repo) restoreDir(d files.Dir, destination string, buffer []byte) error {
	for _, childFile := range d.Files {
		pkg.Log.Debugf("Restoring file %s in %s", childFile.Name, destination)
		if err := utils.CopyFile(r.getPathInRepo(childFile), filepath.Join(destination, childFile.Name), buffer); err != nil {
			if pkg.OmitErrors {
				pkg.Log.Error(err.Error())
				continue
			} else {
				return err
			}
		}
	}

	for _, childDir := range d.Dirs {
		pkg.Log.Debugf("Restoring directory %s in %s", childDir.Name, destination)
		childPath := filepath.Join(destination, childDir.Name)
		if err := os.Mkdir(childPath, 0700); err != nil {
			if pkg.OmitErrors {
				pkg.Log.Errorf("Error restoring folder \"%s\": %s\n", childPath, err.Error())
				continue
			} else {
				return err
			}
		}
		if err := r.restoreDir(childDir, childPath, buffer); err != nil {
			return err
		}
	}

	return nil
}
