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

// parse simple string
func ParseSimpleString(simpleString string) string {
	strArray := strings.Split(simpleString, "\r\n")
	strContent := strArray[0]
	return strContent
}

// parse error
func ParseError(error string) string {
	strArray := strings.Split(error, "\r\n")
	strErrorContent := strArray[0]
	return "(error)" + strErrorContent
}

// parse integer
func ParseInteger(integer string) string {
	strArray := strings.Split(integer, "\r\n")
	strIntContent := strArray[0]
	return strIntContent
}
