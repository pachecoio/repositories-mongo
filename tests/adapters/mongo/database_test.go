package mongo

import (
	"testing"

	testUtils "github.com/pachecoio/repositories-mongo/utils"
)

func TestDatabase_Connection(t *testing.T) {
	db, teardown := testUtils.GetTestDB(t)
	defer teardown()

	if db == nil {
		t.Error("Expected a database instance, got nil")
	}
}

func TestDatabase_GetCollection(t *testing.T) {
	db, teardown := testUtils.GetTestDB(t)
	defer teardown()

	collection := db.GetCollection("test", "app")
	if collection == nil {
		t.Error("Expected a collection instance, got nil")
	}
}

func TestDatabase_InsertSample(t *testing.T) {
	db, teardown := testUtils.GetTestDB(t)
	defer teardown()

	collection := db.GetCollection("test", "app")
	if collection == nil {
		t.Error("Expected a collection instance, got nil")
	}

	type Sample struct {
		Name string `json:"name"`
	}
	sample := Sample{
		Name: "Jon Snow",
	}

	res, err := collection.InsertOne(db.Context, sample)

	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	if res == nil {
		t.Error("Expected a result instance, got nil")
	}

	if res.InsertedID == nil {
		t.Error("Expected an inserted ID, got nil")
	}
}
