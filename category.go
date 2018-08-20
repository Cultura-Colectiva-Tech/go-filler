package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func getCatalogs() {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, urlCatalogs, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", *token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &catalogs)
}

func createCategoryCatalog() {
	client := &http.Client{}

	catalogRequest := CatalogRequest{
		Data: DataRequest{
			Type: "catalogs",
			Attributes: Attributes{
				Name:        "categories",
				Description: "Categories of the CMS",
				Type:        "categories",
			},
		},
	}

	data, err := json.Marshal(catalogRequest)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, urlCatalogs, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", *token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &catalog)
}
