package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	}
}
