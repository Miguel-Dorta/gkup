package files

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/tmp"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"os"
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

		if child.Mode().IsDir() {
			subChild, childFiles, err := NewDir(childPath)
			if err != nil {
				if tmp.OmitErrors {
					os.Stderr.WriteString(err.Error())
					continue
				} else {
					return Dir{}, nil, err
				}
			}
			d.Dirs = append(d.Dirs, subChild)
			fileList = append(fileList, childFiles...)
		} else if child.Mode().IsRegular() {
			subChild, err := NewFile(childPath)
			if err != nil {
				if tmp.OmitErrors {
					os.Stderr.WriteString(err.Error())
					continue
				} else {
					return Dir{}, nil, err
				}
			}
			d.Files = append(d.Files, subChild)
		} else {
			// TODO symlinks and other cases
		}
	}

	return d, append(fileList, d.Files...), nil
}