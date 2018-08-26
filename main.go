package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

const (
	urlCatalogsPrefix     = "https://"
	urlCatalogsSuffix     = ".api.culturacolectiva.com/catalogs"
	urlCatalogsItemSuffix = "/item"
)

func main() {
	dataFlag := flag.String("data", "https://cucodev.culturacolectiva.com/jsoncategory/", "URL for get the data (json) to add")
	token := flag.String("token", "", "Token needed for make the petition")
	environment := flag.String("environment", "dev", "Environment to make the petition {dev, staging}")
	v := flag.Bool("v", false, "Print the version of the program")
	version := flag.Bool("version", false, "Print the version of the program")
	flag.Parse()

	if *v || *version {
		fmt.Printf("go-filler version %s\n", appVersion)
		os.Exit(0)
	}

	if *token == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	urlCatalogs := urlCatalogsPrefix + *environment + urlCatalogsSuffix

	response := makePetition(http.MethodGet, urlCatalogs, nil, token)

	id := ""

	// We search for the id of a catalog which attribute type is categories
	for k, v := range response {
		if k == "data" {
			for _, data := range v.([]interface{}) {
				if data.(map[string]interface{})["attributes"].(map[string]interface{})["type"] == "categories" {
					id = data.(map[string]interface{})["id"].(string)
				}
			}
		}
	}

	if id == "" {
		body := map[string]interface{}{
			"data": map[string]interface{}{
				"type": "catalogs",
				"attributes": map[string]string{
					"name":        "categories",
					"description": "Categories of the CMS 2.0",
					"type":        "categories",
				},
			},
		}

		// Maybe this can be achieved with diferent approach, but for now, works
		bodyCasted, _ := json.Marshal(body)

		response = makePetition(http.MethodPost, urlCatalogs, bodyCasted, token)

		for k, v := range response {
			if k == "id" {
				id = v.(string)
			}
		}
	}

	responseArray := makePetitionResponseArray(http.MethodGet, *dataFlag, nil, nil)

	total := len(responseArray)

	urlCatalogsItem := urlCatalogs + urlCatalogsItemSuffix

	for k, v := range responseArray {
		v["data"].(map[string]interface{})["attributes"].(map[string]interface{})["parent"] = id
		name := v["data"].(map[string]interface{})["attributes"].(map[string]interface{})["name"]

		fmt.Printf("Processing %d of %d: Name: %s\n", k+1, total, name)
		body, _ := json.Marshal(v)
		_ = makePetition(http.MethodPost, urlCatalogsItem, body, token)
	}
}
