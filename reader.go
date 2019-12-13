package rediscli

import (
	"errors"
	"fmt"
    "strings"
    "strconv"
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
	if !strings.Contains(line, "\r\n") && !isValidPrefix(line[0]) {
		return line, err
	}

    if line[0] == '+' {
        return ParseSimpleString(line[1:]), nil        
    } else if line[0] == '-' {
        return ParseError(line[1:]), nil
    } else if line[0] == ':' {
        return ParseInteger(line[1:]), nil
    } else if line[0] == '$' {
        res, errBulk := ProcessBulkString(context, line[1:])
        if errBulk != nil {
            return res, errBulk
        } else {
            return res, nil
        }
    } else if line[0] == '*' {
        return "", nil
    }
	return "", nil
}

func ProcessBulkString(ctx *RedisContext, line string) (string, error) {
	intBulkLen := parseLen(line)
	if intBulkLen > 0 {
		var strNewLine = ReadLine(ctx)
		intCrlfPos := strings.Index(strNewLine, "\r\n")
		return fmt.Sprintf("\"%s\"", strNewLine[:intCrlfPos]), nil
	} else if intBulkLen == 0{
		return "", nil
	} else {
        return "nil", nil
    }
}

func parseLen(line string) int {
    strArray := strings.Split(line, "\r\n")
    intLen, _ := strconv.Atoi(strArray[0])
    return intLen
}

func isValidPrefix(prefix byte) bool {
	switch prefix {
        case '-' :
        case '$' :
        case '*' :
        case '+' :
        case ':' : return true;
        default: return false; 
    }
    return false
}
