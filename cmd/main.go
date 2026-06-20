package main

import (
	"fmt"

	"github.com/umizu/rail"
)

func main() {
	rcon, err := rail.Connect("127.0.0.1:25575", "secretpassword")
	if err != nil {
		panic(err)
	}

	resp, err := rcon.Send("list") // list players online
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	fmt.Println(resp.Payload)
}
