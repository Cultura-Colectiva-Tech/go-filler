package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

func articlesLogic() {
	url := urlPrefix + *environmentFlag + urlSuffix + *typeContentFlag

	body := "month=" + *monthFlag + "&year=" + *yearFlag + "&type=" + *typePostFlag

	// Get all the articles
	responseArray := makePetitionResponseArray(http.MethodPost, *dataFlag, []byte(body), nil)

	total := len(responseArray)

	if total < 1 {
		fmt.Println("There is no data to be processed")
		os.Exit(0)
	}

	for k, v := range responseArray {
		// Article URL
		dataUrl := v["url"].(string)

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("Processing: %s of %s, URL: %s\n", green(k+1), green(total), green(dataUrl))

		// Get article data
		data := makePetition(http.MethodGet, dataUrl, nil, nil)

		dataBytes, _ := json.Marshal(data)

		response := makePetition(http.MethodPost, url, dataBytes, tokenFlag)

		_, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// ID of the created article
		id := response["id"]

		// Get Data from article to make petitions accordingly
		metaInfoItems := []string{"tags", "references"}

		for _, meta := range metaInfoItems {
			data := data["data"].(map[string]interface{})
			attributes := data["attributes"].(map[string]interface{})
			metaInfo := attributes["meta"].(map[string]interface{})[meta].([]interface{})
			for _, item := range metaInfo {
				element := strings.TrimSuffix(meta, "s")
				element = strings.Title(element)

				urlWithId := url + "/" + id.(string) + "/create" + element

				attributes := make(map[string]interface{})

				entity := ""

				if meta == "tags" {
					name := item.(map[string]interface{})["name"].(string)
					attributes = map[string]interface{}{
						"name": name,
					}

					entity = name
				} else if meta == "references" {
					title := item.(map[string]interface{})["title"].(string)
					url := item.(map[string]interface{})["url"].(string)
					attributes = map[string]interface{}{
						"title": title,
						"url":   url,
					}

					entity = title
				}

				data := map[string]interface{}{
					"data": map[string]interface{}{
						"type":       meta,
						"attributes": attributes,
					},
				}

				dataCasted, _ := json.Marshal(data)

				fmt.Printf("Processing: %s: %s\n", element, green(entity))

				response := makePetition(http.MethodPost, urlWithId, dataCasted, tokenFlag)

				if _, err := json.Marshal(response); err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
			}
		}
	}
}
