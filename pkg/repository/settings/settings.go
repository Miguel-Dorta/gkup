package settings

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup/internal"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
)

const FileName = "settings.toml"

type Settings struct {
	Version string `toml:"version"`
	HashAlgorithm string `toml:"hash_algorithm"`
}

func Read(path string) (Settings, error) {
	// Read data
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Settings{}, &os.PathError{
			Op: "read settings",
			Path: path,
			Err: err,
		}
	}

	// Parse settings
	s := Settings{}
	if err := toml.Unmarshal(data, &s); err != nil {
		return Settings{}, fmt.Errorf("error parsing settings: %s", err)
	}

	// Check if the fields have info
	if s.Version == "" || s.HashAlgorithm == "" {
		return Settings{}, errors.New("incomplete information in settings")
	}
	return s, nil
}

func Write(path, hashAlgorithm string) error {
	// Serialize settings
	data, err := toml.Marshal(&Settings{
		Version:       internal.Version,
		HashAlgorithm: hashAlgorithm,
	})
	if err != nil {
		return fmt.Errorf("error serializing settings: %s", err)
	}

	// Write data
	if err := ioutil.WriteFile(path, data, pkg.DefaultFilePerm); err != nil {
		return &os.PathError{
			Op:   "write settings",
			Path: path,
			Err:  err,
		}
	}
	return nil
}
