package context

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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
	var line = ""
	var err error
	var reader = bufio.NewReader(context.Con)
	for {
		line, err = reader.ReadString('\n')
		if len(line) == 0 || err != nil {
			return ""
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			break
		}
	}
	fmt.Println(line)
	return line
}