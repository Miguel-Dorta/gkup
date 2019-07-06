package repo

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"io/ioutil"
)

// backup is a type for saving the files and directories that are backed up.
// It is intended to be saved in json format
type backup struct {
	Version string        `json:"version"`
	Dirs    []files.Dir   `json:"dirs"`
	Files   []*files.File `json:"files"`
}

// readBackup reads and parses the backup from the path provided
func readBackup(path string) (backup, error) {
	var b backup

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return b, fmt.Errorf("cannot read backup file \"%s\": %s", path, err.Error())
	}

	if err = json.Unmarshal(data, &b); err != nil {
		return backup{}, fmt.Errorf("error parsing backup: %s", err.Error())
	}
	return b, nil
}

// writeBackup writes the backup provided in the path provided
func writeBackup(path string, b backup) error {
	data, _ := json.Marshal(b)
	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("cannot write backup to \"%s\": %s", path, err.Error())
	}

	return nil
}
