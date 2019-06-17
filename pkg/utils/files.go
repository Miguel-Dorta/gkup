package utils

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/tmp"
	"io"
	"os"
	"strings"
)

// CopyFile copies a file from origin path to destiny path
func CopyFile(origin, destiny string) error {
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

	if _, err = io.CopyBuffer(destinyFile, originFile, tmp.CopyBuf); err != nil {
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

// ListDir lists the directory from the path provided
func ListDir(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdir(-1)
}
