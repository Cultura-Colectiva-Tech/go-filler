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

// Request Catalog
type DataRequest struct {
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type CatalogRequest struct {
	Data DataRequest `json:"data"`
}

// Request CatalogItem
type CatalogItemAttributesRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Slug        string   `json:"slug"`
	Scope       []string `json:"scope"`
	Parent      string   `json:"parent"`
}

type CatalogItemDataRequest struct {
	Type       string                       `json:"type"`
	Attributes CatalogItemAttributesRequest `json:"attributes"`
}

type CatalogItemRequest struct {
	Data CatalogItemDataRequest `json:"data"`
}

// Response Catalogs Data
type Attributes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CreatedAt   string
	UpdatedAt   string
}

type Catalog struct {
	Type       string
	Id         string
	Attributes Attributes
}

// Response Catalogs Metadata
type Paginate struct {
	PageCount int
	Page      int
	Limit     int
}

type Metadata struct {
	Paginate Paginate
}

// Response Catalogs
type CatalogsResponseGet struct {
	Metadata Metadata
	Data     []Catalog
}

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
