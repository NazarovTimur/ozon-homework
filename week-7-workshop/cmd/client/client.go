package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func main() {
	var request struct {
		Key   string
		Value string
	}

	for i := 0; ; i++ {
		request.Value = strconv.Itoa(i)
		request.Key = strconv.Itoa(i)

		body, _ := json.Marshal(request)

		r, err := http.Post("http://localhost:8080/set", "application/json", bytes.NewReader(body))
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(r.Body)
		r.Body.Close()

		time.Sleep(time.Nanosecond)
	}
}
