package check

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/threadSafe"
	"time"
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

func printStatusAsync(list *threadSafe.StringList, quit <-chan bool) {
	seconds := time.NewTicker(time.Second).C
	for {
		select {
		case <-quit:
			_, _ = statusWriter.Write([]byte{'\n'})
			_, _ = errorWriter.Write([]byte{'\n'})
			return
		case <-seconds:
			printStatus(list.GetPosUnsafe(), list.GetLenUnsafe())
		}
	}
}

func printStatus(processed, total int) {
	var b []byte
	if jsonOutput {
		b = getStatusJSON(processed, total)
	} else {
		b = getStatusTXT(processed, total)
	}
	_, _ = statusWriter.Write(b)
}

func getStatusTXT(processed, total int) []byte {
	return []byte(fmt.Sprintf("\rProcessed files: %d of %d", processed, total))
}

func getStatusJSON(processed, total int) []byte {
	data, _ := json.Marshal(statusJSON{
		Type:      "status",
		Processed: processed,
		Total:     total,
	})
	return append(data, '\n', 0)
}

func printError(err error) {
	var b []byte
	if jsonOutput {
		b = getErrorJSON(err)
	} else {
		b = getErrorTXT(err)
	}
	_, _ = errorWriter.Write(b)
}

func getErrorTXT(err error) []byte {
	return []byte("\r" + err.Error() + "\n")
}

func getErrorJSON(err error) []byte {
	data, _ := json.Marshal(errorJSON{
		Type: "error",
		Err:  err.Error(),
	})
	return append(data, '\n', 0)
}
