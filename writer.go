package rediscli

import "fmt"

func SendCommand(ctx *RedisContext) (bool, error) {
	strRespCommand := GetRespStrOfCmd(ctx.command)
	var _, err = fmt.Fprintf(ctx.conn, strRespCommand)
	if err != nil {
		fmt.Println("write to server fail")
		return false, err
	}
	return true, err
}
