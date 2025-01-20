package main

import "github.com/popliop/goPage/pkg/server"

func main() {
	apiServer := server.NewAPIServer(":80")
	apiServer.Run()
}
