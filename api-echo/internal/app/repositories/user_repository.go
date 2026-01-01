package repositories

import (
	"context"
	"errors"

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

		// ID に一致するアクティブなユーザーを取得
		FindByID(ctx context.Context, id int32) (*model.Users, error)
		// ユーザーを取得
		FindByEmail(ctx context.Context, email string) (*model.Users, error)

		// --- Update --- //

		// ユーザーを更新
		Update(ctx context.Context, tx *gorm.DB, user *model.Users) error

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
	u := query.Use(r.WriterDB).Users
	if tx != nil {
		u = query.Use(tx).Users
	}

	return u.WithContext(ctx).
		Omit(field.AssociationFields).
		Create(user)
}

// ID に一致するアクティブなユーザーを取得
func (r *userRepository) FindByID(ctx context.Context, id int32) (*model.Users, error) {
	u := query.Use(r.ReaderDB).Users

	result, err := u.WithContext(ctx).
		Where(
			u.ID.Eq(id),
			u.IsActive.Is(true),
		).Take()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return result, nil
}

// ユーザーを取得
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.Users, error) {
	u := query.Use(r.ReaderDB).Users

	result, err := u.WithContext(ctx).
		Where(u.Email.Eq(email)).
		Take()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return result, nil
}

// ユーザーを更新
func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user *model.Users) error {
	u := query.Use(r.WriterDB).Users
	if tx != nil {
		u = query.Use(tx).Users
	}

	_, err := u.WithContext(ctx).
		Where(u.ID.Eq(user.ID)).
		Updates(user)

	return err
}
