package repositories

import (
	"context"
	"time"

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

		// 有効なトークンを取得
		FindActiveTokenByUserID(ctx context.Context, user *model.Users, now time.Time) ([]*model.RefreshTokens, error)

		// --- Update --- //

		// トークンを無効化
		Revoke(ctx context.Context, tx *gorm.DB, tokenIDs []int32, now time.Time) error

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
	rt := query.Use(r.WriterDB).RefreshTokens
	if tx != nil {
		rt = query.Use(tx).RefreshTokens
	}

	return rt.WithContext(ctx).
		Omit(field.AssociationFields).
		Create(token)
}

// 有効なトークンを取得
func (r *tokenRepository) FindActiveTokenByUserID(
	ctx context.Context, user *model.Users, now time.Time,
) ([]*model.RefreshTokens, error) {
	rt := query.Use(r.ReaderDB).RefreshTokens

	return rt.WithContext(ctx).
		Where(
			rt.UserID.Eq(user.ID),
			rt.RevokedAt.IsNull(),
			rt.ExpiresAt.Gte(now),
		).
		Find()
}

// トークンを無効化
func (r *tokenRepository) Revoke(
	ctx context.Context, tx *gorm.DB, tokenIDs []int32, now time.Time,
) error {
	rt := query.Use(r.WriterDB).RefreshTokens
	if tx != nil {
		rt = query.Use(tx).RefreshTokens
	}

	_, err := rt.WithContext(ctx).
		Where(rt.ID.In(tokenIDs...)).
		UpdateSimple(rt.RevokedAt.Value(now))

	return err
}
