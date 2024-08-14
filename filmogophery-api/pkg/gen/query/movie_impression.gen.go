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

	"filmogophery/pkg/gen/model"
)

func newMovieImpression(db *gorm.DB, opts ...gen.DOOption) movieImpression {
	_movieImpression := movieImpression{}

	_movieImpression.movieImpressionDo.UseDB(db, opts...)
	_movieImpression.movieImpressionDo.UseModel(&model.MovieImpression{})

	tableName := _movieImpression.movieImpressionDo.TableName()
	_movieImpression.ALL = field.NewAsterisk(tableName)
	_movieImpression.ID = field.NewInt32(tableName, "id")
	_movieImpression.MovieID = field.NewInt32(tableName, "movie_id")
	_movieImpression.Status = field.NewBool(tableName, "status")
	_movieImpression.Rating = field.NewBool(tableName, "rating")
	_movieImpression.Note = field.NewString(tableName, "note")
	_movieImpression.Movie = movieImpressionHasOneMovie{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Movie", "model.Movie"),
		Genres: struct {
			field.RelationField
			Movies struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("Movie.Genres", "model.Genre"),
			Movies: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Movie.Genres.Movies", "model.Movie"),
			},
		},
		Poster: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Movie.Poster", "model.Poster"),
		},
		Series: struct {
			field.RelationField
			Poster struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("Movie.Series", "model.MovieSeries"),
			Poster: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Movie.Series.Poster", "model.Poster"),
			},
		},
	}

	_movieImpression.fillFieldMap()

	return _movieImpression
}

type movieImpression struct {
	movieImpressionDo

	ALL     field.Asterisk
	ID      field.Int32
	MovieID field.Int32
	Status  field.Bool
	Rating  field.Bool
	Note    field.String
	Movie   movieImpressionHasOneMovie

	fieldMap map[string]field.Expr
}

func (m movieImpression) Table(newTableName string) *movieImpression {
	m.movieImpressionDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m movieImpression) As(alias string) *movieImpression {
	m.movieImpressionDo.DO = *(m.movieImpressionDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *movieImpression) updateTableName(table string) *movieImpression {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewInt32(table, "id")
	m.MovieID = field.NewInt32(table, "movie_id")
	m.Status = field.NewBool(table, "status")
	m.Rating = field.NewBool(table, "rating")
	m.Note = field.NewString(table, "note")

	m.fillFieldMap()

	return m
}

func (m *movieImpression) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *movieImpression) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 6)
	m.fieldMap["id"] = m.ID
	m.fieldMap["movie_id"] = m.MovieID
	m.fieldMap["status"] = m.Status
	m.fieldMap["rating"] = m.Rating
	m.fieldMap["note"] = m.Note

}

func (m movieImpression) clone(db *gorm.DB) movieImpression {
	m.movieImpressionDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m movieImpression) replaceDB(db *gorm.DB) movieImpression {
	m.movieImpressionDo.ReplaceDB(db)
	return m
}

type movieImpressionHasOneMovie struct {
	db *gorm.DB

	field.RelationField

	Genres struct {
		field.RelationField
		Movies struct {
			field.RelationField
		}
	}
	Poster struct {
		field.RelationField
	}
	Series struct {
		field.RelationField
		Poster struct {
			field.RelationField
		}
	}
}

func (a movieImpressionHasOneMovie) Where(conds ...field.Expr) *movieImpressionHasOneMovie {
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

func (a movieImpressionHasOneMovie) WithContext(ctx context.Context) *movieImpressionHasOneMovie {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a movieImpressionHasOneMovie) Session(session *gorm.Session) *movieImpressionHasOneMovie {
	a.db = a.db.Session(session)
	return &a
}

func (a movieImpressionHasOneMovie) Model(m *model.MovieImpression) *movieImpressionHasOneMovieTx {
	return &movieImpressionHasOneMovieTx{a.db.Model(m).Association(a.Name())}
}

type movieImpressionHasOneMovieTx struct{ tx *gorm.Association }

func (a movieImpressionHasOneMovieTx) Find() (result *model.Movie, err error) {
	return result, a.tx.Find(&result)
}

func (a movieImpressionHasOneMovieTx) Append(values ...*model.Movie) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a movieImpressionHasOneMovieTx) Replace(values ...*model.Movie) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a movieImpressionHasOneMovieTx) Delete(values ...*model.Movie) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a movieImpressionHasOneMovieTx) Clear() error {
	return a.tx.Clear()
}

func (a movieImpressionHasOneMovieTx) Count() int64 {
	return a.tx.Count()
}

type movieImpressionDo struct{ gen.DO }

func (m movieImpressionDo) Debug() *movieImpressionDo {
	return m.withDO(m.DO.Debug())
}

func (m movieImpressionDo) WithContext(ctx context.Context) *movieImpressionDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m movieImpressionDo) ReadDB() *movieImpressionDo {
	return m.Clauses(dbresolver.Read)
}

func (m movieImpressionDo) WriteDB() *movieImpressionDo {
	return m.Clauses(dbresolver.Write)
}

func (m movieImpressionDo) Session(config *gorm.Session) *movieImpressionDo {
	return m.withDO(m.DO.Session(config))
}

func (m movieImpressionDo) Clauses(conds ...clause.Expression) *movieImpressionDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m movieImpressionDo) Returning(value interface{}, columns ...string) *movieImpressionDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m movieImpressionDo) Not(conds ...gen.Condition) *movieImpressionDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m movieImpressionDo) Or(conds ...gen.Condition) *movieImpressionDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m movieImpressionDo) Select(conds ...field.Expr) *movieImpressionDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m movieImpressionDo) Where(conds ...gen.Condition) *movieImpressionDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m movieImpressionDo) Order(conds ...field.Expr) *movieImpressionDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m movieImpressionDo) Distinct(cols ...field.Expr) *movieImpressionDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m movieImpressionDo) Omit(cols ...field.Expr) *movieImpressionDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m movieImpressionDo) Join(table schema.Tabler, on ...field.Expr) *movieImpressionDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m movieImpressionDo) LeftJoin(table schema.Tabler, on ...field.Expr) *movieImpressionDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m movieImpressionDo) RightJoin(table schema.Tabler, on ...field.Expr) *movieImpressionDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m movieImpressionDo) Group(cols ...field.Expr) *movieImpressionDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m movieImpressionDo) Having(conds ...gen.Condition) *movieImpressionDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m movieImpressionDo) Limit(limit int) *movieImpressionDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m movieImpressionDo) Offset(offset int) *movieImpressionDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m movieImpressionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *movieImpressionDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m movieImpressionDo) Unscoped() *movieImpressionDo {
	return m.withDO(m.DO.Unscoped())
}

func (m movieImpressionDo) Create(values ...*model.MovieImpression) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m movieImpressionDo) CreateInBatches(values []*model.MovieImpression, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m movieImpressionDo) Save(values ...*model.MovieImpression) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m movieImpressionDo) First() (*model.MovieImpression, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.MovieImpression), nil
	}
}

func (m movieImpressionDo) Take() (*model.MovieImpression, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.MovieImpression), nil
	}
}

func (m movieImpressionDo) Last() (*model.MovieImpression, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.MovieImpression), nil
	}
}

func (m movieImpressionDo) Find() ([]*model.MovieImpression, error) {
	result, err := m.DO.Find()
	return result.([]*model.MovieImpression), err
}

func (m movieImpressionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.MovieImpression, err error) {
	buf := make([]*model.MovieImpression, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m movieImpressionDo) FindInBatches(result *[]*model.MovieImpression, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m movieImpressionDo) Attrs(attrs ...field.AssignExpr) *movieImpressionDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m movieImpressionDo) Assign(attrs ...field.AssignExpr) *movieImpressionDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m movieImpressionDo) Joins(fields ...field.RelationField) *movieImpressionDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m movieImpressionDo) Preload(fields ...field.RelationField) *movieImpressionDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m movieImpressionDo) FirstOrInit() (*model.MovieImpression, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.MovieImpression), nil
	}
}

func (m movieImpressionDo) FirstOrCreate() (*model.MovieImpression, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.MovieImpression), nil
	}
}

func (m movieImpressionDo) FindByPage(offset int, limit int) (result []*model.MovieImpression, count int64, err error) {
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

func (m movieImpressionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m movieImpressionDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m movieImpressionDo) Delete(models ...*model.MovieImpression) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *movieImpressionDo) withDO(do gen.Dao) *movieImpressionDo {
	m.DO = *do.(*gen.DO)
	return m
}