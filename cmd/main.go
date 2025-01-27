package main

import (
	"fmt"
	"os"

	"github.com/popliop/goPage/pkg/server"
)

func main() {
	fmt.Println("v4")
	data := os.Getenv("GPT_API_KEY")

	fmt.Println(data)

	apiServer := server.NewAPIServer(":80")
	apiServer.Run()

}
