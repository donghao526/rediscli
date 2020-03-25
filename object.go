package rediscli

// create a string object
func CreateStringObject(task *RedisReaderTask, str string) {
	objStr := new(RedisObject)
	objStr.str_value = str
	objStr.obj_type = TYPE_STRING
	if task.parent != nil {
		task.parent.obj.member[task.idx] = objStr
	} else {
		task.obj = objStr
	}
}

// create an error redis object
func CreateErrorObject(task *RedisReaderTask, err string) {
	objErr := new(RedisObject)
	objErr.str_value = err
	objErr.obj_type = TYPE_ERROR
	task.obj =  objErr
}

// create an error redis object
func CreateBulkObject(task *RedisReaderTask, bulk string) {
	objBulk := new(RedisObject)
	objBulk.str_value = bulk
	objBulk.obj_type = TYPE_BULK
	if task.parent != nil {
		task.parent.obj.member[task.idx] = objBulk
	} else {
		task.obj =  objBulk
	}
}

// create an integer object
func CreateIntegerObject(task *RedisReaderTask, value int) {
	objInt := new(RedisObject)
	objInt.int_value = value
	objInt.obj_type = TYPE_INTEGER
	task.obj =  objInt
}

//  create a nil object
func CreateNilObject(task *RedisReaderTask) {
	objNil := new(RedisObject)
	objNil.str_value = "nil"
	objNil.obj_type = TYPE_NIL
	if task.parent != nil {
		task.parent.obj.member[task.idx] = objNil
	} else {
		task.obj =  objNil
	}
}

// create a array object
func CreateArrayObject(task *RedisReaderTask, size int) {
	objArray := new(RedisObject)
	objArray.obj_type = TYPE_ARRAY
	objArray.size = size
	objArray.member = make([]*RedisObject, size)
	task.obj = objArray
}