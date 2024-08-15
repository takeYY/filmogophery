package movie

import (
	"context"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"filmogophery/pkg/gen/model"
	"filmogophery/pkg/gen/query"
)

type (
	IQueryRepository interface {
		FindByID(ctx context.Context, id *int32) (*model.Movie, error)
		Find(ctx context.Context) ([]*model.Movie, error)
	}
	ICommandRepository interface {
		Save(movie *model.Movie) (*model.Movie, error)

		GetMediaIdByCode(ctx context.Context, code *string) (*int32, error)
		UpdateImpression(ctx context.Context, movieImpression *model.MovieImpression) (*gen.ResultInfo, error)
		SaveRecord(ctx context.Context, watchRecord *model.MovieWatchRecord) (*model.MovieWatchRecord, error)
	}

	MovieRepository struct {
		DB *gorm.DB
	}
)

func (r *MovieRepository) FindByID(ctx context.Context, id *int32) (*model.Movie, error) {
	m := query.Use(r.DB).Movie

	movie, err := m.WithContext(ctx).
		Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).
		Where(m.ID.Eq(*id)).
		First()
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (r *MovieRepository) Find(ctx context.Context) ([]*model.Movie, error) {
	m := query.Use(r.DB).Movie

	movies, err := m.WithContext(ctx).
		Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).
		Find()
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) Save(movie *model.Movie) (*model.Movie, error) {
	var err error
	defer func() {
		if err != nil {
			r.DB.Rollback()
		} else {
			r.DB.Commit()
		}
	}()

	m := query.Use(r.DB).Movie
	err = m.Create(movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

// Code から Media ID を取得する
func (r *MovieRepository) GetMediaIdByCode(ctx context.Context, code *string) (*int32, error) {
	wm := query.Use(r.DB).WatchMedia

	watchMedia, err := wm.WithContext(ctx).Where(wm.Code.Eq(*code)).First()
	if err != nil {
		return nil, err
	}
	return &watchMedia.ID, nil
}

func (r *MovieRepository) UpdateImpression(ctx context.Context, movieImpression *model.MovieImpression) (*gen.ResultInfo, error) {
	var err error
	defer func() {
		if err != nil {
			r.DB.Rollback()
		} else {
			r.DB.Commit()
		}
	}()

	mi := query.Use(r.DB).MovieImpression

	var result gen.ResultInfo
	result, err = mi.WithContext(ctx).Where(mi.ID.Eq(movieImpression.ID)).Updates(movieImpression)

	return &result, err
}

func (r *MovieRepository) SaveRecord(ctx context.Context, watchRecord *model.MovieWatchRecord) (*model.MovieWatchRecord, error) {
	var err error
	defer func() {
		if err != nil {
			r.DB.Rollback()
		} else {
			r.DB.Commit()
		}
	}()

	mwr := query.Use(r.DB).MovieWatchRecord

	err = mwr.WithContext(ctx).Create(watchRecord)

	return watchRecord, err
}
