package rediscli

import (
	"fmt"
	"strings"
)

// get the resp protocol string for the command
func GetRespStrOfCmd(command string) string {
	arrKeyWords := getKeyWordsOfCmd(command)
	return buildRespString(arrKeyWords)
}

// get key words of the cmd
func getKeyWordsOfCmd(command string) []string {
	var keywords = strings.Fields(command)
	return keywords
}

// build the resp protocol string for the command
func buildRespString(arrKeyWords []string) string {
	argc := len(arrKeyWords)
	var respStr = ""
	respStr += fmt.Sprintf("*%d\r\n", argc)

	for _, value := range arrKeyWords {
		if value != "" {
			respStr += fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
		}
	}
	return respStr
}

// parse reply
func ParseReply(reply *RedisObject) string {
	out := ""
	switch reply.obj_type {
	case TYPE_STRING:
		fallthrough
	case TYPE_ERROR:
		out += reply.str_value
	case TYPE_BULK:
		out += fmt.Sprintf("\"%s\"", reply.str_value)
	case TYPE_INTEGER:
		out += fmt.Sprintf("(integer) %d", reply.int_value)
	case TYPE_NIL:
		out += fmt.Sprintf("(integer) %s", reply.str_value)
	case TYPE_ARRAY:
		out = parseArray(reply)
	}
	return out
}

func parseArray(reply *RedisObject) string {
	out := ""
	i := 0
	format := "%d) %s\n"
	for i < reply.size {
		if i == reply.size - 1 {
			format = "%d) %s"
		}
		out += fmt.Sprintf(format, i + 1, ParseReply(reply.member[i]))
		i++
	}
	return out
}