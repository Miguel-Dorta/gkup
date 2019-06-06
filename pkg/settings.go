package pkg

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

const (
	settingsFilename = "settings.toml"
)

var (
	Version string
	ReadSymLink = false
	Verbose = false
	BufferSize = 4*1024*1024
	OmitErrors = false
)

type settings struct {
	Version       string `toml:version`
	HashAlgorithm string `toml:hashAlgorithm`
}

func saveSettings(path, hashAlgorithm string) error {
	data, _ := toml.Marshal(settings{
		Version:       Version,
		HashAlgorithm: hashAlgorithm,
	})

	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("cannot write settings in \"%s\": %s", path, err.Error())
	}
	return nil
}

func loadSettings(path string) (sett settings, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return sett, fmt.Errorf("cannot read settings in \"%s\": %s", path, err.Error())
	}

	if err = toml.Unmarshal(data, &sett); err != nil {
		return settings{}, fmt.Errorf("cannot parse settings from \"%s\": %s", path, err.Error())
	}
	return sett, nil
}
