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


