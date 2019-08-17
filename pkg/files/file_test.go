package files_test

import (
	"bytes"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"testing"
)

func TestGetFileFromName(t *testing.T) {
	validFileNames := []struct{
		filename string
		expectedResult files.File
	}{
		{
			"6878758d2aed4646ee85cb88e7c60817edb51fcaec25f1a9c4939deadb93fba6-165869644",
			files.File{
				Name:"6878758d2aed4646ee85cb88e7c60817edb51fcaec25f1a9c4939deadb93fba6-165869644",
				Size: 165869644,
				Hash: []byte{104, 120, 117, 141, 42, 237, 70, 70, 238, 133, 203, 136, 231, 198, 8, 23, 237, 181, 31, 202, 236, 37, 241, 169, 196, 147, 157, 234, 219, 147, 251, 166},
			},
		},
		{
			"bdfb11f19e51f360e391b24a12274eb3-2135",
			files.File{
				Name: "bdfb11f19e51f360e391b24a12274eb3-2135",
				Size: 2135,
				Hash: []byte{189, 251, 17, 241, 158, 81, 243, 96, 227, 145, 178, 74, 18, 39, 78, 179},
			},
		},
		{
			"6a04e805e2080321ace61f0403e94b7ef75c06c4394a94a1f528d242387ef3f487baa2a9b891629b986cc95a9332bad681d8adfbb20352dd8a89da61a2fce974-531",
			files.File{
				Name: "6a04e805e2080321ace61f0403e94b7ef75c06c4394a94a1f528d242387ef3f487baa2a9b891629b986cc95a9332bad681d8adfbb20352dd8a89da61a2fce974-531",
				Size: 531,
				Hash: []byte{106, 4, 232, 5, 226, 8, 3, 33, 172, 230, 31, 4, 3, 233, 75, 126, 247, 92, 6, 196, 57, 74, 148, 161, 245, 40, 210, 66, 56, 126, 243, 244, 135, 186, 162, 169, 184, 145, 98, 155, 152, 108, 201, 90, 147, 50, 186, 214, 129, 216, 173, 251, 178, 3, 82, 221, 138, 137, 218, 97, 162, 252, 233, 116},
			},
		},
	}

	invalidFileNames := []string{
		"6878758d2aed4646ee85cb88e7c60817edb51fcaec25f1a9c4939deadb93fbax-165869644",
		"6878758d2aed4646ee85cb88e7c60817edb51fcaec25f1a9c4939deadb93fbax-1658696a4",
		"6878758d2aed4646ee85cb88e7c60817edb51fcaec25f1a9c493-deadb93fbax-165869644",
		"6878758d2aed4646ee85cb88e7c60817edb51fcaec25f19c4939deadb93fbax-165869-644",
		"Generic-filename.txt",
		"",
	}

	for _, valid := range validFileNames {
		actualResult, err := files.GetFileFromName(valid.filename)
		if err != nil {
			t.Errorf("error found in valid filename.\n-> Filename: %s\nError: %s", valid.filename, err)
			continue
		}

		if actualResult.Name != valid.expectedResult.Name {
			t.Errorf("expected name is not what was found.\n-> Expected: %s\n-> Found: %s", valid.expectedResult.Name, actualResult.Name)
		}

		if actualResult.Size != valid.expectedResult.Size {
			t.Errorf("expected size is not what was found.\n-> Expected: %d\n-> Found: %d", valid.expectedResult.Size, actualResult.Size)
		}

		if !bytes.Equal(actualResult.Hash, valid.expectedResult.Hash) {
			t.Errorf("expected hash is not what was found.\n-> Expected: %+v\n-> Found: %+v", valid.expectedResult.Hash, actualResult.Hash)
		}
	}

	for _, invalid := range invalidFileNames {
		_, err := files.GetFileFromName(invalid)
		if err == nil {
			t.Errorf("invalid filename \"%s\" was parsed into a file", invalid)
		}
	}
}
