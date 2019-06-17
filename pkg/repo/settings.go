package repo

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

// settings is a type for saving the general settings of a repo.
// It is intended to be saved in toml format
type settings struct {
	Version       string `toml:"version"`
	HashAlgorithm string `toml:"hashAlgorithm"`
}

// readSettings reads and parses the settings from the paths provided
func readSettings(path string) (sett settings, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return sett, fmt.Errorf("cannot read settings in \"%s\": %s", path, err.Error())
	}

	if err = toml.Unmarshal(data, &sett); err != nil {
		return settings{}, fmt.Errorf("cannot parse settings from \"%s\": %s", path, err.Error())
	}
	return sett, nil
}

// writeSettings writes the settings provided in the path provided
func writeSettings(path string, sett settings) error {
	data, _ := toml.Marshal(sett)

	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("cannot write settings in \"%s\": %s", path, err.Error())
	}
	return nil
}
