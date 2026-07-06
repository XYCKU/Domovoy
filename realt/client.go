package realt

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const endpoint = "https://realt.by/bff/graphql"

//go:embed queries/search.graphql
var searchQuery string

type graphqlRequest struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type SearchFilter struct {
	Category  int
	Rooms     []string
	PriceTo   string
	PriceType string
	TownUUID  string
}

type SearchResult struct {
	Success bool `json:"success"`
	Body    struct {
		Pagination Pagination `json:"pagination"`
		Results    []Object   `json:"results"`
	} `json:"body"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
}

type Object struct {
	UUID          string  `json:"uuid"`
	Code          int     `json:"code"`
	Category      int     `json:"category"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
	RaiseDate     string  `json:"raiseDate"`
	Price         float64 `json:"price"`
	PriceCurrency int     `json:"priceCurrency"`
	PricePerM2    float64 `json:"pricePerM2"`
	Rooms         int     `json:"rooms"`
	AreaTotal     float64 `json:"areaTotal"`
	AreaLiving    float64 `json:"areaLiving"`
	AreaKitchen   float64 `json:"areaKitchen"`
	Storey        int     `json:"storey"`
	Storeys       int     `json:"storeys"`
	BuildingYear  int     `json:"buildingYear"`

	TownName     string `json:"townName"`
	DistrictName string `json:"stateDistrictName"`
	StreetName   string `json:"streetName"`
	HouseNumber  int    `json:"houseNumber"`
	Address      string `json:"address"`

	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

type Client struct {
	http  *http.Client
	Debug bool
}

func NewClient() *Client {
	return &Client{http: &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) Search(ctx context.Context, filter SearchFilter, page int) (*SearchResult, error) {
	reqBody := []graphqlRequest{
		{
			OperationName: "searchObjects",
			Query:         searchQuery,
			Variables: map[string]interface{}{
				"data": map[string]interface{}{
					"where": map[string]interface{}{
						"category":       filter.Category,
						"rooms":          filter.Rooms,
						"priceTo":        filter.PriceTo,
						"priceType":      filter.PriceType,
						"priceMeterType": "all",
						"addressV2": []map[string]string{
							{
								"townUuid": filter.TownUUID,
							},
						},
					},
					"pagination": map[string]interface{}{
						"page":     page,
						"pageSize": 30,
					},
					"sort": []map[string]string{
						{"by": "paymentStatus", "order": "DESC"},
						{"by": "priority", "order": "DESC"},
						{"by": "raiseDate", "order": "DESC"},
						{"by": "updatedAt", "order": "DESC"},
					},
					"extraFields":       nil,
					"isReactAdaptiveUA": false,
				},
			},
		},
	}

	raw, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		fmt.Println(string(raw))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error closing body: %v", err)
		}
	}()

	b, _ := io.ReadAll(resp.Body)

	if c.Debug {
		fmt.Printf("Status: %d: %s\n", resp.StatusCode, string(b))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var out []struct {
		Data struct {
			SearchObjects SearchResult `json:"searchObjects"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("empty graphql response")
	}

	return &out[0].Data.SearchObjects, nil
}
