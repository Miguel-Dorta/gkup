package pkg

import (
	"encoding/json"
	"io/ioutil"
)

type backup struct {
	Version string `json:version`
	Dirs []dir `json:dirs`
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

func writeBackup(path string, d dir) error {
	data, err := json.Marshal(backup{
		Version: Version,
		Dirs: d.Dirs,
		Files: d.Files,
	})
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0600)
}
