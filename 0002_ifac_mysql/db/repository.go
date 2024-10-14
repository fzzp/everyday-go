package db

type Repository struct {
	Products ProductRepo
	Category CategoryRepo
}

func NewRepository(db Queryable) *Repository {
	return &Repository{
		Products: NewProductRepo(db),
		Category: NewCategoryRepo(db),
	}
}
