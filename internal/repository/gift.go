package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Gift struct {
	ID         int     `json:"id" bjson:"_id"`
	Name       string  `json:"name" bjson:"name"`
	MeliUrl    string  `json:"meli_url" bjson:"meli_url"`
	Price      float64 `json:"price" bjson:"price"`
	BuyerEmail string  `json:"buyer_email" bjson:"buyer_email"`
}

var (
	CollectionConn *mongo.Collection
	CollectionName string = "gift"
	DatabaseName   string = "mariage"
)

func SavSaveGift(ctx context.Context, client *mongo.Client) (any, error) {
	coll := client.Database(DatabaseName).Collection(CollectionName)
	gift := Gift{
		ID:   0,
		Name: "Presente teste", MeliUrl: "www.com.br",
		Price:      5.50,
		BuyerEmail: "",
	}

	result, err := coll.InsertOne(ctx, gift)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func GetGift(ctx context.Context, client *mongo.Client, giftId any) (Gift, error) {
	coll := client.Database(DatabaseName).Collection(CollectionName)
	// mongoID, err := primitive.ObjectIDFromHex("0")
	// if err != nil {
	// 	return Gift{}, err
	// }

	// filter := bson.D{{"id", 0}}

	var gift Gift
	err := coll.FindOne(ctx, bson.M{"_id": giftId}).Decode(&gift)
	if err != nil {
		return Gift{}, err
	}
	return gift, nil
}
