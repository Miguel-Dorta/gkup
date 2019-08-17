package utils

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg"
	"io"
	"os"
)

// CopyFile copies a file from origin path to destiny path
func CopyFile(origin, destiny string, buffer []byte) error {
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

	pkg.Log.Debugf("Copying file %s to %s", origin, destiny)
	if _, err = io.CopyBuffer(destinyFile, originFile, buffer); err != nil {
		errStr := fmt.Sprintf("Error copying file from %s to %s: %s", origin, destiny, err.Error())
		pkg.Log.Error(errStr)

		if err = destinyFile.Close(); err == nil {
			pkg.Log.Debugf("File %s closed", destiny)
			if err = os.Remove(destiny); err == nil {
				pkg.Log.Debugf("File %s removed", destiny)
				return errors.New(errStr)
			}
		}

		errStr = fmt.Sprintf("%s\n-> There's a corrupt file in \"%s\". Please, remove it", errStr, destiny)
		pkg.Log.Error(errStr)
		return errors.New(errStr)
	}

	if err = destinyFile.Close(); err != nil {
		return fmt.Errorf("error closing file in \"%s\", please check it", err)
	}
	return nil
}

// ListDir lists the directory from the path provided
func ListDir(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdir(-1)
}

// IsHidden returns whether the file provided is hidden
func IsHidden(name string) bool {
	return isHidden(name)
}
