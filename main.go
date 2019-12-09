package main

import (
	"bufio"
	"fmt"
	"github.com/donghao526/rediscli/context"
	"github.com/donghao526/rediscli/str"
	"os"
	"strings"
)


func main()  {
	ctx := context.GetRedisContext("127.0.0.1", "6379")
	for true {
		cmdLinePrompt(ctx)
		var input = getInput()
		var words = str.ParseUserCommand(input)
		var res = str.BuildRespStr(words)
		context.WriteToServer(res, ctx)
		var response = context.ReadFromServer(ctx)
		fmt.Println(response)
		fmt.Println(str.ParseServerResponse(response))
	}
}

func getInput() string {
	var input = ""
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return input
	}
	var words = strings.Split(input, "\n")
	return words[0]
}

func cmdLinePrompt(ctx *context.RedisContext) {
	fmt.Printf("%s:%s> ", ctx.Ip, ctx.Port)
}
