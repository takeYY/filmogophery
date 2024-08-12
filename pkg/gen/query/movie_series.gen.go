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

func newMovieSeries(db *gorm.DB, opts ...gen.DOOption) movieSeries {
	_movieSeries := movieSeries{}

	_movieSeries.movieSeriesDo.UseDB(db, opts...)
	_movieSeries.movieSeriesDo.UseModel(&model.MovieSeries{})

	tableName := _movieSeries.movieSeriesDo.TableName()
	_movieSeries.ALL = field.NewAsterisk(tableName)
	_movieSeries.ID = field.NewInt32(tableName, "id")
	_movieSeries.Name = field.NewString(tableName, "name")
	_movieSeries.PosterID = field.NewInt32(tableName, "poster_id")
	_movieSeries.Poster = movieSeriesHasOnePoster{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Poster", "model.Poster"),
	}

	_movieSeries.fillFieldMap()

	return _movieSeries
}

type movieSeries struct {
	movieSeriesDo

	ALL      field.Asterisk
	ID       field.Int32
	Name     field.String
	PosterID field.Int32
	Poster   movieSeriesHasOnePoster

	fieldMap map[string]field.Expr
}

func (m movieSeries) Table(newTableName string) *movieSeries {
	m.movieSeriesDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m movieSeries) As(alias string) *movieSeries {
	m.movieSeriesDo.DO = *(m.movieSeriesDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *movieSeries) updateTableName(table string) *movieSeries {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewInt32(table, "id")
	m.Name = field.NewString(table, "name")
	m.PosterID = field.NewInt32(table, "poster_id")

	m.fillFieldMap()

	return m
}

func (m *movieSeries) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *movieSeries) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 4)
	m.fieldMap["id"] = m.ID
	m.fieldMap["name"] = m.Name
	m.fieldMap["poster_id"] = m.PosterID

}

func (m movieSeries) clone(db *gorm.DB) movieSeries {
	m.movieSeriesDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m movieSeries) replaceDB(db *gorm.DB) movieSeries {
	m.movieSeriesDo.ReplaceDB(db)
	return m
}

type movieSeriesHasOnePoster struct {
	db *gorm.DB

	field.RelationField
}

func (a movieSeriesHasOnePoster) Where(conds ...field.Expr) *movieSeriesHasOnePoster {
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

func (a movieSeriesHasOnePoster) WithContext(ctx context.Context) *movieSeriesHasOnePoster {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a movieSeriesHasOnePoster) Session(session *gorm.Session) *movieSeriesHasOnePoster {
	a.db = a.db.Session(session)
	return &a
}

func (a movieSeriesHasOnePoster) Model(m *model.MovieSeries) *movieSeriesHasOnePosterTx {
	return &movieSeriesHasOnePosterTx{a.db.Model(m).Association(a.Name())}
}

type movieSeriesHasOnePosterTx struct{ tx *gorm.Association }

func (a movieSeriesHasOnePosterTx) Find() (result *model.Poster, err error) {
	return result, a.tx.Find(&result)
}

func (a movieSeriesHasOnePosterTx) Append(values ...*model.Poster) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a movieSeriesHasOnePosterTx) Replace(values ...*model.Poster) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a movieSeriesHasOnePosterTx) Delete(values ...*model.Poster) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a movieSeriesHasOnePosterTx) Clear() error {
	return a.tx.Clear()
}

func (a movieSeriesHasOnePosterTx) Count() int64 {
	return a.tx.Count()
}

type movieSeriesDo struct{ gen.DO }

func (m movieSeriesDo) Debug() *movieSeriesDo {
	return m.withDO(m.DO.Debug())
}

func (m movieSeriesDo) WithContext(ctx context.Context) *movieSeriesDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m movieSeriesDo) ReadDB() *movieSeriesDo {
	return m.Clauses(dbresolver.Read)
}

func (m movieSeriesDo) WriteDB() *movieSeriesDo {
	return m.Clauses(dbresolver.Write)
}

func (m movieSeriesDo) Session(config *gorm.Session) *movieSeriesDo {
	return m.withDO(m.DO.Session(config))
}

func (m movieSeriesDo) Clauses(conds ...clause.Expression) *movieSeriesDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m movieSeriesDo) Returning(value interface{}, columns ...string) *movieSeriesDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m movieSeriesDo) Not(conds ...gen.Condition) *movieSeriesDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m movieSeriesDo) Or(conds ...gen.Condition) *movieSeriesDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m movieSeriesDo) Select(conds ...field.Expr) *movieSeriesDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m movieSeriesDo) Where(conds ...gen.Condition) *movieSeriesDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m movieSeriesDo) Order(conds ...field.Expr) *movieSeriesDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m movieSeriesDo) Distinct(cols ...field.Expr) *movieSeriesDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m movieSeriesDo) Omit(cols ...field.Expr) *movieSeriesDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m movieSeriesDo) Join(table schema.Tabler, on ...field.Expr) *movieSeriesDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m movieSeriesDo) LeftJoin(table schema.Tabler, on ...field.Expr) *movieSeriesDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m movieSeriesDo) RightJoin(table schema.Tabler, on ...field.Expr) *movieSeriesDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m movieSeriesDo) Group(cols ...field.Expr) *movieSeriesDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m movieSeriesDo) Having(conds ...gen.Condition) *movieSeriesDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m movieSeriesDo) Limit(limit int) *movieSeriesDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m movieSeriesDo) Offset(offset int) *movieSeriesDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m movieSeriesDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *movieSeriesDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m movieSeriesDo) Unscoped() *movieSeriesDo {
	return m.withDO(m.DO.Unscoped())
}

func (m movieSeriesDo) Create(values ...*model.MovieSeries) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m movieSeriesDo) CreateInBatches(values []*model.MovieSeries, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m movieSeriesDo) Save(values ...*model.MovieSeries) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m movieSeriesDo) First() (*model.MovieSeries, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.MovieSeries), nil
	}
}

func (m movieSeriesDo) Take() (*model.MovieSeries, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.MovieSeries), nil
	}
}

func (m movieSeriesDo) Last() (*model.MovieSeries, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.MovieSeries), nil
	}
}

func (m movieSeriesDo) Find() ([]*model.MovieSeries, error) {
	result, err := m.DO.Find()
	return result.([]*model.MovieSeries), err
}

func (m movieSeriesDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.MovieSeries, err error) {
	buf := make([]*model.MovieSeries, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m movieSeriesDo) FindInBatches(result *[]*model.MovieSeries, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m movieSeriesDo) Attrs(attrs ...field.AssignExpr) *movieSeriesDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m movieSeriesDo) Assign(attrs ...field.AssignExpr) *movieSeriesDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m movieSeriesDo) Joins(fields ...field.RelationField) *movieSeriesDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m movieSeriesDo) Preload(fields ...field.RelationField) *movieSeriesDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m movieSeriesDo) FirstOrInit() (*model.MovieSeries, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.MovieSeries), nil
	}
}

func (m movieSeriesDo) FirstOrCreate() (*model.MovieSeries, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.MovieSeries), nil
	}
}

func (m movieSeriesDo) FindByPage(offset int, limit int) (result []*model.MovieSeries, count int64, err error) {
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

func (m movieSeriesDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m movieSeriesDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m movieSeriesDo) Delete(models ...*model.MovieSeries) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *movieSeriesDo) withDO(do gen.Dao) *movieSeriesDo {
	m.DO = *do.(*gen.DO)
	return m
}
