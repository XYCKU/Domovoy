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

	data, err := client.Search(ctx, 1)
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}

	fmt.Println(string(data))
}
