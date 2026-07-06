package main

import (
	"Domovoy/internal/storage"
	"Domovoy/realt"
	"context"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://realty:realty@localhost:5432/realty?sslmode=disable"
		log.Print("DATABASE_URL environment variable not set")
	}

	store, err := storage.New(ctx, dsn)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer store.Close()

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

	newCount := 0
	for _, obj := range res.Body.Results {
		isNew, err := store.Upsert(ctx, obj)
		if err != nil {
			log.Printf("insert failed: %v", err)
			continue
		}
		if isNew {
			newCount++
			log.Printf("NEW: %s | %.0f | %d rooms | %s", obj.UUID, obj.Price, obj.Rooms, obj.Address)
		}
	}

	log.Printf("done: %d new listings", newCount)
}
