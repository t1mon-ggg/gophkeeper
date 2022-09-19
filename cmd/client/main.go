package main

import (
	"github.com/t1mon-ggg/gophkeeper/pkg/client"
)

func main() {
	client := client.New()
	client.Start()
	client.WG().Wait()
}
