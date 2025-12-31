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
	ITokenRepository interface {
		// --- Create --- //

		// トークンを作成
		Save(ctx context.Context, tx *gorm.DB, user *model.RefreshTokens) error

		// --- Read --- //

		// --- Update --- //

		// --- Delete --- //
	}
	tokenRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

func NewTokenRepository(db *gorm.DB) ITokenRepository {
	return &tokenRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
}

// トークンを作成
func (r *tokenRepository) Save(ctx context.Context, tx *gorm.DB, token *model.RefreshTokens) error {
	rv := query.Use(r.WriterDB).RefreshTokens
	if tx != nil {
		rv = query.Use(tx).RefreshTokens
	}

	return rv.WithContext(ctx).
		Omit(field.AssociationFields).
		Create(token)
}
