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
    var strArray = strings.Split(simpleString, "\r\n")
    var strSimpleStringContent = strArray[0]
    var strLen = len(strSimpleStringContent)
    return strSimpleStringContent[1:strLen]
}

// parse error
func ParseError(error string) string {
    var strArray = strings.Split(error, "\r\n")
    var strErrorStringContent = strArray[0]
    var strLen = len(strErrorStringContent)
    return "(error)" + strErrorStringContent[1:strLen]
}