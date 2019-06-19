package repo

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/logger"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"os"
	"path/filepath"
)

// Create creates the structure for the repo in the path provided
func (r *Repo) Create(hashAlgorithm string) error {
	logger.Log.Infof("Creating directory in %s", r.path)
	logger.Log.Debug("Checking if exists something in the repo path")
	stat, err := os.Stat(r.path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error getting stats from \"%s\": %s", r.path, err.Error())
		}

		logger.Log.Debug("Creating repo folder")
		if err = os.MkdirAll(r.path, 0700); err != nil {
			return fmt.Errorf("error creating directory \"%s\": %s", r.path, err.Error())
		}
	}

	logger.Log.Debug("Checking if it's a directory")
	if !stat.IsDir() {
		//TODO read symlink
		return fmt.Errorf("\"%s\" is not a directory", r.path)
	}

	logger.Log.Debug("Checking if it's empty")
	if children, err := utils.ListDir(r.path); err != nil {
		return fmt.Errorf("error listing \"%s\": %s", r.path, err.Error())
	} else if len(children) != 0 {
		return fmt.Errorf("\"%s\" is not empty", r.path)
	}

	logger.Log.Debug("Creating backup folder")
	if err := os.Mkdir(r.backupFolder, 0700); err != nil {
		return fmt.Errorf("error creating subdirectory \"%s\": %s", r.backupFolder, err.Error())
	}

	logger.Log.Debug("Creating files folder")
	if err := os.Mkdir(r.filesFolder, 0700); err != nil {
		return fmt.Errorf("error creating subdirectory \"%s\": %s", r.filesFolder, err.Error())
	}
	for i:=0x0; i<=0xff; i++ {
		path := filepath.Join(r.filesFolder, fmt.Sprintf("%02x", i))
		logger.Log.Debugf("Creating subfolder %s", path)
		if err := os.Mkdir(path, 0700); err != nil {
			return fmt.Errorf("error creating subdirectory \"%s\": %s", path, err.Error())
		}
	}

	logger.Log.Debug("Creating settings.toml")
	if err := writeSettings(r.settingsPath, settings{pkg.Version, hashAlgorithm}); err != nil {
		return err
	}
	return nil
}
