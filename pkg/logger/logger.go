package logger

import "github.com/Miguel-Dorta/logolang"

var (
	// Log is the global object for logging operation
	Log = logolang.NewLogger()

	// OmitErrors is the global value for knowing when the user prefers to omit errors
	OmitErrors = false
)
