package pkg

import (
	"io"
	"os"
	"path/filepath"
)

type dir struct {
	Name string `json:name`
	Dirs []dir `json:dirs`
	Files []file `json:files`
}

type file struct {
	Name string `json:name`
	Size int64 `json:size`
	Hash []byte `json:hash`
	realPath string //Private so it won't be saved in the backup.json
}

func listFilesRecursive(path, name string) (dir, error) {
	childs, err := listDir(path)
	if err != nil {
		return dir{}, err //TODO
	}

	d := dir {
		Name: name,
		// Avoid unnecessary slice resizing
		Files: make([]file, 0, 10),
		Dirs: make([]dir, 0, 10),
	}

	for _, child := range childs {
		// Avoid unnecessary function calls
		childMode := child.Mode()
		childName := child.Name()
		childPath := filepath.Join(path, childName)

		if childMode.IsDir() {
			subChild, err := listFilesRecursive(childPath, childName)
			if err != nil {
				return dir{}, err //TODO if omit error is active, return slice, else stop function ¿maybe?
			}
			d.Dirs = append(d.Dirs, subChild)
		} else if childMode.IsRegular() {
			subChild, err := getFile(childPath)
			if err != nil {
				return dir{}, err //TODO if omit error is active, return slice, else stop function ¿maybe?
			}
			d.Files = append(d.Files, subChild)
		} else {
			// TODO symlinks and other cases
		}
	}

	return d, nil
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

	f.realPath = path
	f.Name = filepath.Base(path)
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
	if err != nil {
		return
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
