package testUtils

import (
	"context"
	"testing"

	mim "github.com/pachecoio/inmemory-mongo"
	"github.com/pachecoio/repositories-mongo/adapters/mongo"
)

func GetInMemoryServer(ctx context.Context) *mim.Server {
	server, err := mim.Start(ctx, "6.0.5")
	if err != nil {
		panic(err)
	}
	return server
}

// GetTestDB returns a test database instance and a teardown function
// to be called after the test is done.
// This DB is connected to a local MongoDB instance.
func GetTestDB(t *testing.T) (*mongo.DatabaseClient, func()) {
	ctx := context.Background()
	mongoServer := GetInMemoryServer(ctx)
	t.Setenv("MONGODB_URI", mongoServer.URI())
	db := mongo.NewDatabaseClient(ctx)
	if db == nil {
		t.Error("Expected a database instance, got nil")
	}
	return db, func() {
		db.Disconnect()
	}
}
