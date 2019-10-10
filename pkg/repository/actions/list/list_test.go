package list_test

import (
	"github.com/Miguel-Dorta/gkup/pkg/repository/actions/list"
	"strings"
	"testing"
)

var (
	expectedTXT = `[no-name]
- 1998/12/28 01:23:45
- 2018/10/08 16:43:32

custom_name
- 0001/01/01 00:00:00
- 2099/12/31 23:59:59

MyPC
- 1969/12/31 23:59:59
- 1970/01/01 00:00:00

mypc
- 2038/01/19 03:14:07
- 2038/01/19 03:14:08

`
	expectedJSON = `{"snapshots":[{"name":"","times":[914808225,1539017012]},{"name":"custom_name","times":[-62135596800,4102444799]},{"name":"MyPC","times":[-1,0]},{"name":"mypc","times":[2147483647,2147483648]}]}
`
)

func TestList(t *testing.T) {
	var actualTXT, actualJSON = &strings.Builder{}, &strings.Builder{}
	testdataPath := "testdata"

	// Test TXT export
	err := list.List(testdataPath, false, actualTXT)
	if err != nil {
		t.Errorf("error found listing TXT: %s", err)
	} else if expectedTXT != actualTXT.String() {
		t.Errorf("TXT doesn't match the expected result\n-> Expected: %s\n-> Found: %s", expectedTXT, actualTXT.String())
	}

	// Test JSON export
	err = list.List(testdataPath, true, actualJSON)
	if err != nil {
		t.Errorf("error found listing JSON: %s", err)
	} else if expectedJSON != actualJSON.String() {
		t.Errorf("JSON doesn't match the expected result\n-> Expected: %s\n-> Found: %s", expectedJSON, actualJSON.String())
	}
}
