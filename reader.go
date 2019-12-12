package rediscli

import (
    "fmt"
    "strings"
    "errors"
)

func ReadLine(ctx *RedisContext) string {
    res, err := ctx.reader.ReadString('\n')
    if err != nil {
        return ""
    }
    return res
}

func ReadReply(context *RedisContext) (string, error) {

    // read line
    line := ReadLine(context)
    err := errors.New("the response invalid")
    if !strings.Contains(line, "\r\n") {
        return line, err
    }

    var res = ""
    switch line[0] {
    case '+':
        return ParseSimpleString(line), nil
    case '-':
        return ParseError(line), nil
    case ':':
        return ParseInteger(line), nil
    case '$':
        return GetBulkString(context, line), nil
    }
    return res, nil
}

func ParseInteger(integer string) string {
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
