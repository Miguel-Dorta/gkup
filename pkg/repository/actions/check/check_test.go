package check_test

import (
	"bytes"
	"encoding/json"
	"github.com/Miguel-Dorta/gkup/pkg/repository/actions/check"
	"regexp"
	"strconv"
	"testing"
)

type statusJSON struct {
	Type      string `json:"type"`
	Processed int    `json:"processed"`
	Total     int    `json:"total"`
}

type errorJSON struct {
	Type string `json:"type"`
	Err  string `json:"error"`
}

var errs = map[string]bool{
	"get data from filename testdata/files/96/random-file: error decoding hash from name: encoding/hex: invalid byte: U+0072 'r'": false, // Invalid name
	"hashes don't match in file testdata/files/64/64010f4b56691585c72c27067a86cd4e447c9c73fe8218e75d92d503ea05d6ad-5120": false, // Modified content
	"hashes don't match in file testdata/files/a7/a78fca73d3cf76c693af93c01d4ca3a2b810fd05e20aafa045a235d610069433-5120": false, // Modified filename (hash)
	"sizes don't match in file testdata/files/98/98d9c7552d34d9a179b919836555a408db32451fc9d50eaa86a38428dfd7bd9c-5120": false, // Modified content
	"sizes don't match in file testdata/files/cc/ccfd79d83959bb758db3a771100ba9e15e249b2e7348f3ae574fc484bc342241-5121": false, // Modified filename (size)
}

func TestCheck(t *testing.T) {
	var statusWriter, errorWriter = bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	if err := check.Check("testdata", 128*1024, false, statusWriter, errorWriter); err != nil {
		t.Fatalf("error checking files with JSON==false: %s", err)
	}
	checkStatusTXT(statusWriter.Bytes(), t)
	checkErrorTXT(errorWriter.Bytes(), t)

	statusWriter.Reset()
	errorWriter.Reset()
	for k := range errs {
		errs[k] = false
	}

	if err := check.Check("testdata", 128*1024, true, statusWriter, errorWriter); err != nil {
		t.Fatalf("error checking files with JSON==true: %s", err)
	}
	checkStatusJSON(statusWriter.Bytes(), t)
	checkErrorJSON(errorWriter.Bytes(), t)
}

func checkStatusJSON(status []byte, t *testing.T) {
	parts := bytes.Split(status, []byte{0})
	for _, part := range parts {
		if len(part) == 0 || string(part) == "\n" {
			continue
		}

		var stat statusJSON
		if err := json.Unmarshal(part, &stat); err != nil {
			t.Errorf("cannot unmarshal status msg \"%s\": %s", string(part), err)
			continue
		}
		if stat.Type != "status" {
			t.Errorf("incorrect type in status: %s", stat.Type)
			continue
		}
		if stat.Processed < 0 || stat.Processed > stat.Total {
			t.Errorf("incorrect processed/total in status json: %d - %d", stat.Processed, stat.Total)
		}
	}
}

func checkStatusTXT(status []byte, t *testing.T) {
	statusTxtRegex := regexp.MustCompile("^Processed files: (\\d+) of (\\d+)[\\n]?$")
	parts := bytes.Split(status, []byte{'\r'})
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if !statusTxtRegex.Match(part) {
			t.Errorf("status txt doesn't match regex: %s", string(part))
			continue
		}

		subParts := statusTxtRegex.FindSubmatch(part)[1:]
		processed, err := strconv.Atoi(string(subParts[0]))
		if err != nil {
			t.Errorf("cannot parse \"processed\" in status txt \"%s\": %s", string(part), err)
			continue
		}
		total, err := strconv.Atoi(string(subParts[1]))
		if err != nil {
			t.Errorf("cannot parse \"total\" in status txt \"%s\": %s", string(part), err)
			continue
		}
		if processed < 0 || processed > total {
			t.Errorf("incorrect processed/total in status txt: %d - %d", processed, total)
		}
	}
}

func checkErrorJSON(errors []byte, t *testing.T) {
	parts := bytes.Split(errors, []byte{0})
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		var e errorJSON
		if err := json.Unmarshal(part, &e); err != nil {
			t.Errorf("cannot unmarshal error msg \"%s\": %s", string(part), err)
			continue
		}
		if e.Type != "error" {
			t.Errorf("incorrect type in error: %s", e.Type)
			continue
		}

		if _, exists := errs[e.Err]; !exists {
			t.Errorf("unexpected errorJSON: %s", e.Err)
			continue
		}
		errs[e.Err] = true
	}

	for k, v := range errs {
		if !v {
			t.Errorf("not catched error: %s", k)
		}
	}
}

func checkErrorTXT(errors []byte, t *testing.T) {
	parts := bytes.Split(errors, []byte{'\n'})
	for i := range parts {
		part := string(parts[i])
		if len(part) == 0 || part[0] == '\n' {
			continue
		}
		part = part[1:]

		if _, exists := errs[part]; !exists {
			t.Errorf("unexpected errorTXT: %s", part)
			continue
		}
		errs[part] = true
	}

	for k, v := range errs {
		if !v {
			t.Errorf("not catched error: %s", k)
		}
	}
}
