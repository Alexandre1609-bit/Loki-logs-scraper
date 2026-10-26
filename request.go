package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Root struct {
	ObtainedStatus string `json:"status"`
	Data           Data   `json:"data"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

func main() {
	/* Fonction de test
	   	Test := `{
	     		"status": "success",
	     		"data": {
	       		"resultType": "streams",
	       		"result": [
	         			{
	           			"stream": { "service_name": "checkout-api" },
	           			"values": [
	             				[ "1782008380038503000", "checkout complete trace_id=5b8..." ]
	           			]
	         			}
	       		]
	     		}
	   	}`

	   	var response Root
	   	err := json.Unmarshal([]byte(Test), &response)
	   	if err != nil {
	   		panic(err)
	   	}

	   	fmt.Println(response)
	*/
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	startTime := time.Now().Add(-5 * time.Second)

	baseURL, err := url.Parse("http://localhost:3100/loki/api/v1/query_range")
	if err != nil {
		log.Fatalln(err)
	}

	params := url.Values{}
	params.Add("query", "{job=~\".+\"}")

	for range ticker.C {
		endTime := time.Now()
		params.Set("start", strconv.FormatInt(startTime.UnixNano(), 10))
		params.Set("end", strconv.FormatInt(endTime.UnixNano(), 10))
		baseURL.RawQuery = params.Encode()
		resp, err := http.Get(baseURL.String())

		if err != nil {
			log.Print(err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			res, err := io.ReadAll(resp.Body)
			log.Print(string(res), err)
			resp.Body.Close()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			log.Print(err)
			continue
		}

		sb := string(body)
		log.Println(sb)

		log.Println(baseURL.String())

		startTime = endTime
		resp.Body.Close()
	}

}
