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
	url := urlPrefix + *environment + urlSuffix + *typeContent

	data := makePetition(http.MethodGet, *dataFlag, nil, nil)

	dataBytes, _ := json.Marshal(data)

	response := makePetition(http.MethodPost, url, dataBytes, token)

	_, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	green := color.New(color.FgGreen).SprintFunc()

	fmt.Printf("Processing: %s\n", green(*dataFlag))
}
