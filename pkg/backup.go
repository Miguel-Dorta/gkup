package pkg

import (
	"encoding/json"
	"io/ioutil"
)

type backup struct {
	Version string `json:version`
	Files []file   `json:files`
}

func readBackup(path string) (b backup, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &b)
	return
}

func writeBackup(path string, files []file) error {
	data, err := json.Marshal(backup{Version: Version, Files: files})
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0600)
}
