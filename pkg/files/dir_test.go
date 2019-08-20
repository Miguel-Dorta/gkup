package files_test

import (
	"encoding/json"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"sort"
	"strings"
	"testing"
)

func TestNewDir(t *testing.T) {
	d, fileList, err := files.NewDir("../../test")
	if err != nil {
		t.Fatalf("error listing working directory: %s", err)
	}

	// Sort
	d = sortDir(d)
	sort.Slice(fileList, func(i, j int) bool {
		return strings.ToLower(fileList[i].RealPath) < strings.ToLower(fileList[j].RealPath)
	})

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("error converting dir to JSON: %s", err)
	}

	// This string is checked by hand and it's a pain in the ass to check it
	expectedResult := "{\"name\":\"test\",\"dirs\":[{\"name\":\"aYVxBryiOoTL\",\"dirs\":[],\"files\":[{\"name\":\"d7MFNsZ72ZDi\",\"size\":1024,\"hash\":null},{\"name\":\"EwG9iXlHBMzw\",\"size\":1024,\"hash\":null},{\"name\":\"sRXDiPdbkQfx\",\"size\":1024,\"hash\":null},{\"name\":\"XnkUi8zHyYQS\",\"size\":1024,\"hash\":null}]},{\"name\":\"gbjQzr86G0iI\",\"dirs\":[{\"name\":\"nf9jUiAOYBMl\",\"dirs\":[{\"name\":\"sd50GOkQxkwz\",\"dirs\":[],\"files\":[{\"name\":\"AgCOOVeuHi5M\",\"size\":1024,\"hash\":null},{\"name\":\"c9Eoceyr4Ez5\",\"size\":1024,\"hash\":null},{\"name\":\"hE90zLtfdJgL\",\"size\":1024,\"hash\":null},{\"name\":\"mHMyKkS7R0Sc\",\"size\":1024,\"hash\":null}]}],\"files\":[{\"name\":\"3hCnntC9lcoP\",\"size\":1024,\"hash\":null},{\"name\":\"Wkw6RQcKun1I\",\"size\":1024,\"hash\":null}]}],\"files\":[{\"name\":\"4jCqKSh1Qbuo\",\"size\":1024,\"hash\":null},{\"name\":\"dEzW3eCn1Jbp\",\"size\":1024,\"hash\":null},{\"name\":\"PcZuqYiW7CWW\",\"size\":1024,\"hash\":null},{\"name\":\"Sj3O4c8klFwl\",\"size\":1024,\"hash\":null},{\"name\":\"Tck0MYvYNSmN\",\"size\":1024,\"hash\":null}]},{\"name\":\"Ls9xvOjzEG7f\",\"dirs\":[{\"name\":\"EV6n9EEI901l\",\"dirs\":[],\"files\":[{\"name\":\"LvPz1HJQMECi\",\"size\":1024,\"hash\":null},{\"name\":\"P3Smv4U19Mnb\",\"size\":1024,\"hash\":null},{\"name\":\"rnVMwZGoiyG7\",\"size\":1024,\"hash\":null},{\"name\":\"ymDAygI4nxVj\",\"size\":1024,\"hash\":null}]}],\"files\":[{\"name\":\"d6pl9ZIEN4KF\",\"size\":1024,\"hash\":null},{\"name\":\"NJFxnlLvxZu6\",\"size\":1024,\"hash\":null},{\"name\":\"pDf1D7y9mN5Z\",\"size\":1024,\"hash\":null},{\"name\":\"pVjdvDtk2Elg\",\"size\":1024,\"hash\":null},{\"name\":\"uRQBk47cOcT6\",\"size\":1024,\"hash\":null}]}],\"files\":[{\"name\":\"2d3f8vsTFvQB\",\"size\":1024,\"hash\":null},{\"name\":\"2XOdc2DsuIrz\",\"size\":1024,\"hash\":null},{\"name\":\"794ixIkN6CLN\",\"size\":1024,\"hash\":null},{\"name\":\"eGbEQTqQytsl\",\"size\":1024,\"hash\":null},{\"name\":\"MnZCGblTujPr\",\"size\":1024,\"hash\":null},{\"name\":\"WKMM0YnkArCi\",\"size\":1024,\"hash\":null}]}"
	if string(data) != expectedResult {
		t.Fatalf("results does not match\n-> Expected: %s\n-> Found: %s", expectedResult, string(data))
	}
}

func sortDir(d files.Dir) files.Dir {
	sort.Slice(d.Dirs, func(i, j int) bool {
		return strings.ToLower(d.Dirs[i].Name) < strings.ToLower(d.Dirs[j].Name)
	})

	sort.Slice(d.Files, func(i, j int) bool {
		return strings.ToLower(d.Files[i].Name) < strings.ToLower(d.Files[j].Name)
	})

	for i := range d.Dirs {
		d.Dirs[i] = sortDir(d.Dirs[i])
	}

	return d
}
