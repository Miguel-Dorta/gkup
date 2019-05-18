package pkg

import (
	"io"
	"os"
	"path/filepath"
)

type file struct {
	Path string `json:path`
	Size int64 `json:size`
	Hash []byte `json:hash`
}

func listFilesRecursive(path string) ([]file, error) {
	childs, err := listDir(path)
	if err != nil {
		return nil, err //TODO
	}

	files := make([]file, 0, 10)
	for _, child := range childs {
		childMode := child.Mode()
		childPath := filepath.Join(path, child.Name())

		if childMode.IsDir() {
			subChilds, err := listFilesRecursive(childPath)
			if err != nil {
				return nil, err //TODO
			}
			files = append(files, subChilds...)
		} else if childMode.IsRegular() {
			f, err := getFile(childPath)
			if err != nil {
				//TODO
			}
			files = append(files, f)
		} else {
			//TODO
		}
	}

	return files, nil
}

func getFile(path string) (f file, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return
	}
	f.Size = stat.Size()

	f.Hash, err = hashFile(path)
	if err != nil {
		return
	}

	f.Path = path
	return
}

var copyBuf = make([]byte, BufferSize)
func copyFile(origin, destiny string) (err error) {
	originFile, err := os.Open(origin)
	if err != nil {
		return
	}
	defer originFile.Close()

	destinyFile, err := os.Create(destiny)
	if err != nil {
		return
	}
	defer destinyFile.Close()

	_, err = io.CopyBuffer(destinyFile, originFile, copyBuf)
	return
}

func listDir(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdir(-1)
}
