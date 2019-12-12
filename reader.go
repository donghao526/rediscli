package rediscli

import (
	"errors"
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

func ReadReply(context *RedisContext) (string, error) {

	// read line
	line := ReadLine(context)
	err := errors.New("the response invalid")
	if !strings.Contains(line, "\r\n") && !isValidPrefix(line[0]) {
		return line, err
	}

	var res = ""
	switch line[0] {
	case '+':
		return ParseSimpleString(line[1:]), nil
	case '-':
		return ParseError(line[1:]), nil
	case ':':
		return ParseInteger(line[1:]), nil
	case '$':
		return ProcessBulkString(context, line[1:]), nil
	}
	return res, nil
}

func ProcessBulkString(ctx *RedisContext, line string) (string, error) {
    strBulk := ""
    err := errors.New("read failed")
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
    intBase := 0
    intStart := 0
    bolPositive := true

    if line[0] == '-' {
        bolPositive = false
        intStart = 1
    }
    
	for i := intStart; i < len(line); i++ {
		if line[i] >= '0' && line[i] <= '9' {
			intBase = intBase*10 + (int)(line[i]-'0')
		} else {
			break
		}
	}
    
    if bolPositive == false {
        return -1 * intBase
    }
}

func isValidPrefix(prefix char) bool {
	switch prefix {
        '-' :
        '$' :
        '*' :
        '+' :
        ':' : return true;
        default: return false; 
	}
}
