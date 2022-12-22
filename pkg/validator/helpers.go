package validator

import (
	"strings"
)

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}
