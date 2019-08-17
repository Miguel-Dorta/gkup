package files

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"path/filepath"
)

// Dir represents an abstraction of a directory
type Dir struct {
	Name  string  `json:"name"`
	Dirs  []Dir   `json:"dirs"`
	Files []*File `json:"files"`
}

// NewDir returns a Dir object that represents the complete structure from the path provided
// and a slice of File objects containing all the files from that structure
func NewDir(path string) (Dir, []*File, error) {
	// Check if it's a directory
	children, err := utils.ListDir(path)
	if err != nil {
		return Dir{}, nil, fmt.Errorf("cannot list \"%s\": %s", path, err.Error())
	}

	var fileList []*File
	d := Dir{
		Name:  filepath.Base(path),
		Files: make([]*File, 0, pkg.SliceSmallCapacity),
		Dirs:  make([]Dir, 0, pkg.SliceSmallCapacity),
	}

	for _, child := range children {
		childPath := filepath.Join(path, child.Name())

		// Omit if hidden
		if pkg.OmitHidden && utils.IsHidden(child.Name()) {
			pkg.Log.Debugf("omitting hidden file %s", childPath)
			continue
		}

		if child.Mode().IsDir() { // If child is a directory, list it, and add it to this directory list of directories, and its files to the filelist.
			pkg.Log.Debugf("Listing directory %s", childPath)
			subChild, childFiles, err := NewDir(childPath)
			if err != nil {
				if pkg.OmitErrors {
					pkg.Log.Error(err.Error())
					continue
				} else {
					return Dir{}, nil, err
				}
			}
			d.Dirs = append(d.Dirs, subChild)
			fileList = append(fileList, childFiles...)

		} else if child.Mode().IsRegular() { // If child is a file, add it to this directory list of files
			pkg.Log.Debugf("Listing file %s", childPath)
			subChild, err := NewFile(childPath)
			if err != nil {
				if pkg.OmitErrors {
					pkg.Log.Error(err.Error())
					continue
				} else {
					return Dir{}, nil, err
				}
			}
			d.Files = append(d.Files, subChild)

		} else { // If child is neither a directory nor a file, omit it
			pkg.Log.Debugf("omitting unsupported file %s", childPath)
		}
	}

	return d, append(fileList, d.Files...), nil
}
