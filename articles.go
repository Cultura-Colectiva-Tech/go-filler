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

func fillArticleFromJSON(jsons []string) {
	var data [](map[string]interface{})
	for _, v := range jsons {
		j := map[string]interface{}{
			"url": v,
		}

		data = append(data, j)
	}

	if *initIndexFlag > 0 {
		fmt.Println("Don't support index")
		os.Exit(0)
	}

	process(data, *initIndexFlag)
	os.Exit(0)
}

func articlesLogic() {

	body := "month=" + *monthFlag + "&year=" + *yearFlag + "&type=" + *typePostFlag

	initIndex := *initIndexFlag

	// Get all the articles
	responseArray, err := makePetitionResponseArray(http.MethodPost, *dataFlag, []byte(body), nil)

	if err != nil {
		formatError("Can't process request. The error was", err)
		os.Exit(0)
	}

	total := len(responseArray)

	if initIndex != 0 {
		if initIndex > total {
			fmt.Println("Index range not valid")
			os.Exit(0)
		}

		initIndex = (*initIndexFlag - 1)
	}

	if total < 1 {
		fmt.Println("There is no data to be processed")
		os.Exit(0)
	}

	process(responseArray, initIndex)
}

/**
 * Generic function to process data
 */
func process(data []map[string]interface{}, index int) {
	url := urlPrefix + *environmentFlag + urlSuffix + *typeContentFlag
	total := len(data)
	init := index

	for k, v := range data[init:] {
		// Article URL
		dataURL := v["url"].(string)

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("Processing: %s of %s, URL: %s\n", green(k+(init+1)), green(total), green(dataURL))

		// Get article data
		data, dataError := makePetition(http.MethodGet, dataURL, nil, nil)
		if dataError != nil {
			log.Print(dataError)
			continue
		}

		dataBytes, _ := json.Marshal(data)

		response, artError := makePetition(http.MethodPost, url, dataBytes, tokenFlag)
		if artError != nil {
			log.Print(artError)
			continue
		}

		_, marshalError := json.Marshal(response)
		if marshalError != nil {
			log.Print(marshalError)
			continue
		}

		// ID of the created article
		id := response["id"]

		// Get Data from article to make petitions accordingly
		supplementInfo := []string{"author", "facebook", "editor", "seo", "ads", "meta"}

		// Info
		infoItems := map[string][]string{
			"meta": []string{"tags", "references", "properties"},
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

						urlWithID := url + "/" + id.(string) + "/createKeyword"

						makePetition(http.MethodPost, urlWithID, dataCasted, tokenFlag)
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

				urlWithID := url + "/" + id.(string)

				response, error := makePetition(http.MethodPatch, urlWithID, dataCasted, tokenFlag)

				if _, err := json.Marshal(response); error != nil {
					log.Println(err)
					continue
				}
			}

			for _, meta := range infoItems[articleInfo] {
				info := attributes[articleInfo].(map[string]interface{})[meta].([]interface{})

				for _, item := range info {
					element := strings.TrimSuffix(meta, "s")
					element = strings.Title(element)

					if element == "Propertie" {
						element = "Property"
					}

					urlWithID := url + "/" + id.(string) + "/create" + element

					attributes := make(map[string]interface{})

					entity := ""

					if meta == "tags" {
						name := item.(map[string]interface{})["name"].(string)
						slug := item.(map[string]interface{})["slug"].(string)
						attributes = map[string]interface{}{
							"name": name,
							"slug": slug,
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
					} else if meta == "properties" {
						coverImage := item.(map[string]interface{})["coverImage"].(string)
						attributes = map[string]interface{}{
							"name":  "coverImage",
							"value": coverImage,
						}

						entity = coverImage
					}

					data := map[string]interface{}{
						"data": map[string]interface{}{
							"type":       meta,
							"attributes": attributes,
						},
					}

					dataCasted, _ := json.Marshal(data)

					fmt.Printf("Processing: %s: %s\n", element, green(entity))

					response, metaError := makePetition(http.MethodPost, urlWithID, dataCasted, tokenFlag)

					if _, err := json.Marshal(response); metaError != nil {
						log.Println(err)
						continue
					}
				}
			}
		}
	}
}
