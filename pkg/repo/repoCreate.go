package repo

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"os"
	"path/filepath"
)

// Create creates the structure for the repo in the path provided
func (r *Repo) Create(hashAlgorithm string) error {
	// Check if it's a directory
	if stat, err := os.Stat(r.path); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error getting stats from \"%s\": %s", r.path, err.Error())
		}
		// If not exists, create it
		if err = os.MkdirAll(r.path, 0700); err != nil {
			return fmt.Errorf("error creating directory \"%s\": %s", r.path, err.Error())
		}
	} else if !stat.IsDir() {
		//TODO read symlink
		return fmt.Errorf("\"%s\" is not a directory", r.path)
	}

	// Check if it's empty
	if children, err := utils.ListDir(r.path); err != nil {
		return fmt.Errorf("error listing \"%s\": %s", r.path, err.Error())
	} else if len(children) != 0 {
		return fmt.Errorf("\"%s\" is not empty", r.path)
	}

	// Make backup folder
	if err := os.Mkdir(r.backupFolder, 0700); err != nil {
		return fmt.Errorf("error creating subdirectory \"%s\": %s", r.backupFolder, err.Error())
	}

	// Make file folder and subFolders
	if err := os.Mkdir(r.filesFolder, 0700); err != nil {
		return fmt.Errorf("error creating subdirectory \"%s\": %s", r.filesFolder, err.Error())
	}
	for i:=0x0; i<=0xff; i++ {
		path := filepath.Join(r.filesFolder, fmt.Sprintf("%02x", i))
		if err := os.Mkdir(path, 0700); err != nil {
			return fmt.Errorf("error creating subdirectory \"%s\": %s", path, err.Error())
		}
	}

	// Write settings.toml
	if err := writeSettings(r.settingsPath, settings{pkg.Version, hashAlgorithm}); err != nil {
		return err
	}
	return nil
}
