package pkg

import (
	"strings"
)

func stringBuilderAppend(b *strings.Builder, s ...string) {
	for _, str := range s {
		b.WriteString(str)
	}
}
