package rediscli

import (
    "fmt"
    "strings"
)

func ReadLine(ctx *RedisContext) string {
    res, err := ctx.reader.ReadString('\n')
    if err != nil {
        return ""
    }
    return res
}

func ReadReply(context *RedisContext) string {
    var line = ReadLine(context)
    var res = ""
    switch line[0] {
    case '+':
        return GetSimpleString(line)
    case '-':
        return GetErrorString(line)
    case ':':
        return GetIntegerString(line)
    case '$':
        return GetBulkString(context, line)
    }
    return res
}

func GetSimpleString(simpleString string) string {
    var strArray = strings.Split(simpleString, "\r\n")
    var strSimpleStringContent = strArray[0]
    var strLen = len(strSimpleStringContent)
    return strSimpleStringContent[1:strLen]
}

func GetErrorString(error string) string {
    var strArray = strings.Split(error, "\r\n")
    var strErrorStringContent = strArray[0]
    var strLen = len(strErrorStringContent)
    return "(error)" + strErrorStringContent[1:strLen]
}

func GetIntegerString(integer string) string {
    var strArray = strings.Split(integer, "\r\n")
    var strIntegerStringContent = strArray[0]
    var strLen = len(strIntegerStringContent)
    return strIntegerStringContent[1:strLen]
}

func GetBulkString(ctx *RedisContext, line string) string {
    intBulkLen := parseLen(line[1:])
    if intBulkLen > 0 {
        var strNewLine = ReadLine(ctx)
        intCrlfPos := strings.Index(strNewLine, "\r\n")
        return fmt.Sprintf("\"%s\"", strNewLine[:intCrlfPos])
    } else {
        return "(empty list or set)"
    }
}

func parseLen(line string) int {
    base := 0
    for i := 0; i < len(line); i++ {
        if line[i] >= '0' && line[i] <= '9' {
            base = base*10 + (int)(line[i]-'0')
        } else {
            break
        }
    }
    return base
}
