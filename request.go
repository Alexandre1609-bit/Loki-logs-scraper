package main

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

	mask()
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
	   		panic(err) Panic ici mais on récupère l'erreur dans notre boucle en dessous
	   	}

	   	fmt.Println(response)
	#########
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	startTime := time.Now().Add(-5 * time.Second)

	baseURL, err := url.Parse("http://localhost:3100/loki/api/v1/query_range")
	if err != nil {
		log.Fatalln(err)
	}

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

		err = json.Unmarshal(body, &RootJson)
		if err != nil {
			sb := string(body)
			log.Println(sb)
			resp.Body.Close()
			log.Print(err)
			continue
		}

		for _, res := range RootJson.Data.Result {
			log.Println(res.Stream["namespace"], res.Stream["pod"])
			for _, v := range res.Values {
				log.Println(v[0], v[1])
				v[1] = strings.TrimSpace(v[1])
			}
		}

		startTime = endTime
		resp.Body.Close()
	}
	*/
}
