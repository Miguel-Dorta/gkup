package pkg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type dir struct {
	Name string `json:"name"`
	Dirs []dir `json:"dirs"`
	Files []*file `json:"files"`
}

type file struct {
	Name string `json:"name"`
	Size int64 `json:"size"`
	Hash []byte `json:"hash"`
	realPath string `json:"-"`
}

func listFilesRecursive(path string) (dir, []*file, error) {
	children, err := listDir(path)
	if err != nil {
		return dir{}, nil, fmt.Errorf("cannot list \"%s\": %s", path, err.Error())
	}

	var fileList []*file
	d := dir {
		Name: filepath.Base(path),
		Files: make([]*file, 0, 10),
		Dirs: make([]dir, 0, 10),
	}

	for _, child := range children {
		childPath := filepath.Join(path, child.Name())

		if child.Mode().IsDir() {
			subChild, childFiles, err := listFilesRecursive(childPath)
			if err != nil {
				if OmitErrors {
					os.Stderr.WriteString(err.Error())
					continue
				} else {
					return dir{}, nil, err
				}
			}
			d.Dirs = append(d.Dirs, subChild)
			fileList = append(fileList, childFiles...)
		} else if child.Mode().IsRegular() {
			subChild, err := getFile(childPath)
			if err != nil {
				if OmitErrors {
					os.Stderr.WriteString(err.Error())
					continue
				} else {
					return dir{}, nil, err
				}
			}
			d.Files = append(d.Files, subChild)
		} else {
			// TODO symlinks and other cases
		}
	}

	return d, append(fileList, d.Files...), nil
}

var copyBuf = make([]byte, BufferSize)
func copyFile(origin, destiny string) error {
	originFile, err := os.Open(origin)
	if err != nil {
		return fmt.Errorf("cannot open file \"%s\": %s", origin, err.Error())
	}
	defer originFile.Close()

	destinyFile, err := os.Create(destiny)
	if err != nil {
		return fmt.Errorf("cannot create file in \"%s\": %s", destiny, err.Error())
	}
	defer destinyFile.Close()

	if _, err = io.CopyBuffer(destinyFile, originFile, copyBuf); err != nil {
		var errStr strings.Builder
		errStr.Grow(1000)

		stringBuilderAppend(&errStr,
			"Error copying file from \"", origin, "\" to \"", destiny, "\": ", err.Error(),
			"\n-> DESCRIPTION: ", err.Error(),
			"\n-> CLOSED: ",
		)
		if err = destinyFile.Close(); err == nil {
			errStr.WriteString("yes - REMOVED: ")
			if err = os.Remove(destiny); err == nil {
				errStr.WriteString("yes")
				return errors.New(errStr.String())
			}
		}

		stringBuilderAppend(&errStr, "no\n-> There is a corrupt file in \"", destiny, "\". Please, remove it")
		return errors.New(errStr.String())
	}

	return destinyFile.Close()
}

func listDir(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdir(-1)
}
