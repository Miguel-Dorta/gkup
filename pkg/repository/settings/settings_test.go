package settings_test

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/internal"
	"github.com/Miguel-Dorta/gkup/pkg/repository/settings"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	internal.Version = "v1.0.0-alpha+20191231235959"
	path := filepath.Join(os.TempDir(), fmt.Sprintf("gkup_pkg_repository_settings_TestWrite_%d.toml", time.Now().UnixNano()))
	defer os.Remove(path)

	if err := settings.Write(path, "sha256"); err != nil {
		t.Fatalf("write error in path %s: %s", path, err)
	}
}

func TestRead(t *testing.T) {
	// Check valid
	if err := checkReadValid("testdata/correct.toml", "v1.0.0-alpha+20191231235959", "sha256"); err != nil {
		t.Errorf("error in testdata/correct.toml: %s", err)
	}
	if err := checkReadValid("testdata/extrainfo.toml", "1.0.0", "md5"); err != nil {
		t.Errorf("error in testdata/extrainfo.toml: %s", err)
	}

	// Check invalid
	if sett := checkReadInvalid("testdata/empty.toml"); sett != nil {
		t.Errorf("not error in testdata/empty.toml: %+v", sett)
	}
	if sett := checkReadInvalid("testdata/invalid.toml"); sett != nil {
		t.Errorf("not error in testdata/invalid.toml: %+v", sett)
	}
	if sett := checkReadInvalid("testdata/lacks_info.toml"); sett != nil {
		t.Errorf("not error in testdata/lacks_info.toml: %+v", sett)
	}
	if sett := checkReadInvalid("testdata/non_existing.toml"); sett != nil {
		t.Errorf("not error in testdata/non_existing.toml: %+v", sett)
	}
}

// checkReadValid returns nil if the file is valid and matches the inputs, error otherwise
func checkReadValid(path, version, hashAlgorithm string) error {
	sett, err := settings.Read(path)
	if err != nil {
		return fmt.Errorf("read error: %s", err)
	}

	if sett.Version != version {
		return fmt.Errorf("versions don't match: expected %s, found %s", version, sett.Version)
	}
	if sett.HashAlgorithm != hashAlgorithm {
		return fmt.Errorf("hash_algorithms don't match: expected %s, found %s", hashAlgorithm, sett.HashAlgorithm)
	}
	return nil
}

// checkReadInvalid returns nil if the file is invalid, a Settings object otherwise
func checkReadInvalid(path string) *settings.Settings {
	if sett, err := settings.Read(path); err == nil {
		return &sett
	}
	return nil
}
