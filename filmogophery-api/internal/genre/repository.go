package genre

import (
	"context"
	"filmogophery/internal/db"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"

	"gorm.io/gorm"
)

type (
	IQueryRepository interface {
		FindByName(ctx context.Context, name []string) ([]*model.Genre, error)
	}

	GenreRepository struct {
		DB *gorm.DB
	}
)

func NewQueryRepository() *IQueryRepository {
	var queryRepo IQueryRepository = &GenreRepository{
		DB: db.READER_DB,
	}
	return &queryRepo
}

func (r *GenreRepository) FindByName(ctx context.Context, name []string) ([]*model.Genre, error) {
	g := query.Use(r.DB).Genre

	genres, err := g.WithContext(ctx).Where(g.Name.In(name...)).Find()
	if err != nil {
		return nil, err
	}

	return genres, nil
}
