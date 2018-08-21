package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	CATALOG_ITEMS = "data-catalog-items-categories-request.json"
)

// Control flow variables
var (
	token                 *string
	urlCatalogsPrefix     = "https://"
	urlCatalogsSuffix     = ".api.culturacolectiva.com/catalogs"
	urlCatalogsItemSuffix = "/item"
	urlCatalogs           = ""
	urlCatalogsItem       = ""
	catalogs              CatalogsResponseGet
	catalog               Catalog
	catalogId             string
	catalogItems          []CatalogItemRequest
)

func main() {
	token = flag.String("token", "", "Token needed for make the petition")
	environment := flag.String("environment", "dev", "Environment to make the petition {dev, staging}")
	flag.Parse()

	if *token == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	urlCatalogs = urlCatalogsPrefix + *environment + urlCatalogsSuffix
	urlCatalogsItem = urlCatalogs + urlCatalogsItemSuffix

	getCatalogs()
	getCategoryIdCatalog()
	createCatalogItems()

	for _, v := range catalogItems {
		makeAPICalls(v)
	}
}

func makeAPICalls(item CatalogItemRequest) {
	client := &http.Client{}

	itemByte, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, urlCatalogsItem, bytes.NewBuffer(itemByte))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", *token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)

	fmt.Printf("%s", response)
}

func getCategoryIdCatalog() {
	if len(catalogs.Data) == 0 {
		createCategoryCatalog()
		catalogId = catalog.Id
	}

	if len(catalogs.Data) > 0 {
		for _, v := range catalogs.Data {
			if v.Attributes.Type == "categories" {
				catalogId = v.Id
			}
		}
	}
}
