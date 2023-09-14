package mongo

import (
	"testing"

	"github.com/pachecoio/repositories-mongo/adapters/mongo"
	testUtils "github.com/pachecoio/repositories-mongo/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

type Sample struct {
	Name string `json:"name"`
}

type SampleRepository struct {
	mongo.Repository[Sample]
}

func NewSampleRepository(db *mongo.DatabaseClient, databaseName string) *SampleRepository {
	return &SampleRepository{
		Repository: *mongo.NewRepository[Sample](db, databaseName),
	}
}

type CustomFilters struct {
	Name string `bson:"name"`
}

func (f *CustomFilters) ToQuery() any {
	return bson.D{{"name", f.Name}}
}

func TestRepository_Create(t *testing.T) {
	db, teardown := testUtils.GetTestDB(t)
	defer teardown()

	repo := NewSampleRepository(db, "test")

	sample := &Sample{
		Name: "Jon Snow",
	}
	insertedId, err := repo.Create(sample)
	assert.Nil(t, err)
	assert.NotEmpty(t, insertedId)
}

func TestRepository_Filter(t *testing.T) {
	db, teardown := testUtils.GetTestDB(t)
	defer teardown()

	repo := NewSampleRepository(db, "test")

	sample := &Sample{
		Name: "Jon Snow",
	}
	insertedId, err := repo.Create(sample)
	assert.Nil(t, err)
	assert.NotEmpty(t, insertedId)

	filters := &mongo.DefaultFilters{}
	samples, err := repo.Filter(filters)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(samples))

	customFilters := &CustomFilters{
		Name: "Jon Snow",
	}

	samples, err = repo.Filter(customFilters)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(samples))

	customFilters = &CustomFilters{
		Name: "Arya Stark",
	}
	samples, err = repo.Filter(customFilters)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(samples))
}

func TestRepository_Count(t *testing.T) {
	db, teardown := testUtils.GetTestDB(t)
	defer teardown()

	repo := NewSampleRepository(db, "test")

	sample := &Sample{
		Name: "Jon Snow",
	}
	insertedId, err := repo.Create(sample)
	assert.Nil(t, err)
	assert.NotEmpty(t, insertedId)

	filters := &mongo.DefaultFilters{}
	count, err := repo.Count(filters)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	customFilters := &CustomFilters{
		Name: "Jon Snow",
	}
	count, err = repo.Count(customFilters)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	customFilters = &CustomFilters{
		Name: "Arya Stark",
	}
	count, err = repo.Count(customFilters)
	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}

type SampleUpdater struct {
	Name string `json:"name"`
}

func (u *SampleUpdater) ToUpdate() any {
	return bson.M{
		"$set": bson.M{
			"name": u.Name,
		},
	}
}

func TestRepository_Update(t *testing.T) {
	db, teardown := testUtils.GetTestDB(t)
	defer teardown()

	repo := NewSampleRepository(db, "test")

	sample := &Sample{
		Name: "Jon Snow",
	}
	insertedId, err := repo.Create(sample)
	assert.Nil(t, err)
	assert.NotEmpty(t, insertedId)

	data := &SampleUpdater{
		Name: "Arya Stark",
	}
	err = repo.Update(insertedId, data)
	assert.Nil(t, err)

	filters := &CustomFilters{
		Name: "Arya Stark",
	}
	samples, err := repo.Filter(filters)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(samples))
}

func TestRepository_Delete(t *testing.T) {
	db, teardown := testUtils.GetTestDB(t)
	defer teardown()

	repo := NewSampleRepository(db, "test")

	sample := &Sample{
		Name: "Jon Snow",
	}
	insertedId, err := repo.Create(sample)
	assert.Nil(t, err)
	assert.NotEmpty(t, insertedId)

	err = repo.Delete(insertedId)
	assert.Nil(t, err)

	filters := &CustomFilters{
		Name: "Jon Snow",
	}
	samples, err := repo.Filter(filters)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(samples))
}
