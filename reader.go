package rediscli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type redisReader struct {
	err    int    /* Error flags, 0 when there is no error */
	errstr string /* String representation of error when applicable */
	buf    string /* Read buffer */
	pos    int    /* Buffer cursor */
	len    int    /* Buffer length */
	maxbuf int    /* Max length of unused buffer */
	ridx   int
	rtask  [10]RedisContext
}

type redisReaderTask struct {
}

func processItem() int {
	return REDIS_OK
}

func ReadBuffer(ctx *RedisContext) int {
	var newbuf [1024 * 16]byte
	var nread, _ = ctx.reader.Read(newbuf[:])

	if nread > 0 {
		ctx.len = nread
		copy(ctx.buf[:], newbuf[:])
		return REDIS_OK
	} else if nread < 0 {
		return REDIS_ERR
	}

	return REDIS_OK
}

/*
 * read line from server
 */
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

func ProcessArray(ctx *RedisContext, line string) (string, error) {
	return "", nil
}

/*
 * process the bulk string
 */
func ProcessBulkString(ctx *RedisContext, line string) (string, error) {
	intBulkLen := parseLen(line)
	if intBulkLen > 0 {
		var strNewLine = ReadLine(ctx)
		intCrlfPos := strings.Index(strNewLine, "\r\n")
		return fmt.Sprintf("\"%s\"", strNewLine[:intCrlfPos]), nil
	} else if intBulkLen == 0 {
		return "", nil
	} else {
		return "nil", nil
	}
}

/*
 * get the length of the bulk string
 */
func parseLen(line string) int {
	strArray := strings.Split(line, "\r\n")
	intLen, _ := strconv.Atoi(strArray[0])
	return intLen
}

/*
 * check the prefix is valid
 */
func isValidPrefix(prefix byte) bool {
	switch prefix {
	case '-':
	case '$':
	case '*':
	case '+':
	case ':':
		return true
	default:
		return false
	}
	return false
}
