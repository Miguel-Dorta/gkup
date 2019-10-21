package check

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/threadSafe"
	"os"
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

func printStatus(list *threadSafe.StringList, json bool, quit <-chan bool) {
	var printStatusFunc func(processed, total int)
	if json {
		printStatusFunc = printStatusJSON
	} else {
		printStatusFunc = printStatusTXT
	}
	seconds := time.NewTicker(time.Second).C

	for {
		select {
		case <-quit:
			return
		case <-seconds:
			printStatusFunc(list.GetPosUnsafe(), list.GetLenUnsafe())
		}
	}
}

func printStatusTXT(processed, total int) {
	fmt.Printf("\rProcessed files: %d of %d", processed, total)
}

func printStatusJSON(processed, total int) {
	data, _ := json.Marshal(statusJSON{
		Type:      "status",
		Processed: processed,
		Total:     total,
	})
	_, _ = os.Stdout.Write(append(data, '\n', 0))
}

func printError(err error, json bool) {
	if json {
		data, _ := json.Marshal(errorJSON{
			Type: "error",
			Err:  err.Error(),
		})
		_, _ = os.Stdout.Write(append(data, '\n', 0))
	} else {
		_, _ = os.Stderr.WriteString("\r" + err.Error() + "\n")
	}
}
