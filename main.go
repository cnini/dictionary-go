package main

import (
	"dictionary-go/server"
)

func main() {
	server.Start()
	server.ListenAndServe()
}
