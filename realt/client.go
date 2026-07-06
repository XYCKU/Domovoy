package realt

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
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

type graphqlResponse struct {
	Data json.RawMessage `json:"data"`
}

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{http: &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) Search(ctx context.Context, page int) (json.RawMessage, error) {
	reqBody := []graphqlRequest{
		{
			OperationName: "searchObjects",
			Query:         searchQuery,
			Variables: map[string]interface{}{
				"data": map[string]interface{}{
					"where": map[string]interface{}{
						"category":       5,
						"rooms":          []string{"2"},
						"priceTo":        "90000",
						"priceType":      "840",
						"priceMeterType": "all",
						"addressV2": []map[string]string{
							{
								"townUuid": "4cb07174-7b00-11eb-8943-0cc47adabd66",
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

	fmt.Println(string(raw))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Realt-Client", "www@6.14.4")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response: %d\n", resp.StatusCode)
		}
	}(resp.Body)

	b, _ := io.ReadAll(resp.Body)

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Println(string(b))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var out []graphqlResponse
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("empty graphql response")
	}

	return out[0].Data, nil
}
