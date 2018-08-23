package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func makePetition(method, url string, body []byte, token *string) map[string]interface{} {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	if token != nil {
		req.Header.Add("Authorization", *token)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := make(map[string]interface{})

	err = json.Unmarshal(bodyResponse, &response)
	if err != nil {
		log.Fatal(err)
	}

	for k := range response {
		if k == "errors" {
			data, _ := json.Marshal(response)
			log.Fatalf("The server has responded with: %s", data)
		}
	}

	return response
}

// The unique difference between this function and the `makePetittion` is the response and what is downloaded
// These functionalities can be made by only one function
func makePetitionResponseArray(method, url string, body []byte, token *string) []map[string]interface{} {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	if token != nil {
		req.Header.Add("Authorization", *token)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := make([]map[string]interface{}, 0)

	err = json.Unmarshal(bodyResponse, &response)
	if err != nil {
		log.Fatal(err)
	}

	return response
}
