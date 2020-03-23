package rediscli

import (
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
	reply   *RedisObject
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
	cur := &r.rstack[r.ridx]
	p, err := readChar(r)
	if err == REDIS_OK {
		switch p {
		case '-':
			cur.obj_type = TYPE_ERROR
		case '+':
			cur.obj_type = TYPE_STRING
		case ':':
			cur.obj_type = TYPE_INTEGER
		case '$':
			cur.obj_type = TYPE_BULK
		default:
		}
	} else {
		return REDIS_ERR
	}

	switch cur.obj_type {
	case TYPE_INTEGER:
		fallthrough
	case TYPE_ERROR:
		fallthrough
	case TYPE_STRING:
		fallthrough
	case TYPE_BULK:
		processLineItem(r)
	}
	return REDIS_OK
}

func processLineItem(r *RedisReader) int {
	task := &r.rstack[r.ridx]

	if task.obj_type == TYPE_STRING || task.obj_type == TYPE_ERROR {
		strLine, err := readLine(r)
		if err == REDIS_ERR {
			return REDIS_ERR
		}
		if task.obj_type == TYPE_STRING {
			task.obj = CreateStringObject(strLine)
		}
		if task.obj_type == TYPE_ERROR {
			task.obj = CreateErrorObject(strLine)
		}
	} else if task.obj_type == TYPE_BULK {
		strLen := readLen(r)
		r.cur_pos  = r.cur_pos + 2
		bulk := readBytes(r, strLen)
		task.obj = CreateBulkObject(string(bulk[:]))
	}

	if r.ridx == 0 {
		r.reply = task.obj
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

// read bytes in redis read buf
func readBytes(r *RedisReader, bytes int) []byte {
	if r.len - r.cur_pos >= bytes {
		t := r.buf[r.cur_pos : r.cur_pos+bytes]
		r.cur_pos += bytes
		return t
	}
	return nil
}

// read the len of the bulk or elements
func readLen(r *RedisReader) int {
	pos := r.cur_pos
	res := 0
	for pos < r.len {
		if r.buf[pos] > 0x30 && r.buf[pos] <= 0x39 {
			res = res * 10 + (int)(r.buf[pos] - 0x30)
			pos = pos + 1
		} else {
			r.cur_pos += pos - r.cur_pos
			return res
		}
	}
	return pos
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

func RedisGetReply(ctx *RedisContext) (int, *RedisObject) {
	var reply *RedisObject
	reply = nil
	if REDIS_OK != ReadRedisReply(ctx) {
		return REDIS_ERR, reply
	}

	r := ctx.replyReader
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

	if r.ridx == -1 {
		reply = r.reply
	}
	return REDIS_OK, reply
}

/*
 * @brief read redis reply to reader
 */
func ReadRedisReply(ctx *RedisContext) int {
	var newbuf [1024 * 16]byte
	var nread, _ = ctx.reader.Read(newbuf[:])
	var newReader RedisReader

	if nread > 0 {
		newReader.len = nread
		newReader.cur_pos = 0
		copy(newReader.buf[:], newbuf[:])
		ctx.replyReader = &newReader
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
