package mongo

import (
	"context"
	"fmt"

	"github.com/pachecoio/repositories-mongo/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoLib "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DefaultFilters struct{}

func (f *DefaultFilters) ToQuery() any {
	return bson.D{}
}

type Repository[T base.Model] struct {
	db         *DatabaseClient
	Collection *mongoLib.Collection
	Context    context.Context
}

func NewRepository[T base.Model](db *DatabaseClient, databaseName string) *Repository[T] {
	modelName := fmt.Sprintf("%T", *new(T))
	return &Repository[T]{
		db:         db,
		Collection: db.GetCollection(databaseName, modelName),
		Context:    context.Background(),
	}
}

func (r *Repository[T]) Create(model *T) (string, error) {
	sess, err := r.Collection.Database().Client().StartSession()
	if err != nil {
		return "", err
	}
	defer sess.EndSession(r.Context)

	createdId := ""
	sess.WithTransaction(r.Context, func(sessCtx mongoLib.SessionContext) (interface{}, error) {
		res, err := r.Collection.InsertOne(r.Context, model)
		if err != nil {
			return "", err
		}
		createdId = res.InsertedID.(primitive.ObjectID).Hex()
		return createdId, nil
	})

	return createdId, err

}

func (r *Repository[T]) Filter(filters base.Filters, opts ...base.FilterOptions) ([]T, error) {
	items := make([]T, 0)

	if filters == nil {
		filters = &DefaultFilters{}
	}

	mongoOptions := options.Find()
	if len(opts) > 0 {
		if opts[0].Limit > 0 {
			mongoOptions.SetLimit(int64(opts[0].Limit))
		}
		if opts[0].Offset > 0 {
			mongoOptions.SetSkip(int64(opts[0].Offset))
		}
		if opts[0].Sort != nil {
			mongoOptions.SetSort(opts[0].Sort)
		}
	}

	cursor, err := r.Collection.Find(r.Context, filters.ToQuery(), mongoOptions)
	if err != nil {
		return nil, err
	}
	for cursor.Next(r.Context) {
		item := *new(T)
		err := cursor.Decode(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository[T]) Count(filters base.Filters) (int, error) {
	if filters == nil {
		filters = &DefaultFilters{}
	}
	count, err := r.Collection.CountDocuments(r.Context, filters.ToQuery())
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *Repository[T]) FindOne(filters base.Filters) (*T, error) {
	if filters == nil {
		filters = &DefaultFilters{}
	}
	item := new(T)
	err := r.Collection.FindOne(r.Context, filters.ToQuery()).Decode(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *Repository[T]) Get(id any) (*T, error) {
	item := new(T)
	_id, err := primitive.ObjectIDFromHex(id.(string))
	err = r.Collection.FindOne(r.Context, map[string]primitive.ObjectID{"_id": _id}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *Repository[T]) Update(id any, data base.PartialUpdate[T]) error {
	_id, err := primitive.ObjectIDFromHex(id.(string))
	filterById := bson.D{{
		"_id", _id,
	}}
	dataToUpdate := data.ToUpdate()

	sess, err := r.Collection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(r.Context)

	sess.WithTransaction(r.Context, func(sessCtx mongoLib.SessionContext) (interface{}, error) {
		_, err = r.Collection.UpdateOne(r.Context, filterById, dataToUpdate)
		return nil, err
	})

	return err
}

func (r *Repository[T]) Delete(id any) error {
	_id, err := primitive.ObjectIDFromHex(id.(string))
	if err != nil {
		return err
	}

	sess, err := r.Collection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(r.Context)

	sess.WithTransaction(r.Context, func(sessCtx mongoLib.SessionContext) (interface{}, error) {
		_, err = r.Collection.DeleteOne(r.Context, map[string]primitive.ObjectID{"_id": _id})
		return nil, err
	})
	return err
}
