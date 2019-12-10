package rediscli

import (
	"fmt"
	"net"
)

func ConnectToServer(ctx *RedisContext) (bool, error) {
	conn, err := net.Dial("tcp", ctx.ip + ":" + ctx.port)
	if err != nil {
		return false, err
	}

	fmt.Println("success")
	ctx.conn = conn
	return true, err
}
