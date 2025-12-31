package repositories

import (
	"context"

	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IUserRepository interface {
		// --- Create --- //

		// ユーザーを作成
		Save(ctx context.Context, tx *gorm.DB, user *model.Users) error

		// --- Read --- //

		// --- Update --- //

		// --- Delete --- //
	}
	userRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
}

// ユーザーを作成
func (r *userRepository) Save(ctx context.Context, tx *gorm.DB, user *model.Users) error {
	rv := query.Use(r.WriterDB).Users
	if tx != nil {
		rv = query.Use(tx).Users
	}

	return rv.WithContext(ctx).
		Omit(field.AssociationFields).
		Create(user)
}
