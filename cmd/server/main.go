package main

import (
	"github.com/t1mon-ggg/gophkeeper/pkg/server"
)

func main() {
	server := server.New()
	server.Start()
	server.WG().Wait()
}
