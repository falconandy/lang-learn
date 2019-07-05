package main

import (
	"github.com/falconandy/lang-learn/server"
)

func main() {
	srv := server.NewServer()
	srv.Start()
}
