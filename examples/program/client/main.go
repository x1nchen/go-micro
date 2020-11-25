package main

import (
	"fmt"

	"github.com/asim/nitro/app"
	"github.com/asim/nitro/examples/program/types"
)

func main() {
	// create a new client program
	cli := app.New()
	cli.Name("client")

	var rsp types.Response
	// execute a function call
	err := cli.Execute("helloworld", "Handler.Call", &types.Request{Name: "Alice"}, &rsp)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print the response
	fmt.Println(rsp.Message)
}
