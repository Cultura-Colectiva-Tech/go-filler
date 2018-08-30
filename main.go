package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	urlPrefix             = "https://"
	urlSuffix             = ".api.culturacolectiva.com/"
	urlCatalogsSuffix     = "catalogs"
	urlCatalogsItemSuffix = "/item"
)

var (
	dataFlag    *string
	token       *string
	environment *string
	name        *string
	description *string
	typeContent *string
)

func main() {
	dataFlag = flag.String("data", "https://cucodev.culturacolectiva.com/jsoncategory/", "URL for get the data (json) to add")
	token = flag.String("token", "", "Token needed for make the petition")
	environment = flag.String("environment", "dev", "Environment to make the petition {dev, staging}")
	v := flag.Bool("v", false, "Print the version of the program")
	version := flag.Bool("version", false, "Print the version of the program")
	name = flag.String("name", "categories", "Name of catalog")
	description = flag.String("description", "Description of Catalog", "Description of catalog")
	typeContent = flag.String("type", "categories", "Type of Catalog's content to create")
	flag.Parse()

	if *v || *version {
		fmt.Printf("go-filler version %s\n", appVersion)
		os.Exit(0)
	}

	if *token == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *typeContent == "categories" || *typeContent == "tags" {
		catalogsLogic()
	}

	if *typeContent == "articles" {
		articlesLogic()
	}
}
