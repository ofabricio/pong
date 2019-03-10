package main

import (
	"net/http"

	"github.com/ofabricio/pong"
)

func main() {
	http.HandleFunc("/", pong.Default(http.NotFound))
	panic(http.ListenAndServe(":4000", nil))
}
