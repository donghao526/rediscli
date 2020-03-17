package main

import (
	"fmt"
	"github.com/donghao526/rediscli"
)

func main() {
	ctx := rediscli.GetRedisContext("127.0.0.1", "6379")

	if ctx == nil {
		fmt.Println("get context failed.")
		return
	} else {
		for true {
			cmd := rediscli.GetUserInputCmd(ctx)
			if cmd == "" {
				continue
			} else if cmd == "exit" {
				return
			} else {
				status, err := rediscli.SendCommand(ctx)
				if status == false {
					fmt.Println(err)
					continue
				} else {
					//rediscli.ReadReply(ctx)
//					err := rediscli.ReadRedisReply(ctx)
					rediscli.RedisGetReply(ctx)
					//if err != rediscli.REDIS_OK {
					//	fmt.Println(err)
					//	continue
					//}
					//rediscli.PrintReply(strReply)
				}
			}
		}
	}
}
