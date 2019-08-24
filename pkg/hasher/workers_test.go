package hasher

import (
	"encoding/hex"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"github.com/Miguel-Dorta/gkup/pkg/threadSafe"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// This tests are made to test the single behaviour of every process.
// This means that this test will not be concurrent.
// For concurrent testing go to multiHasher_test.go

var hashes = []struct {
	path, hash string
}{
	{
		path: "../../test/2d3f8vsTFvQB",
		hash: "60bfd58e6499978c2ff66727ea95cbd6d9f6b42e4ee808ca7c1a3e8ffa955a3f",
	},
	{
		path: "../../test/2XOdc2DsuIrz",
		hash: "066e69f54a1040037f07df53e7e5cf419ecaca2254432ca522d2b48421c588c8",
	},
	{
		path: "../../test/794ixIkN6CLN",
		hash: "1d002351a45f5d2dd01694701759110fa32f702b66c7ee71adaba10d48ec9e22",
	},
	{
		path: "../../test/aYVxBryiOoTL/d7MFNsZ72ZDi",
		hash: "7dfb75ac5d328a0680946ec85d7047b587ed59011e94ff9450ccdc4e706d09a2",
	},
	{
		path: "../../test/aYVxBryiOoTL/EwG9iXlHBMzw",
		hash: "f3c2f01397e94f06ed8ac586af26c1a2870a9cd44b8d867d6bc045cdf809843c",
	},
	{
		path: "../../test/aYVxBryiOoTL/sRXDiPdbkQfx",
		hash: "918251ffa8cadce4c91520cd606d83527d3bec20c502edb8904904d911fb3127",
	},
	{
		path: "../../test/aYVxBryiOoTL/XnkUi8zHyYQS",
		hash: "d6d51f16f0da298a824bb0bf3aef78e1cc7d108604f733c14bee5693b564f06e",
	},
	{
		path: "../../test/eGbEQTqQytsl",
		hash: "86d45a7de3519d50501f0822d05c042ca9c0df70a6e6735384054e51967d165b",
	},
	{
		path: "../../test/gbjQzr86G0iI/4jCqKSh1Qbuo",
		hash: "7cb2ff54bb1e27f433bd83f153c577d440896dc26bb2d07a1a2bac775872507d",
	},
	{
		path: "../../test/gbjQzr86G0iI/dEzW3eCn1Jbp",
		hash: "ba271bae5c59f1e547642a8e9069cc2bb6206d4fb756fdfd96226fa3464c1ac8",
	},
	{
		path: "../../test/gbjQzr86G0iI/nf9jUiAOYBMl/3hCnntC9lcoP",
		hash: "97a2a0477e4129d41a5b09f12bd1e0075f118458202bbe6434296aeddb03c78c",
	},
	{
		path: "../../test/gbjQzr86G0iI/nf9jUiAOYBMl/sd50GOkQxkwz/AgCOOVeuHi5M",
		hash: "1243ceb13341c88bcc96cc33c799c8c2e8236b136cea71c9331569ac792ee8a0",
	},
	{
		path: "../../test/gbjQzr86G0iI/nf9jUiAOYBMl/sd50GOkQxkwz/c9Eoceyr4Ez5",
		hash: "929185768ea22f2572949905e75d0bb2dfe51299234cedd42dc24c83c9a417f1",
	},
	{
		path: "../../test/gbjQzr86G0iI/nf9jUiAOYBMl/sd50GOkQxkwz/hE90zLtfdJgL",
		hash: "891f5dcd42b100bcea5c38522e0fa35eb1b9e23665b98665183680dcb1bf6ae0",
	},
	{
		path: "../../test/gbjQzr86G0iI/nf9jUiAOYBMl/sd50GOkQxkwz/mHMyKkS7R0Sc",
		hash: "12b958a11b7db4144188aa4bed476164b5e8a82a844927b2c855bea4ac7099f9",
	},
	{
		path: "../../test/gbjQzr86G0iI/nf9jUiAOYBMl/Wkw6RQcKun1I",
		hash: "f51c5b8325a8c5b6884783b823eae2cb4f52bf959213bd68b955f1770cbc68fd",
	},
	{
		path: "../../test/gbjQzr86G0iI/PcZuqYiW7CWW",
		hash: "7232b55e2dd9891f2fbd933b2e586629b33ace81172ea7d32ecc9bb8b9b294cf",
	},
	{
		path: "../../test/gbjQzr86G0iI/Sj3O4c8klFwl",
		hash: "b2396f640c596e8f15df93df7e3d8a09e85559b92bf3929defd8294761c1df53",
	},
	{
		path: "../../test/gbjQzr86G0iI/Tck0MYvYNSmN",
		hash: "441c491d8dc11bcd2db9e04bbf8b683f2714a2bb8d9ec15b63571c60ca3bd7a0",
	},
	{
		path: "../../test/Ls9xvOjzEG7f/d6pl9ZIEN4KF",
		hash: "140a779f1494f51e0e30c9d1a7023a3a65974f6bc9685e2c81e9b9f5bef829af",
	},
	{
		path: "../../test/Ls9xvOjzEG7f/EV6n9EEI901l/LvPz1HJQMECi",
		hash: "7d043a25178152af71c995356cb07982dc356af8e63ad0d1441a8b4fc030dd3d",
	},
	{
		path: "../../test/Ls9xvOjzEG7f/EV6n9EEI901l/P3Smv4U19Mnb",
		hash: "a9b108723186be4c2a9eb8191c9010bb80be8d3c291643c118aa42875b6961b0",
	},
	{
		path: "../../test/Ls9xvOjzEG7f/EV6n9EEI901l/rnVMwZGoiyG7",
		hash: "13388f0abdacdffc838cd34c5840461cf2295863cced39f8b9109bb562bdc075",
	},
	{
		path: "../../test/Ls9xvOjzEG7f/EV6n9EEI901l/ymDAygI4nxVj",
		hash: "aea4dd685007e46e551ade220f537a20466d5625fd15658ffb04bc2be4d6c376",
	},
	{
		path: "../../test/Ls9xvOjzEG7f/NJFxnlLvxZu6",
		hash: "762d988d94288e5d8342c41f0ecb549a9a99a20ea2d44b56005656bf49f6c286",
	},
	{
		path: "../../test/Ls9xvOjzEG7f/pDf1D7y9mN5Z",
		hash: "725cb5acdcb1b8fea546f051dcb2d3e072f6daffba3106315888eb994710108e",
	},
	{
		path: "../../test/Ls9xvOjzEG7f/pVjdvDtk2Elg",
		hash: "797cd11b3401a2aeb1b2666390c8393c103dcab25fe70ffc200680cb8c8570b7",
	},
	{
		path: "../../test/Ls9xvOjzEG7f/uRQBk47cOcT6",
		hash: "800b1a394069f5582dfd45c62b76271aa3ab0a2b37a1fc5c586babb1d8ecb8e0",
	},
	{
		path: "../../test/MnZCGblTujPr",
		hash: "f1b881586c4a3a667da3866a85460bae935bc694d7d938c20913e7ead6f0a8ab",
	},
	{
		path: "../../test/WKMM0YnkArCi",
		hash: "6fb792e1e7def0eca4051975a6149f5b9ea5e0e3fa055c95c670e89640aae302",
	},
}

func TestFileChecker(t *testing.T) {
	// Create tmp dir for working in it
	tmpDir := "/tmp/gkup-TestFileChecker"
	err := os.MkdirAll(tmpDir, 0777)
	if err != nil {
		t.Fatalf("Cannot create folder in /tmp: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	// Copy all test files in tmp dir with specific name
	fileList := make([]string, 0, len(hashes))
	for _, hash := range hashes {
		data, err := ioutil.ReadFile(hash.path)
		if err != nil {
			t.Errorf("cannot read file \"%s\": %s", hash.path, err)
			continue
		}
		fileName := filepath.Join(tmpDir, hash.hash + "-1024")
		if err = ioutil.WriteFile(fileName, data, 0777); err != nil {
			t.Errorf("cannot write file in %s: %s", tmpDir, err)
			continue
		}
		fileList = append(fileList, fileName)
	}

	// Create all variables necessary for the test
	h, err := New("sha256")
	if err != nil {
		t.Fatal("sha256 was not a valid hash for creating hasher")
	}
	var fuse fuse
	var wg sync.WaitGroup
	wg.Add(1)

	// Do the actual test
	h.fileChecker(threadSafe.NewStringList(fileList), &fuse, &wg)

	// Check if there are false positives
	if fuse {
		t.Fatal("False errors found")
	}

	// Test it now for real errors
	fileList = append(fileList, filepath.Join(tmpDir, "random-invalid.name"))
	wg.Add(1) // We can reuse the same not-triggered fuse and not-waited wg
	h.fileChecker(threadSafe.NewStringList(fileList), &fuse, &wg)
	if !fuse {
		t.Fatal("Fuse was not triggered with actual errors")
	}
}

// TODO separate the validation algorithm of CheckFiles for testing it

func TestFileGetter(t *testing.T) {
	// Create hasher
	h, err := New("sha256")
	if err != nil {
		t.Fatal("sha256 was not a valid hash for creating hasher")
	}

	// Create list of paths and list where results will be saved
	pathList := make([]string, len(hashes))
	fileList := threadSafe.NewFileList(make([]*files.File, 0, len(hashes)))
	for i := range hashes {
		pathList[i] = hashes[i].path
	}

	// Get files (actual test)
	if err = h.fileGetter(threadSafe.NewStringList(pathList), fileList); err != nil {
		t.Fatalf("error getting files: %s", err)
	}

	// Check if the files were got correctly
	for i, f := range fileList.GetList() {
		if f == nil {
			t.Logf("skipping i == %d", i)
			continue
		}
		if f.Name != filepath.Base(hashes[i].path) {
			t.Errorf(
				"filenames don't match in file from path \"%s\"\n-> Expected: %s\n-> Found: %s",
				hashes[i].path, filepath.Base(hashes[i].path), f.Name,
			)
		}
		if f.RealPath != hashes[i].path {
			t.Errorf("paths don't match\n-> Expected: %s\n-> Found: %s", hashes[i].path, f.RealPath)
		}
		if f.Size != 1024 {
			t.Errorf("sizes don't match\n-> Expected: 1024\n-> Found: %d", f.Size)
		}
		if hex.EncodeToString(f.Hash) != hashes[i].hash {
			t.Errorf("hashes don't match\n-> Expected: %s\n-> Found: %s", hashes[i].hash, hex.EncodeToString(f.Hash))
		}
	}
}

func TestFileHasher(t *testing.T) {
	// Create file list
	fileList := make([]*files.File, 0, len(hashes))
	for _, hash := range hashes {
		f, err := files.NewFile(hash.path)
		if err != nil {
			t.Fatalf("error creating file in path \"%s\": %s", hash.path, err)
		}
		fileList = append(fileList, f)
	}
	safeFileList := threadSafe.NewFileList(fileList)

	// Create hasher
	h, err := New("sha256")
	if err != nil {
		t.Fatal("sha256 was not a valid hash for creating hasher")
	}

	// Do the actual test
	if err = h.fileHasher(safeFileList); err != nil {
		t.Fatalf("error hashing files: %s", err)
	}

	// Check if they were hashed correctly
	for i, f := range safeFileList.GetList() {
		if hex.EncodeToString(f.Hash) != hashes[i].hash {
			t.Errorf("hashes don't match in file \"%s\"\n-> Expected: %s\n-> Found: %s",
				hashes[i].path, hashes[i].hash, hex.EncodeToString(f.Hash),
			)
		}
	}
}
