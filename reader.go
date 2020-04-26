package rediscli

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
	size int
	member []*RedisObject
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
		case '*':
			cur.obj_type = TYPE_ARRAY
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
	case TYPE_ARRAY:
		processAggregate(r)
	}
	return REDIS_OK
}

// process aggregate types, such as map, set, hash
// but map should double on elements, because of key and value
func processAggregate(r *RedisReader) int {
	task := &r.rstack[r.ridx]

	// get the size of the aggregate object
	size := readLen(r)
	task.elements = size

	// step over "\r\"
	r.cur_pos += 2
	root := r.ridx == 0

	if task.obj_type == TYPE_ARRAY {
		CreateArrayObject(task, size)
	}

	if size > 0 {
		r.ridx++
		r.rstack[r.ridx].elements = - 1
		r.rstack[r.ridx].idx = 0
		r.rstack[r.ridx].parent = task
	}

	if root {
		r.reply = task.obj
	}
	return REDIS_OK
}

// process line item
func processLineItem(r *RedisReader) int {
	task := &r.rstack[r.ridx]

	if task.obj_type == TYPE_STRING || task.obj_type == TYPE_ERROR {
		strLine, err := readLine(r)
		if err == REDIS_ERR {
			return REDIS_ERR
		}
		if task.obj_type == TYPE_STRING {
			CreateStringObject(task, strLine)
		}
		if task.obj_type == TYPE_ERROR {
			CreateErrorObject(task, strLine)
		}
	} else if task.obj_type == TYPE_BULK || task.obj_type == TYPE_INTEGER{
		strLen := readLen(r)
		if strLen == -1 {
			CreateNilObject(task)
		} else if task.obj_type == TYPE_BULK {
			r.cur_pos += 2
			bulk := readBytes(r, strLen)
			CreateBulkObject(task, string(bulk[:]))
		} else if task.obj_type == TYPE_INTEGER {
			CreateIntegerObject(task, strLen)
		}
		r.cur_pos += 2
	}

	if r.ridx == 0 {
		r.reply = task.obj
	}
	moveToNextTask(r)
	return REDIS_OK
}

/**
 * @brief move to next protocol parse task
 * @param RedisReader
 * @return Error Code
 */
func moveToNextTask(r *RedisReader) int {
	for r.ridx >=0 {
		if r.ridx == 0 {
			r.ridx = -1
			return REDIS_OK
		}
		cur := &r.rstack[r.ridx]
		prv := &r.rstack[r.ridx - 1]
		if prv.obj_type != TYPE_ARRAY {
			return REDIS_ERR
		}

		if cur.idx == prv.elements - 1 {
			r.ridx--
		} else {
			cur.idx++
			cur.obj_type = -1
			cur.elements = -1
			return REDIS_OK
		}
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
	flag := 1
	for pos < r.len {
		if r.buf[pos] > 0x30 && r.buf[pos] <= 0x39 {
			res = res * 10 + (int)(r.buf[pos] - 0x30)
			pos = pos + 1
		} else if r.buf[pos] == 45 {
			flag = -1
			pos = pos + 1
		} else {
			r.cur_pos += pos - r.cur_pos
			res = res * flag
			return res
		}
	}
	return pos
}

// read line
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

// read a byte
func readChar(r *RedisReader) (byte, int) {
	if r.len >= r.cur_pos+1 {
		t := r.buf[r.cur_pos]
		r.cur_pos++
		return t, REDIS_OK
	}
	return '0', REDIS_ERR
}

// get redis reply
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