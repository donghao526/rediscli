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
	var intCount = len(arrKeyWords)
	var strResp = ""
	strResp += fmt.Sprintf("*%d\r\n", intCount)

	for _, value := range arrKeyWords {
		if value != "" {
			strResp += fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
		}
	}
	return strResp
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
		out += fmt.Sprintf("%d", reply.int_value)
	}
	return out
}