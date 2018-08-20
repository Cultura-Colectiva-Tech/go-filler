package main

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
