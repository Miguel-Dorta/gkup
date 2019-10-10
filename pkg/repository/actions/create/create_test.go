package create_test

import (
	"github.com/Miguel-Dorta/gkup/internal"
	"github.com/Miguel-Dorta/gkup/pkg/repository/actions/create"
	"github.com/Miguel-Dorta/gkup/pkg/repository/settings"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

var testingPath = filepath.Join(os.TempDir(), "gkup_pkg_repository_TestCreate")

func init() {
	internal.Version = "v1.0.0"
}

func createCases(t *testing.T) {
	// Empty dir
	if err := os.MkdirAll(filepath.Join(testingPath, "empty_dir"), 0777); err != nil {
		t.Fatal("cannot create case empty_dir case")
	}

	// Non-empty dir
	if err := os.MkdirAll(filepath.Join(testingPath, "non_empty_dir"), 0777); err != nil {
		t.Fatal("cannot create case non_empty_dir case dir")
	}
	if err := createEmptyFile(filepath.Join(testingPath, "non_empty_dir", "file")); err != nil {
		t.Fatal("cannot create case non_empty_dir case file")
	}

	// File
	if err := createEmptyFile(filepath.Join(testingPath, "file")); err != nil {
		t.Fatal("cannot create case file")
	}

	//TODO symlink case
}

func checkGoodCase(caseStr string, t *testing.T) {
	path := filepath.Join(testingPath, caseStr)
	if err := create.Create(path, "sha256"); err != nil {
		t.Errorf("error creating case %s: %s", caseStr, err)
		return
	}
	if !isWellFormed(path) {
		t.Errorf("case %s is not well formed", caseStr)
	}
}

func checkBadCase(caseStr string, t *testing.T) {
	path := filepath.Join(testingPath, caseStr)
	if err := create.Create(path, "sha256"); err == nil {
		t.Errorf("not error detected in case %s", caseStr)
	}
}

func TestCreate(t *testing.T) {
	defer os.RemoveAll(testingPath)
	createCases(t)

	checkGoodCase("non_existing", t)
	checkGoodCase("empty_dir", t)

	checkBadCase("non_empty_dir", t)
	checkBadCase("file", t)
	checkBadCase(filepath.Join("file", "fileChild"), t)
}

//////////////////////////
//// HELPER FUNCTIONS ////
//////////////////////////

func createEmptyFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

func listDir(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	sort.Strings(names)

	return names, nil
}

func equalsStringsSlice(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func isWellFormed(path string) bool {
	// Check top level
	topLevel, err := listDir(path)
	if err != nil {
		return false
	}
	if !equalsStringsSlice(topLevel, []string{"files", "settings.toml", "snapshots"}) {
		return false
	}

	// Check backup dir
	filesSubDirs, err := listDir(filepath.Join(path, "files"))
	if err != nil {
		return false
	}
	if !equalsStringsSlice(filesSubDirs, []string{
		"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "0a", "0b", "0c", "0d", "0e", "0f",
		"10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1a", "1b", "1c", "1d", "1e", "1f",
		"20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "2a", "2b", "2c", "2d", "2e", "2f",
		"30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "3a", "3b", "3c", "3d", "3e", "3f",
		"40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "4a", "4b", "4c", "4d", "4e", "4f",
		"50", "51", "52", "53", "54", "55", "56", "57", "58", "59", "5a", "5b", "5c", "5d", "5e", "5f",
		"60", "61", "62", "63", "64", "65", "66", "67", "68", "69", "6a", "6b", "6c", "6d", "6e", "6f",
		"70", "71", "72", "73", "74", "75", "76", "77", "78", "79", "7a", "7b", "7c", "7d", "7e", "7f",
		"80", "81", "82", "83", "84", "85", "86", "87", "88", "89", "8a", "8b", "8c", "8d", "8e", "8f",
		"90", "91", "92", "93", "94", "95", "96", "97", "98", "99", "9a", "9b", "9c", "9d", "9e", "9f",
		"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9", "aa", "ab", "ac", "ad", "ae", "af",
		"b0", "b1", "b2", "b3", "b4", "b5", "b6", "b7", "b8", "b9", "ba", "bb", "bc", "bd", "be", "bf",
		"c0", "c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8", "c9", "ca", "cb", "cc", "cd", "ce", "cf",
		"d0", "d1", "d2", "d3", "d4", "d5", "d6", "d7", "d8", "d9", "da", "db", "dc", "dd", "de", "df",
		"e0", "e1", "e2", "e3", "e4", "e5", "e6", "e7", "e8", "e9", "ea", "eb", "ec", "ed", "ee", "ef",
		"f0", "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "fa", "fb", "fc", "fd", "fe", "ff",
	}) {
		return false
	}

	// Check settings content
	sett, err := settings.Read(filepath.Join(path, "settings.toml"))
	if err != nil {
		return false
	}
	return sett.HashAlgorithm == "sha256" && sett.Version == "v1.0.0"
}
