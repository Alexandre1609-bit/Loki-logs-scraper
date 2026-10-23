package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	startTime := time.Now().Add(-5 * time.Second)

	for range ticker.C {
		endTime := time.Now()
		fmt.Println(startTime, endTime)
		startTime = endTime
	}

	baseURL, err := url.Parse("http://localhost:3100/loki/api/v1/query_range")
	if err != nil {
		log.Fatalln(err)
	}

	params := url.Values{}
	params.Add("query", "{job=\"varlogs\"}")
	params.Add("start", "1609459200000000000")
	params.Add("end", "1609462800000000000")

	baseURL.RawQuery = params.Encode()

	fmt.Println(baseURL.String())

	resp, err := http.Get(baseURL.String())
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		resp.Body.Close()
	}

	sb := string(body)
	log.Println(sb)

}
