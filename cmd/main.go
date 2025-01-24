package main

import "github.com/popliop/goPage/pkg/server"

func main() {
	apiServer := server.NewAPIServer("localhost:80")
	apiServer.Run()
}
