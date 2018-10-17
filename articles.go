package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func articlesLogic() {
	url := urlPrefix + *environmentFlag + urlSuffix + *typeContentFlag

	body := "month=" + *monthFlag + "&year=" + *yearFlag + "&type=" + *typePostFlag

	initIndex := *initIndexFlag

	// Get all the articles
	responseArray := makePetitionResponseArray(http.MethodPost, *dataFlag, []byte(body), nil)

	total := len(responseArray)

	if initIndex != 1 {
		if initIndex <= 1 || initIndex > total {
			fmt.Println("Index range not valid")
			os.Exit(0)
		}
	}

	initIndex = (*initIndexFlag - 1)

	if total < 1 {
		fmt.Println("There is no data to be processed")
		os.Exit(0)
	}

	for k, v := range responseArray[initIndex:] {

		// Article URL
		dataUrl := v["url"].(string)

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("Processing: %s of %s, URL: %s\n", green(k+*initIndexFlag), green(total), green(dataUrl))

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
		supplementInfo := []string{"author", "facebook", "editor", "seo", "ads", "meta"}

		// Info
		infoItems := map[string][]string{
			"meta": []string{"tags", "references"},
		}

		for _, articleInfo := range supplementInfo {
			data := data["data"].(map[string]interface{})
			attributes := data["attributes"].(map[string]interface{})

			if articleInfo != "meta" {
				info := attributes[articleInfo].(map[string]interface{})

				// Custom logic for the ads
				if articleInfo == "ads" {
					for key, value := range info {
						// We exclude the already string values
						if key != "adUnit" {
							info[key] = strconv.FormatBool(value.(bool))
						}
					}
				}

				if articleInfo == "seo" {
					for _, keywordName := range info["keywords"].([]interface{}) {

						data := map[string]interface{}{
							"data": map[string]interface{}{
								"type": "keywords",
								"attributes": map[string]interface{}{
									"name": keywordName,
								},
							},
						}

						dataCasted, _ := json.Marshal(data)

						urlWithId := url + "/" + id.(string) + "/createKeyword"

						makePetition(http.MethodPost, urlWithId, dataCasted, tokenFlag)
					}
				}

				data := map[string]interface{}{
					"data": map[string]interface{}{
						"type": "articles",
						"attributes": map[string]interface{}{
							articleInfo: info,
						},
					},
				}

				fmt.Printf("Processing: %s\n", green(articleInfo))

				dataCasted, _ := json.Marshal(data)

				urlWithId := url + "/" + id.(string)

				response := makePetition(http.MethodPatch, urlWithId, dataCasted, tokenFlag)

				if _, err := json.Marshal(response); err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
			}

			for _, meta := range infoItems[articleInfo] {
				info := attributes[articleInfo].(map[string]interface{})[meta].([]interface{})

				for _, item := range info {
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
}
