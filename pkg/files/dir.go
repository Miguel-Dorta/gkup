package files

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/logger"
	"github.com/Miguel-Dorta/gkup/pkg/tmp"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"path/filepath"
)

// Dir represents an abstraction of a directory
type Dir struct {
	Name string   `json:"name"`
	Dirs []Dir    `json:"dirs"`
	Files []*File `json:"files"`
}

// NewDir returns a Dir object that represents the complete structure from the path provided
// and a slice of File objects containing all the files from that structure
func NewDir(path string) (Dir, []*File, error) {
	children, err := utils.ListDir(path)
	if err != nil {
		return Dir{}, nil, fmt.Errorf("cannot list \"%s\": %s", path, err.Error())
	}

	var fileList []*File
	d := Dir{
		Name: filepath.Base(path),
		Files: make([]*File, 0, 10),
		Dirs: make([]Dir, 0, 10),
	}

	for _, child := range children {
		childPath := filepath.Join(path, child.Name())

		if tmp.OmitHidden {
			isHidden, err := utils.IsHidden(childPath, child.Name())
			if err != nil {
				if logger.OmitErrors {
					logger.Log.Errorf("cannot determine if path \"%s\" is hidden: %s", childPath, err.Error())
					continue
				} else {
					return Dir{}, nil, fmt.Errorf("error determining if path \"%s\" is hidden: %s", childPath, err.Error())
				}
			}

			if isHidden {
				logger.Log.Debugf("omitting hidden file %s", childPath)
				continue
			}
		}

		if utils.IsSymLink(child.Mode()) {
			solvedChild, err := utils.ResolveSymlink(childPath)
			if err != nil {
				if logger.OmitErrors {
					logger.Log.Error(err.Error())
					continue
				} else {
					return Dir{}, nil, err
				}
			}
			child = solvedChild
		}

		if child.Mode().IsDir() {
			logger.Log.Debugf("Listing directory %s", childPath)
			subChild, childFiles, err := NewDir(childPath)
			if err != nil {
				if logger.OmitErrors {
					logger.Log.Error(err.Error())
					continue
				} else {
					return Dir{}, nil, err
				}
			}
			d.Dirs = append(d.Dirs, subChild)
			fileList = append(fileList, childFiles...)
		} else if child.Mode().IsRegular() {
			logger.Log.Debugf("Listing file %s", childPath)
			subChild, err := NewFile(childPath)
			if err != nil {
				if logger.OmitErrors {
					logger.Log.Error(err.Error())
					continue
				} else {
					return Dir{}, nil, err
				}
			}
			d.Files = append(d.Files, subChild)
		}
	}

	return d, append(fileList, d.Files...), nil
}
