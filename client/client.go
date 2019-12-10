package main

import (
	"fmt"
	"github.com/donghao526/rediscli"
)

func main()  {
	ctx := rediscli.GetRedisContext("127.0.0.1", "6379")

	if ctx == nil {
		fmt.Println("get context failed.")
		return
	} else {
		for true {
			cmd := rediscli.GetUserInputCmd(ctx)
			if cmd == "" {
				continue
			} else {

				status, err := rediscli.SendCommand(ctx)
				if status == false {
					fmt.Println(err)
					continue
				} else {
					strReply := rediscli.ReadReply(ctx)
					rediscli.PrintReply(strReply)
				}
			}
		}
	}
}


