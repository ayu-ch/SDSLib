package main

import (
	"fmt"

	"github.com/ayu-ch/SDSLib/pkg/api"
)

func main() {
	fmt.Println("Started the API server")
	api.Start()
}
