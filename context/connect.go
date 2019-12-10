package context

import (
	"fmt"
	"net"
)

func ConnectServer(context *RedisContext) {
	conn, err := net.Dial("tcp", context.Ip + ":" +context.Port)
	if err != nil {
		return
	}

	fmt.Println("success")
	context.Con = conn
	return
}

func WriteToServer(command string, context *RedisContext) {
	var conn = context.Con
	var _, err = fmt.Fprintf(conn, command)
	if err != nil {
		fmt.Println("write to server fail")
	}
}

