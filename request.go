package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	lokiAddr := flag.String("loki", "http://localhost:3100/loki/api/v1/query_range", "Loki query_range endpoint")
	flag.Parse()
	baseURL, err := url.Parse(*lokiAddr)
	if err != nil {
		log.Fatalln(err)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	startTime := time.Now().Add(-5 * time.Second)

	params := url.Values{}
	params.Add("query", `{namespace=~".+", namespace!="loki"}`)
	params.Add("direction", "forward")

	for range ticker.C {
		var RootJson Root
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
			errBody, err := io.ReadAll(resp.Body)
			log.Print(string(errBody), err)
			resp.Body.Close()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			log.Print(err)
			continue
		}

		err = json.Unmarshal(body, &RootJson)
		if err != nil {
			stringBody := string(body)
			log.Println(stringBody)
			resp.Body.Close()
			log.Print(err)
			continue
		}

		for _, res := range RootJson.Data.Result {
			log.Println(res.Stream["namespace"], res.Stream["pod"])
			for _, v := range res.Values {
				v[1] = strings.TrimSpace(v[1])
				log.Println(v[0], maskIp(v[1]))
			}
		}

		startTime = endTime
		resp.Body.Close()
	}

}
