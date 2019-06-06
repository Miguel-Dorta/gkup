package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type backup struct {
	Version string `json:version`
	Dirs []dir `json:dirs`
	Files []file   `json:files`
}

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

func writeBackup(path string, d dir) error {
	data, err := json.Marshal(backup{
		Version: Version,
		Dirs: d.Dirs,
		Files: d.Files,
	})
	if err != nil {
		return fmt.Errorf("cannot compose backup: %s", err.Error())
	}

	if err = ioutil.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("cannot write backup to \"%s\": %s", path, err.Error())
	}

	return nil
}
