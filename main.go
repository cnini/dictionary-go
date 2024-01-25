package main

import (
	"dictionary-go/server"
	"fmt"
	"strings"
)

func handleErrors(err error) {
	switch {
	case strings.Contains(err.Error(), "Permission denied"):
		fmt.Println("\nERROR | Permission denied:", err)

	case strings.Contains(err.Error(), "Bad request"):
		fmt.Println("\nERROR | Bad request:", err)

	default:
		fmt.Println("\nERROR:", err)
	}
}

func main() {
	errors := make(chan error)

	go func() {
		defer close(errors)

		for err := range errors {
			handleErrors(err)
		}
	}()

	server.Start(errors)
	server.ListenAndServe()
}
