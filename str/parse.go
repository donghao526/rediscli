package str

import (
	"strings"
)

func ParseUserCommand(command string) []string {
	var keywords = strings.Split(command, " ")
	return keywords
}

func ParseServerResponse(response string) string {
	var res = ""
	switch response[0] {
		case '+':
			return getSimpleString(response)
		case '-' :
			return getErrorString(response)
		case ':' :
			return getIntegerString(response)
	}
	return res
}

func getSimpleString(simpleString string) string {
	var strArray = strings.Split(simpleString, "\r\n")
	var strSimpleStringContent = strArray[0]
	var strLen = len(strSimpleStringContent)
	return strSimpleStringContent[1 : strLen]
}

func getErrorString(error string) string {
	var strArray = strings.Split(error, "\r\n")
	var strErrorStringContent = strArray[0]
	var strLen = len(strErrorStringContent)
	return "(error)" + strErrorStringContent[1 : strLen]
}

func getIntegerString(integer string) string {
	var strArray = strings.Split(integer, "\r\n")
	var strIntegerStringContent = strArray[0]
	var strLen = len(strIntegerStringContent)
	return strIntegerStringContent[1 : strLen]
}