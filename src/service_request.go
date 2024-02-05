package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type PayloadEntry struct {
	Index      int         `json:"index"`
	MethodName string      `json:"methodname"`
	Args       interface{} `json:"args"`
}

type ResponseEntry[T any] struct {
	Error bool `json:"error"`
	Data  T    `json:"data"`
}

func makeServiceRequest[T any](client *http.Client, args interface{}, sessionKey string, methodName string) T {
	payload := []PayloadEntry{
		{
			Index:      0,
			MethodName: methodName,
			Args:       args,
		},
	}

	jsonValue, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://lernen.min.uni-hamburg.de/lib/ajax/service.php", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatalf("Error while creating request %s", err)
	}

	q := req.URL.Query()
	q.Add("sesskey", sessionKey)
	q.Add("info", methodName)

	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error while getting course format %s", err)
	}
	defer res.Body.Close()

	resBodyContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error while reading course format %s", err)
	}

	var response []ResponseEntry[T]
	err = json.Unmarshal(resBodyContent, &response)
	if err != nil {
		log.Fatalf("Error while parsing course format %s", err)
	}

	if len(response) != 1 {
		log.Fatalf("Error while parsing course format, response length is not 1")
	}

	responseEntry := response[0]
	if responseEntry.Error {
		log.Fatalf("Error while parsing course format, response error is true")
	}

	return responseEntry.Data
}
