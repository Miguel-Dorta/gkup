package list

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup/api"
	"github.com/Miguel-Dorta/gkup/pkg/repository"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// snapshotNameRegex represents the name that the snapshots file should follow.
var snapshotNameRegex = regexp.MustCompile("^(\\d{4})-(\\d{2})-(\\d{2})_(\\d{2})-(\\d{2})-(\\d{2}).json$")

// List takes the repo path, list all the snapshots of that repo, and writes them in the writer
// provided in an human-readable way or in JSON depending of the bool provided.
func List(path string, inJson bool, writeTo io.Writer) error {
	snapList := make([]*api.Snapshots, 0, 100)
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

// getTXT returns a easily-readable representation of the snapshot list provided.
func getTXT(snapList []*api.Snapshots) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 100))

	for _, snap := range snapList {
		name := snap.Name
		if name == "" {
			name = "[no-name]"
		}
		_, _ = buf.WriteString(name)
		_ = buf.WriteByte('\n')

		for _, unixTime := range snap.Times {
			t := time.Unix(unixTime, 0).UTC()
			Y, M, D := t.Date()
			h, m, s := t.Clock()
			_, _ = fmt.Fprintf(buf, "- %04d/%02d/%02d %02d:%02d:%02d\n", Y, M, D, h, m, s)
		}
		_ = buf.WriteByte('\n')
	}
	return buf.Bytes()
}

// getJSON returns the JSON representation of the snapshot list provided.
func getJSON(snapList []*api.Snapshots) []byte {
	list := api.List{List: make([]api.Snapshots, len(snapList))}
	for i := range snapList {
		list.List[i] = api.Snapshots{
			Name:  snapList[i].Name,
			Times: snapList[i].Times,
		}
	}
	data, _ := json.Marshal(list)
	return append(data, '\n')
}

// getSnapshots list a path and return an snapshot type with the name provided, and a slice of
// the times of the snapshots found in that path.
func getSnapshots(path, name string) (*api.Snapshots, error) {
	fileList, err := utils.ListDir(path)
	if err != nil {
		return nil, &os.PathError{
			Op:   "list snapshots folder",
			Path: path,
			Err:  err,
		}
	}

	snap := api.Snapshots{
		Name:  name,
		Times: make([]int64, 0, len(fileList)),
	}
	for _, f := range fileList {
		if isSnapshot(f) {
			snap.Times = append(snap.Times, getDateOfSnapshot(f.Name()))
		}
	}

	sort.Slice(snap.Times, func(i, j int) bool {
		return snap.Times[i] < snap.Times[j]
	})

	return &snap, nil
}

// getDateOfSnapshot returns an Unix timestamp of the date contained in the name of a snapshot file.
// The name must have been checked with isSnapshot, otherwise it can panic.
func getDateOfSnapshot(name string) int64 {
	panicMsg := "parse error: not checked snapshot: "
	parts := snapshotNameRegex.FindStringSubmatch(name)
	if len(parts) != 7 {
		panic(panicMsg + "unexpected number of parts")
	}

	var dates [6]int
	for i := range dates {
		x, err := strconv.Atoi(parts[i+1])
		if err != nil {
			panic(panicMsg + err.Error())
		}
		dates[i] = x
	}

	return time.Date(dates[0], time.Month(dates[1]), dates[2], dates[3], dates[4], dates[5], 0, time.UTC).Unix()
}

// isSnapshots returns true if the FileInfo provided is a snapshot file
func isSnapshot(fi os.FileInfo) bool {
	return fi.Mode().IsRegular() && snapshotNameRegex.MatchString(fi.Name())
}
