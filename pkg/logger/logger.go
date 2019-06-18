package logger

import "github.com/Miguel-Dorta/logolang"

var Log = logolang.NewLogger(nil, nil, nil, nil)

func init() {
	_ = Log.SetLevel(3)
}
