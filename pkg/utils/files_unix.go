package utils

// isHidden returns whether the name provided is hidden
func isHidden(name string) bool {
	return name[0] == '.'
}
