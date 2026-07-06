package main

import (
	"Domovoy/realt"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := realt.NewClient()

	filter := realt.SearchFilter{
		Category:  5,
		Rooms:     []string{"2"},
		PriceTo:   "90000",
		PriceType: "840",
		TownUUID:  "4cb07174-7b00-11eb-8943-0cc47adabd66",
	}

	res, err := client.Search(ctx, filter, 1)
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}

	fmt.Printf("total: %d\n", res.Body.Pagination.TotalCount)
	for _, obj := range res.Body.Results {
		fmt.Printf("%+v\n", obj)
	}
}
