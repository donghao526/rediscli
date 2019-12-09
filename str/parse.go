package str

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseUserCommand(command string) []string {
	var keywords = strings.Fields(command)
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
		case '$' :
			return getBulkString(response)
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

func getBulkString(arrayString string) string {
	fmt.Println(len(arrayString))
	var intLastCrlf  = strings.LastIndex(arrayString, "\r\n")
	fmt.Println(intLastCrlf)
	var intFirstCrlf = strings.Index(arrayString, "\r\n")
	fmt.Println(intFirstCrlf)
	var strElements = arrayString[1 : intFirstCrlf - 1]
	fmt.Println(strElements)

	var intElement, _ = strconv.Atoi(strElements)
	if intElement < 0 {
		return "nil"
	} else {
		if intFirstCrlf + 2 > intLastCrlf - 2 {
			return ""
		} else {
			return arrayString[intFirstCrlf + 2 : intLastCrlf - 2 ]
		}
	}
}

func getArrayString(bulk string) string {
	var intFirstCrlf = strings.Index(bulk, "\r\n")
	var strElements = bulk[1 : intFirstCrlf - 1]
	var intElement, _ = strconv.Atoi(strElements)

	if intElement > 0 {
		return ""
	} else {
		return "(empty list or set)"
	}
}