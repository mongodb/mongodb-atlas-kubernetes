package contract

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func WithDatabase(prefix string) OptResourceFunc {
	return func(ctx context.Context, resources *TestResources) (*TestResources, error) {
		uri := resources.URI()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MongoDB at %s: %w", uri, err)
		}
		db := client.Database(prefix)
		resources.DatabaseName = prefix
		collectionName := newRandomName(fmt.Sprintf("%s-collection", prefix))
		collection := db.Collection(collectionName)
		resources.CollectionName = collectionName
		_, err = collection.InsertOne(ctx, resources)
		if err != nil {
			return nil, fmt.Errorf("failed to insert test data into MongoDB %s at %s.%s: %w",
				uri, prefix, collectionName, err)
		}
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				log.Printf("Failed to disconnect from MongoDB at %s: %v", uri, err)
			}
		}()
		return resources, nil
	}
}
