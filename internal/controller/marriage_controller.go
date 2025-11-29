package controller

import (
	"context"
	"log"

	"github.com/js-bruno/mariage-api/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InsertGift(ctx context.Context, client *mongo.Client) (any, error) {
	coll := client.Database(repository.DatabaseName).Collection(repository.CollectionName)
	gift := repository.Gift{
		ID:         0,
		Name:       "Presente teste",
		MeliUrl:    "www.com.br",
		Price:      5.50,
		BuyerEmail: "",
	}

	result, err := coll.InsertOne(ctx, gift)
	if err != nil {
		return nil, err
	}
	// oi := result.InsertedID
	log.Printf("%T", result.InsertedID)
	return result.InsertedID, nil
}
