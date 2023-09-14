package base

type Repository[T Model] interface {
	Filter(filters Filters, options ...FilterOptions) ([]T, error)
	FindOne(filters Filters) (*T, error)
	Count(filters Filters) (int, error)
	Get(id any) (*T, error)
	Create(model *T) (any, error)
	Update(id any, data PartialUpdate[T]) error
	Delete(id any) error
}

type PartialUpdate[T Model] interface {
	ToUpdate() any
}

type Filters interface {
	ToQuery() any
}

type FilterOptions struct {
	Limit  int
	Offset int
	Sort   any
}
