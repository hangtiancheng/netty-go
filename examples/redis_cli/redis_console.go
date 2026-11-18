package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/examples/redis_cli/resp"
)

type simpleRedisConsole struct {
	respChan chan *resp.Resp
}

func (s *simpleRedisConsole) HandleActive(ctx netty.ActiveContext) {
	fmt.Println("connected")

	go s.attachConsole(ctx)
	ctx.HandleActive()
}

func (s *simpleRedisConsole) HandleRead(ctx netty.InboundContext, message netty.Message) {
	s.respChan <- message.(*resp.Resp)
}

func (s *simpleRedisConsole) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	fmt.Println("disconnected", ex)

	ctx.HandleInactive(ex)
}

func (s *simpleRedisConsole) attachConsole(ctx netty.HandlerContext) {

	s.respChan = make(chan *resp.Resp, 1)

	fmt.Print(ctx.Channel().RemoteAddr(), ">")

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {

		if text := strings.TrimSpace(stdin.Text()); len(text) > 0 {

			switch text {
			case "exit":
				ctx.Close(fmt.Errorf("user exit"))
				return
			case "help":
				fmt.Println("help information")
			default:
				inputs := strings.Split(text, " ")

				// build command.
				var cmds = make([]resp.Value, 0, len(inputs))
				for _, v := range inputs {
					cmds = append(cmds, resp.BlukString(v))
				}

				// send redis command.
				ctx.Write(cmds)

				// print response
				reply := <-s.respChan

				switch {
				case reply.Null:
					fmt.Println("(empty)")
				case len(reply.Data) > 0:
					fmt.Println(reply.Data)
				case len(reply.Array) > 0:
					for index, rsp := range reply.Array {
						fmt.Println(fmt.Sprintf(`%d) "%s"`, index+1, rsp.Data))
					}
				}
			}
		} else {
			fmt.Println("please input redis command")
		}

		fmt.Print(ctx.Channel().RemoteAddr(), ">")
	}
}
