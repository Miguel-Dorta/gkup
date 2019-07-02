package utils

// isHidden returns whether the name provided is hidden
func isHidden(path, name string) (bool, error) {
	return name[0] == '.', nil
}
