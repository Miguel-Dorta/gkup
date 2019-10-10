package list

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// getTXT returns a easily-readable representation of the snapshot list provided.
func getTXT(snapList []*Snapshots) []byte {
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
func getJSON(snapList []*Snapshots) []byte {
	list := ListJSON{List: make([]Snapshots, len(snapList))}
	for i := range snapList {
		list.List[i] = Snapshots{
			Name:  snapList[i].Name,
			Times: snapList[i].Times,
		}
	}
	data, _ := json.Marshal(list)
	return append(data, '\n')
}
