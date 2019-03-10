package main

import (
	"net/http"

	"github.com/ofabricio/pong"
)

func Handler(next http.HandlerFunc) http.HandlerFunc {
	return pong.Default(next)
}
