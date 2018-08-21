package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func createCatalogItems() {
	jsonFile, err := os.Open(CATALOG_ITEMS)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(jsonBytes, &catalogItems)

	for k, _ := range catalogItems {
		catalogItems[k].Data.Attributes.Parent = catalogId
	}
}
