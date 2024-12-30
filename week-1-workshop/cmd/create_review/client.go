package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/pflag"
	"io"
	"log"
	"net/http"
)

type CustomRT struct {
	R http.RoundTripper
}

func (r *CustomRT) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Print("client is called")

	return r.R.RoundTrip(req)
}

func main() {
	addr := pflag.String("addr", "", "Specify the address of service")
	pflag.Parse()

	if *addr == "" {
		fmt.Println("addr is empty")
		return
	}

	client := http.Client{
		Transport: &CustomRT{R: http.DefaultTransport},
		Timeout:   0,
	}

	var createReviewRequest struct {
		SKU     int64     `json:"sku"`
		Comment string    `json:"comment"`
		UserID  uuid.UUID `json:"user_id"`
	}

	createReviewRequest.SKU = 10
	createReviewRequest.Comment = "Comment"
	createReviewRequest.UserID = uuid.New()

	body, err := json.Marshal(createReviewRequest)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(body)

	respo, err := client.Post(*addr, "application/json", r)
	if err != nil {
		panic(err)
	}

	defer respo.Body.Close()

	responseBody, err := io.ReadAll(respo.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(responseBody))

	return
}
