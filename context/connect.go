package context

import (
	"bufio"
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

func ReadFromServer(context *RedisContext) string {
	var res = ""
	res, err := bufio.NewReader(context.Con).ReadString('\n')
	if err != nil  {
		return res
	}
	return res
}