package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func catalogsLogic() {
	urlCatalogs := urlPrefix + *environment + urlSuffix + urlCatalogsSuffix

	response := makePetition(http.MethodGet, urlCatalogs, nil, token)

	id := ""

	// We search for the id of a catalog which attribute type is categories
	for k, v := range response {
		if k == "data" {
			for _, data := range v.([]interface{}) {
				if data.(map[string]interface{})["attributes"].(map[string]interface{})["type"] == *typeContent {
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
					"name":        *name,
					"description": *description,
					"type":        *typeContent,
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
