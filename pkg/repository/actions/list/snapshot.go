package list

import (
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type ListJSON struct {
	List []Snapshots `json:"snapshots"`
}

type Snapshots struct {
	Name string `json:"name"`
	Times []int64 `json:"times"`
}

// snapshotNameRegex represents the name that the snapshots file should follow.
var snapshotNameRegex = regexp.MustCompile("^(\\d{4})-(\\d{2})-(\\d{2})_(\\d{2})-(\\d{2})-(\\d{2}).json$")

// getSnapshots list a path and return an snapshot type with the name provided, and a slice of
// the times of the snapshots found in that path.
func getSnapshots(path, name string) (*Snapshots, error) {
	fileList, err := utils.ListDir(path)
	if err != nil {
		return nil, &os.PathError{
			Op:   "list snapshots folder",
			Path: path,
			Err:  err,
		}
	}

	snap := Snapshots{
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
