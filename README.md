# Coal

A minecraft RCON client library

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/umizu/coal"
)

func main() {
	client, err := coal.Connect("localhost:25575", "password")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	resp, err := client.Send("list")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Content)
}
```

## Known Limitations

Maximum response length of 4096 bytes. (https://minecraft.wiki/w/RCON#fragmentation)
