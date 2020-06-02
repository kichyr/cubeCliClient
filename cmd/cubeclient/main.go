package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kichyr/cubeCliClient/pkg/cubeclient"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Print("wrong number of cli args: \n cubeclient <host> <port> <token> <scope>")
		os.Exit(1)
	}
	host := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("wrong format of port: %s \n", os.Args[2])
		os.Exit(1)
	}
	token := os.Args[3]
	scope := os.Args[4]
	checkResult, err := cubeclient.CheckToken(host, port, token, scope)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Print(checkResult)
}
