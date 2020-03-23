package rediscli

// create a string object
func CreateStringObject(str string) *RedisObject  {
	strObj := new(RedisObject)
	strObj.str_value = str
	strObj.obj_type = TYPE_STRING
	return strObj
}

// create an error redis object
func CreateErrorObject(err string) *RedisObject{
	strObj := new(RedisObject)
	strObj.str_value = err
	strObj.obj_type = TYPE_ERROR
	return strObj
}