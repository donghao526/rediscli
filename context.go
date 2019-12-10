package rediscli

import (
	"bufio"
	"net"
)

type RedisContext struct {
	ip string
	port string
	conn net.Conn
	reader *bufio.Reader
	command string
}

func GetRedisContext(ip string, port string) *RedisContext {
	var ctx = initContext(ip, port)
	_, err := ConnectToServer(ctx)
	if err != nil {
		return nil
	}
	setContextReader(ctx)
	return ctx
}

func setContextReader(ctx *RedisContext) {
	ctx.reader = bufio.NewReader(ctx.conn)
}

func initContext(ip string, port string) *RedisContext  {
	var ctx RedisContext
	ctx.ip = ip
	ctx.port = port
	ctx.conn = nil
	ctx.reader = nil
	ctx.command = ""
	return &ctx
}