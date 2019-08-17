package repo

import (
	"encoding/hex"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"path/filepath"
)

const (
	BackupFolderName = "backup"
	FilesFolderName  = "files"
	SettingsFileName = "settings.toml"
)

// Repo is a type that represents the abstract structure of a gkup repo.
// It's the place where the backed up files, the backup info, and the backup settings are saved.
type Repo struct {
	sett         *settings
	path         string
	backupFolder string
	filesFolder  string
	settingsPath string
}

// New creates a new Repo object
func New(repoPath string) *Repo {
	return &Repo{
		path:         repoPath,
		backupFolder: filepath.Join(repoPath, BackupFolderName),
		filesFolder:  filepath.Join(repoPath, FilesFolderName),
		settingsPath: filepath.Join(repoPath, SettingsFileName),
	}
}

// LoadSettings loads the repo settings
func (r *Repo) LoadSettings() error {
	pkg.Log.Debug("Loading settings...")
	sett, err := readSettings(r.settingsPath)
	if err != nil {
		return err
	}
	r.sett = &sett
	return nil
}

// getPathInRepo gets the paths where a specific files.File should be stored in the repo
func (r *Repo) getPathInRepo(f *files.File) string {
	hashStr := hex.EncodeToString(f.Hash)
	return filepath.Join(
		r.filesFolder,
		hashStr[:2],
		fmt.Sprintf("%s-%d", hashStr, f.Size),
	)
}
