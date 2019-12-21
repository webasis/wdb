package main

import (
	"fmt"
	"log"
	"time"

	"github.com/webasis/wdb"
)

func main() {
	c := wdb.NewClient("http://localhost:9812")

	cache := ""
	for {
		raw, err := c.Get("hello")
		if err != nil || len(raw) == 0 {
			log.Print("error: retry after 1 second")
			time.Sleep(time.Second)
			continue
		}

		if cache == string(raw) {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		cache = string(raw)
		fmt.Println(cache)
	}
}
