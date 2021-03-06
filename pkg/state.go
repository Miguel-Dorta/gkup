package pkg

import (
	"github.com/Miguel-Dorta/logolang"
	"runtime"
)

var (
	BufferSize      = 4 * 1024 * 1024
	Log             = logolang.NewLogger()
	NumberOfThreads = runtime.NumCPU()
	OmitHidden      = false
	OmitErrors      = false
	Version         string
)
