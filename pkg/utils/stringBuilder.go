package utils

import "strings"

// stringBuilderAppend creates syntactic sugar for appending strings to a strings.Builder
func stringBuilderAppend(b *strings.Builder, s ...string) {
	for _, str := range s {
		b.WriteString(str)
	}
}
