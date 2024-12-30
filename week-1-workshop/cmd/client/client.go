package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"net/http"
)

func main() {
	addr := pflag.String("addr", "", "Specify the address of service")
	pflag.Parse()

	if *addr == "" {
		fmt.Println("addr is empty")
		return
	}

	response, err := http.Get(*addr)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	return
}
