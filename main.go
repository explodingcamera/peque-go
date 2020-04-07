package main

import (
	"fmt"

	"github.com/explodingcamera/peque-go/backends/postgres"
)

func main() {
	queue, err := postgres.Connect(postgres.Options{})
	fmt.Println(err)
	queue.Install()
}
