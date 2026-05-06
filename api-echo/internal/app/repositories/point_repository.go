package repositories

import (
	"context"
	"errors"

	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IPointRepository interface {
		// ユーザーのポイントを取得（なければ作成）
		FindOrCreateByUserID(ctx context.Context, tx *gorm.DB, userID int32) (*model.UserPoints, error)
		// ポイントを加算してレベルを更新
		AddPoints(ctx context.Context, tx *gorm.DB, userID int32, points int32, newLevel int32) (*model.UserPoints, error)
		// ポイント履歴を保存
		SaveHistory(ctx context.Context, tx *gorm.DB, history *model.PointHistory) error
		// ポイント履歴をユーザーIDで取得
		FindHistoryByUserID(ctx context.Context, userID int32, limit, offset int32) ([]*model.PointHistory, error)
	}

	pointRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

func NewPointRepository(db *gorm.DB) IPointRepository {
	return &pointRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
}

// ユーザーのポイントを取得（なければ作成）
func (r *pointRepository) FindOrCreateByUserID(ctx context.Context, tx *gorm.DB, userID int32) (*model.UserPoints, error) {
	up := query.Use(r.WriterDB).UserPoints
	if tx != nil {
		up = query.Use(tx).UserPoints
	}

	// まず既存レコードを検索
	result, err := up.WithContext(ctx).
		Where(up.UserID.Eq(userID)).
		First()
	if err == nil {
		return result, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 存在しない場合は初期値で作成
	defaultLevel := int32(1)
	newRecord := &model.UserPoints{
		UserID:      userID,
		TotalPoints: 0,
		Level:       &defaultLevel,
	}
	if err := up.WithContext(ctx).Omit(field.AssociationFields).Create(newRecord); err != nil {
		return nil, err
	}
	return newRecord, nil
}

// ポイントを加算してレベルを更新
func (r *pointRepository) AddPoints(ctx context.Context, tx *gorm.DB, userID int32, points int32, newLevel int32) (*model.UserPoints, error) {
	db := r.WriterDB
	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"total_points": gorm.Expr("total_points + ?", points),
				"level":        newLevel,
			}),
		}).
		Omit("User").
		Create(&model.UserPoints{
			UserID:      userID,
			TotalPoints: points,
			Level:       &newLevel,
		}).Error
	if err != nil {
		return nil, err
	}

	return r.FindOrCreateByUserID(ctx, tx, userID)
}

// ポイント履歴を保存
func (r *pointRepository) SaveHistory(ctx context.Context, tx *gorm.DB, history *model.PointHistory) error {
	ph := query.Use(r.WriterDB).PointHistory
	if tx != nil {
		ph = query.Use(tx).PointHistory
	}

	return ph.WithContext(ctx).
		Omit(field.AssociationFields).
		Create(history)
}

// ポイント履歴をユーザーIDで取得
func (r *pointRepository) FindHistoryByUserID(ctx context.Context, userID int32, limit, offset int32) ([]*model.PointHistory, error) {
	ph := query.Use(r.ReaderDB).PointHistory

	return ph.WithContext(ctx).
		Where(ph.UserID.Eq(userID)).
		Order(ph.CreatedAt.Desc()).
		Limit(int(limit)).
		Offset(int(offset)).
		Omit(field.AssociationFields).
		Find()
}
