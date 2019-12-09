package context

import "net"

type RedisContext struct {
	Ip string
	Port string
	Con net.Conn
}

func GetRedisContext(ip string, port string) *RedisContext {
	var context = InitContext(ip, port)
	ConnectServer(context)
	return context
}

func InitContext(ip string, port string) *RedisContext  {
	var context RedisContext
	context.Ip = ip
	context.Port = port
	context.Con = nil
	return &context
}

