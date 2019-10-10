package list

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/repository"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// List takes the repo path, list all the snapshots of that repo, and writes them in the writer
// provided in an human-readable way or in JSON depending of the bool provided.
func List(path string, inJson bool, writeTo io.Writer) error {
	snapList := make([]*Snapshots, 0, 100)
	snapshotsFolderPath := filepath.Join(path, repository.SnapshotsFolderName)

	// Add snapshots with no name defined
	noNameSnap, err := getSnapshots(snapshotsFolderPath, "")
	if err != nil {
		return fmt.Errorf("cannot get snapshots: %w", err)
	}
	snapList = append(snapList, noNameSnap)

	// Get file list
	fileList, err := utils.ListDir(snapshotsFolderPath)
	if err != nil {
		return &os.PathError{
			Op:   "list snapshots folder",
			Path: snapshotsFolderPath,
			Err:  err,
		}
	}
	// Iterate folders to get snapshots with name
	for _, f := range fileList {
		if !f.IsDir() {
			continue
		}

		// Append snapshots
		snap, err := getSnapshots(filepath.Join(snapshotsFolderPath, f.Name()), f.Name())
		if err != nil {
			return fmt.Errorf("cannot get snapshots: %w", err)
		}
		snapList = append(snapList, snap)
	}

	// Sort result
	sort.Slice(snapList, func(i, j int) bool {
		iLow := strings.ToLower(snapList[i].Name)
		jLow := strings.ToLower(snapList[j].Name)

		if iLow == jLow {
			return snapList[i].Name < snapList[j].Name
		}
		return iLow < jLow
	})

	// Get data formatted
	var output []byte
	if inJson {
		output = getJSON(snapList)
	} else {
		output = getTXT(snapList)
	}

	// Write output
	if _, err := writeTo.Write(output); err != nil {
		return fmt.Errorf("cannot write list to writer provided: %w", err)
	}
	return nil
}
