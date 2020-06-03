package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kichyr/cubeCliClient/test/testserver"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("wrong number of cli args: \n testserver <port>")
		os.Exit(1)
	}
	srv := testserver.NewTestServer("./test/testserver/test_tokens.json")
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("wrong number format of port: %s", os.Args[1])
		os.Exit(1)
	}
	srv.StartServer(port) // start test server on localhost:8091
}
