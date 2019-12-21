package main

import (
	"net/http"

	"github.com/webasis/wdb"
)

func main() {
	s := wdb.NewServer()
	http.Handle("/", s)
	http.ListenAndServe("localhost:9812", nil)
}
