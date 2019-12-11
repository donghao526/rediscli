package rediscli

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func getInputLine() string {
    var input = ""
    inputReader := bufio.NewReader(os.Stdin)
    input, err := inputReader.ReadString('\n')
    if err != nil {
        return input
    }
    var words = strings.Split(input, "\n")
    return words[0]
}

func cmdLinePrompt(ctx *RedisContext) {
    fmt.Printf("%s:%s> ", ctx.ip, ctx.port)
}

// 获取用户输入的命令
func GetUserInputCmd(ctx *RedisContext) string {
    cmdLinePrompt(ctx)
    command := getInputLine()
    ctx.command = command
    if command == "" || len(strings.Fields(command)) == 0 {
        return ""
    } else {
        return command
    }
}

func PrintReply(reply string) {
    fmt.Println(reply)
}
