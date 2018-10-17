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
	dataFlag        *string
	tokenFlag       *string
	environmentFlag *string
	nameFlag        *string
	descriptionFlag *string
	typeContentFlag *string
	monthFlag       *string
	yearFlag        *string
	typePostFlag    *string
	initIndexFlag   *int
)

func main() {
	v := flag.Bool("v", false, "Print the version of the program")
	version := flag.Bool("version", false, "Print the version of the program")

	dataFlag = flag.String("data", "https://cucodev.culturacolectiva.com/jsoncategory/", "URL for get the data (json) to add")
	tokenFlag = flag.String("token", "", "Token needed for make the petition")
	environmentFlag = flag.String("environment", "dev", "Environment to make the petition {dev, staging}")
	nameFlag = flag.String("name", "categories", "Name of catalog")
	descriptionFlag = flag.String("description", "Description of Catalog", "Description of catalog")
	typeContentFlag = flag.String("type", "categories", "Type of Catalog's content to create")
	monthFlag = flag.String("month", "01", "Month to bring Articles. Default: 01")
	yearFlag = flag.String("year", "2018", "Year to bring Articles. Default: 2018")
	typePostFlag = flag.String("type-post", "video", "Article type to be searched. Default: video")
	initIndexFlag = flag.Int("init-index", 1, "Index to start search Articles")

	flag.Parse()

	if *v || *version {
		fmt.Printf("go-filler version %s\n", appVersion)
		os.Exit(0)
	}

	if *tokenFlag == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *typeContentFlag == "categories" || *typeContentFlag == "tags" {
		catalogsLogic()
	}

	if *typeContentFlag == "articles" {
		articlesLogic()
	}
}
