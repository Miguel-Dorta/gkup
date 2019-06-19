package repo

import (
	"errors"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"github.com/Miguel-Dorta/gkup/pkg/logger"
	"github.com/Miguel-Dorta/gkup/pkg/tmp"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

// RestoreBackup restores the backup made in the date provided in the path provided.
func (r *Repo) RestoreBackup(date, restoreTo string) error {
	if r.sett == nil {
		return errors.New("settings not loaded")
	}

	var b backup
	{
		backupsList, err := utils.ListDir(r.backupFolder)
		if err != nil {
			return err
		}

		logger.Log.Debug("Finding backup")
		var backupPath string
		for _, bac := range backupsList {
			if strings.HasPrefix(bac.Name(), date) {
				backupPath = filepath.Join(r.backupFolder, bac.Name())
				break
			}
		}

		if backupPath == "" {
			return errors.New("backupPath not found")
		}

		logger.Log.Debug("Reading backup")
		if b, err = readBackup(backupPath); err != nil {
			return err
		}
	}

	//TODO check versioning

	logger.Log.Debugf("Restoring backup in %s", restoreTo)
	if err := r.restoreDir(files.Dir{Files: b.Files, Dirs: b.Dirs}, restoreTo); err != nil {
		return err
	}
	return nil
}

// restoreDir restores a specific files.Dir in the path provided.
func (r *Repo) restoreDir(d files.Dir, pathToRestore string) error {
	for _, childFile := range d.Files {
		logger.Log.Debugf("Restoring file %s in %s", childFile.Name, pathToRestore)
		if err := utils.CopyFile(r.getPathInRepo(childFile), filepath.Join(pathToRestore, childFile.Name)); err != nil {
			if tmp.OmitErrors {
				logger.Log.Error(err.Error())
				continue
			} else {
				return err
			}
		}
	}

	for _, childDir := range d.Dirs {
		logger.Log.Debugf("Restoring directory %s in %s", childDir.Name, pathToRestore)
		childPath := filepath.Join(pathToRestore, childDir.Name)
		if err := os.Mkdir(childPath, 0700); err != nil {
			if tmp.OmitErrors {
				logger.Log.Errorf("Error restoring folder \"%s\": %s\n", childPath, err.Error())
				continue
			} else {
				return err
			}
		}
		if err := r.restoreDir(childDir, childPath); err != nil {
			return err
		}
	}

	return nil
}

