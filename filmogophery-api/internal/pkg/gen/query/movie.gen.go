// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
)

func newMovie(db *gorm.DB, opts ...gen.DOOption) movie {
	_movie := movie{}

	_movie.movieDo.UseDB(db, opts...)
	_movie.movieDo.UseModel(&model.Movie{})

	tableName := _movie.movieDo.TableName()
	_movie.ALL = field.NewAsterisk(tableName)
	_movie.ID = field.NewInt32(tableName, "id")
	_movie.Title = field.NewString(tableName, "title")
	_movie.Overview = field.NewString(tableName, "overview")
	_movie.ReleaseDate = field.NewTime(tableName, "release_date")
	_movie.RunTime = field.NewInt32(tableName, "run_time")
	_movie.PosterURL = field.NewString(tableName, "poster_url")
	_movie.SeriesID = field.NewInt32(tableName, "series_id")
	_movie.TmdbID = field.NewInt32(tableName, "tmdb_id")
	_movie.Genres = movieManyToManyGenres{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Genres", "model.Genre"),
		Movies: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Genres.Movies", "model.Movie"),
		},
	}

	_movie.Series = movieHasOneSeries{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Series", "model.MovieSeries"),
	}

	_movie.MovieImpression = movieBelongsToMovieImpression{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("MovieImpression", "model.MovieImpression"),
		Movie: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("MovieImpression.Movie", "model.Movie"),
		},
		WatchRecords: struct {
			field.RelationField
			WatchMedia struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("MovieImpression.WatchRecords", "model.MovieWatchRecord"),
			WatchMedia: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("MovieImpression.WatchRecords.WatchMedia", "model.WatchMedia"),
			},
		},
	}

	_movie.fillFieldMap()

	return _movie
}

type movie struct {
	movieDo

	ALL         field.Asterisk
	ID          field.Int32
	Title       field.String
	Overview    field.String
	ReleaseDate field.Time
	RunTime     field.Int32
	PosterURL   field.String
	SeriesID    field.Int32
	TmdbID      field.Int32
	Genres      movieManyToManyGenres

	Series movieHasOneSeries

	MovieImpression movieBelongsToMovieImpression

	fieldMap map[string]field.Expr
}

func (m movie) Table(newTableName string) *movie {
	m.movieDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m movie) As(alias string) *movie {
	m.movieDo.DO = *(m.movieDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *movie) updateTableName(table string) *movie {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewInt32(table, "id")
	m.Title = field.NewString(table, "title")
	m.Overview = field.NewString(table, "overview")
	m.ReleaseDate = field.NewTime(table, "release_date")
	m.RunTime = field.NewInt32(table, "run_time")
	m.PosterURL = field.NewString(table, "poster_url")
	m.SeriesID = field.NewInt32(table, "series_id")
	m.TmdbID = field.NewInt32(table, "tmdb_id")

	m.fillFieldMap()

	return m
}

func (m *movie) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *movie) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 11)
	m.fieldMap["id"] = m.ID
	m.fieldMap["title"] = m.Title
	m.fieldMap["overview"] = m.Overview
	m.fieldMap["release_date"] = m.ReleaseDate
	m.fieldMap["run_time"] = m.RunTime
	m.fieldMap["poster_url"] = m.PosterURL
	m.fieldMap["series_id"] = m.SeriesID
	m.fieldMap["tmdb_id"] = m.TmdbID

}

func (m movie) clone(db *gorm.DB) movie {
	m.movieDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m movie) replaceDB(db *gorm.DB) movie {
	m.movieDo.ReplaceDB(db)
	return m
}

type movieManyToManyGenres struct {
	db *gorm.DB

	field.RelationField

	Movies struct {
		field.RelationField
	}
}

func (a movieManyToManyGenres) Where(conds ...field.Expr) *movieManyToManyGenres {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a movieManyToManyGenres) WithContext(ctx context.Context) *movieManyToManyGenres {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a movieManyToManyGenres) Session(session *gorm.Session) *movieManyToManyGenres {
	a.db = a.db.Session(session)
	return &a
}

func (a movieManyToManyGenres) Model(m *model.Movie) *movieManyToManyGenresTx {
	return &movieManyToManyGenresTx{a.db.Model(m).Association(a.Name())}
}

type movieManyToManyGenresTx struct{ tx *gorm.Association }

func (a movieManyToManyGenresTx) Find() (result []*model.Genre, err error) {
	return result, a.tx.Find(&result)
}

func (a movieManyToManyGenresTx) Append(values ...*model.Genre) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a movieManyToManyGenresTx) Replace(values ...*model.Genre) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a movieManyToManyGenresTx) Delete(values ...*model.Genre) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a movieManyToManyGenresTx) Clear() error {
	return a.tx.Clear()
}

func (a movieManyToManyGenresTx) Count() int64 {
	return a.tx.Count()
}

type movieHasOneSeries struct {
	db *gorm.DB

	field.RelationField
}

func (a movieHasOneSeries) Where(conds ...field.Expr) *movieHasOneSeries {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a movieHasOneSeries) WithContext(ctx context.Context) *movieHasOneSeries {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a movieHasOneSeries) Session(session *gorm.Session) *movieHasOneSeries {
	a.db = a.db.Session(session)
	return &a
}

func (a movieHasOneSeries) Model(m *model.Movie) *movieHasOneSeriesTx {
	return &movieHasOneSeriesTx{a.db.Model(m).Association(a.Name())}
}

type movieHasOneSeriesTx struct{ tx *gorm.Association }

func (a movieHasOneSeriesTx) Find() (result *model.MovieSeries, err error) {
	return result, a.tx.Find(&result)
}

func (a movieHasOneSeriesTx) Append(values ...*model.MovieSeries) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a movieHasOneSeriesTx) Replace(values ...*model.MovieSeries) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a movieHasOneSeriesTx) Delete(values ...*model.MovieSeries) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a movieHasOneSeriesTx) Clear() error {
	return a.tx.Clear()
}

func (a movieHasOneSeriesTx) Count() int64 {
	return a.tx.Count()
}

type movieBelongsToMovieImpression struct {
	db *gorm.DB

	field.RelationField

	Movie struct {
		field.RelationField
	}
	WatchRecords struct {
		field.RelationField
		WatchMedia struct {
			field.RelationField
		}
	}
}

func (a movieBelongsToMovieImpression) Where(conds ...field.Expr) *movieBelongsToMovieImpression {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a movieBelongsToMovieImpression) WithContext(ctx context.Context) *movieBelongsToMovieImpression {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a movieBelongsToMovieImpression) Session(session *gorm.Session) *movieBelongsToMovieImpression {
	a.db = a.db.Session(session)
	return &a
}

func (a movieBelongsToMovieImpression) Model(m *model.Movie) *movieBelongsToMovieImpressionTx {
	return &movieBelongsToMovieImpressionTx{a.db.Model(m).Association(a.Name())}
}

type movieBelongsToMovieImpressionTx struct{ tx *gorm.Association }

func (a movieBelongsToMovieImpressionTx) Find() (result *model.MovieImpression, err error) {
	return result, a.tx.Find(&result)
}

func (a movieBelongsToMovieImpressionTx) Append(values ...*model.MovieImpression) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a movieBelongsToMovieImpressionTx) Replace(values ...*model.MovieImpression) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a movieBelongsToMovieImpressionTx) Delete(values ...*model.MovieImpression) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a movieBelongsToMovieImpressionTx) Clear() error {
	return a.tx.Clear()
}

func (a movieBelongsToMovieImpressionTx) Count() int64 {
	return a.tx.Count()
}

type movieDo struct{ gen.DO }

func (m movieDo) Debug() *movieDo {
	return m.withDO(m.DO.Debug())
}

func (m movieDo) WithContext(ctx context.Context) *movieDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m movieDo) ReadDB() *movieDo {
	return m.Clauses(dbresolver.Read)
}

func (m movieDo) WriteDB() *movieDo {
	return m.Clauses(dbresolver.Write)
}

func (m movieDo) Session(config *gorm.Session) *movieDo {
	return m.withDO(m.DO.Session(config))
}

func (m movieDo) Clauses(conds ...clause.Expression) *movieDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m movieDo) Returning(value interface{}, columns ...string) *movieDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m movieDo) Not(conds ...gen.Condition) *movieDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m movieDo) Or(conds ...gen.Condition) *movieDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m movieDo) Select(conds ...field.Expr) *movieDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m movieDo) Where(conds ...gen.Condition) *movieDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m movieDo) Order(conds ...field.Expr) *movieDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m movieDo) Distinct(cols ...field.Expr) *movieDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m movieDo) Omit(cols ...field.Expr) *movieDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m movieDo) Join(table schema.Tabler, on ...field.Expr) *movieDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m movieDo) LeftJoin(table schema.Tabler, on ...field.Expr) *movieDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m movieDo) RightJoin(table schema.Tabler, on ...field.Expr) *movieDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m movieDo) Group(cols ...field.Expr) *movieDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m movieDo) Having(conds ...gen.Condition) *movieDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m movieDo) Limit(limit int) *movieDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m movieDo) Offset(offset int) *movieDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m movieDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *movieDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m movieDo) Unscoped() *movieDo {
	return m.withDO(m.DO.Unscoped())
}

func (m movieDo) Create(values ...*model.Movie) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m movieDo) CreateInBatches(values []*model.Movie, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m movieDo) Save(values ...*model.Movie) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m movieDo) First() (*model.Movie, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Movie), nil
	}
}

func (m movieDo) Take() (*model.Movie, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Movie), nil
	}
}

func (m movieDo) Last() (*model.Movie, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Movie), nil
	}
}

func (m movieDo) Find() ([]*model.Movie, error) {
	result, err := m.DO.Find()
	return result.([]*model.Movie), err
}

func (m movieDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Movie, err error) {
	buf := make([]*model.Movie, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m movieDo) FindInBatches(result *[]*model.Movie, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m movieDo) Attrs(attrs ...field.AssignExpr) *movieDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m movieDo) Assign(attrs ...field.AssignExpr) *movieDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m movieDo) Joins(fields ...field.RelationField) *movieDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m movieDo) Preload(fields ...field.RelationField) *movieDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m movieDo) FirstOrInit() (*model.Movie, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Movie), nil
	}
}

func (m movieDo) FirstOrCreate() (*model.Movie, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Movie), nil
	}
}

func (m movieDo) FindByPage(offset int, limit int) (result []*model.Movie, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m movieDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m movieDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m movieDo) Delete(models ...*model.Movie) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *movieDo) withDO(do gen.Dao) *movieDo {
	m.DO = *do.(*gen.DO)
	return m
}