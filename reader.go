package rediscli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type RedisReader struct {
	err     int    /* Error flags, 0 when there is no error */
	errstr  string /* String representation of error when applicable */
	pos     int    /* Buffer cursor */
	maxbuf  int    /* Max length of unused buffer */
	ridx    int    /* read index for the stack*/
	buf     [1024 * 16]byte
	len     int
	cur_pos int
	rstack  [10]RedisReaderTask /*read stack*/
}

type RedisReaderTask struct {
	elements int
	idx      int
	parent   *RedisReaderTask
	obj_type int
	obj      *RedisObject
}

type RedisObject struct {
	obj_type  int
	int_value int
	str_value string
}

/*
 * @brief process item
 * @param RedisReader
 */
func processItem(r *RedisReader) int {
	cur := r.rstack[r.ridx]
	p, err := readChar(r)
	if err == REDIS_OK {
		switch p {
		case '-':
			cur.obj_type = TYPE_ERROR
			break
		case '+':
			cur.obj_type = TYPE_STRING
			break
		case ':':
			cur.obj_type = TYPE_INTEGER
			break
		default:
		}
	} else {
		return REDIS_ERR
	}

	switch cur.obj_type {
	case TYPE_INTEGER:
	case TYPE_ERROR:
	case TYPE_STRING:
		processLineItem(r)
	}
	return REDIS_OK
}

func processLineItem(r *RedisReader) int {
	task := r.rstack[r.ridx]
	strLine, err := readLine(r)
	if err == REDIS_ERR {
		return REDIS_ERR
	}

	if task.obj_type == TYPE_STRING {
		obj := RedisObject{
			obj_type:  TYPE_STRING,
			int_value: 0,
			str_value: strLine,
		}
		task.obj = &obj
	}
	moveToNextTask(r)
	return REDIS_OK
}

func moveToNextTask(r *RedisReader) int {
	if r.ridx == 0 {
		r.ridx = -1
		return REDIS_OK
	}
	return REDIS_OK
}

func readBytes(r *RedisReader, bytes int) []byte {
	if r.len - r.cur_pos >= bytes {
		t := r.buf[r.cur_pos : r.cur_pos+bytes]
		r.cur_pos += bytes
		return t
	}
	return nil
}

func readLine(r *RedisReader) (string, int) {
	newLinePos := seekNewLine(r)
	if newLinePos == -1 {
		return "", REDIS_ERR
	}

	strNewLine := string(r.buf[r.cur_pos : newLinePos])
	r.cur_pos += newLinePos - r.cur_pos + 2
	return strNewLine, REDIS_OK
}

func seekNewLine(r *RedisReader) int{
	pos := r.cur_pos
	for pos <= r.len - 1 {
		if r.buf[pos] == '\r' {
			if (pos < r.len - 1) && (r.buf[pos + 1] == '\n') {
				return pos
			}
		}
		pos++
	}
	return -1
}

func readChar(r *RedisReader) (byte, int) {
	if r.len >= r.cur_pos+1 {
		t := r.buf[r.cur_pos]
		r.cur_pos++
		return t, REDIS_OK
	}
	return '0', REDIS_ERR
}

func RedisGetReply(ctx *RedisContext) int {
	if REDIS_OK != ReadRedisReply(ctx) {
		return REDIS_ERR
	}

	r := ctx.replyReadr
	r.ridx = 0
	r.rstack[0].elements = -1
	r.rstack[0].obj_type = -1
	r.rstack[0].idx = -1
	r.rstack[0].parent = nil

	for r.ridx >= 0 {
		if processItem(r) != REDIS_OK {
			break
		}
	}

	return REDIS_OK
}

/*
 * @brief read redis reply to reader
 */
func ReadRedisReply(ctx *RedisContext) int {
	var newbuf [1024 * 16]byte
	var nread, _ = ctx.reader.Read(newbuf[:])
	fmt.Println(string(newbuf[:]))
	var newReader RedisReader

	if nread > 0 {
		newReader.len = nread
		newReader.cur_pos = 0
		copy(newReader.buf[:], newbuf[:])
		ctx.replyReadr = &newReader
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
