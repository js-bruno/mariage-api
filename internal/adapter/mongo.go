package adapter

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func MongoOpenConnetion() (client *mongo.Client, close func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	mongoURL := fmt.Sprintf("mongodb://admin:password@localhost:%s/mariage?authSource=admin", "27017")

	opts := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	close = func() {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}

	defer cancel()
	return client, close
}
