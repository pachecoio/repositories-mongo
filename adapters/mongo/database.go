package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type DatabaseClient struct {
	uri     string
	Client  *mongo.Client
	Context context.Context
}

func NewDatabaseClient(ctx ...context.Context) *DatabaseClient {
	uri := os.Getenv("MONGODB_URI")
	_ctx := context.Background()
	if len(ctx) > 0 {
		_ctx = ctx[0]
	}
	client, err := mongo.Connect(_ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return &DatabaseClient{
		uri,
		client,
		_ctx,
	}
}

func (d *DatabaseClient) GetCollection(databaseName string, collectionName string) *mongo.Collection {
	return d.Client.Database(databaseName).Collection(collectionName)
}

func (d *DatabaseClient) Disconnect() {
	d.Client.Disconnect(d.Context)
}
