package str

import (
	"strings"
)

func ParseUserCommand(command string) []string {
	var keywords = strings.Split(command, " ")
	return keywords
}