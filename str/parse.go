package str

import (
	"strings"
)

func ParseUserCommand(command string) []string {
	var keywords = strings.Fields(command)
	return keywords
}


