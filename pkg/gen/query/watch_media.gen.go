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

func newWatchMedia(db *gorm.DB, opts ...gen.DOOption) watchMedia {
	_watchMedia := watchMedia{}

	_watchMedia.watchMediaDo.UseDB(db, opts...)
	_watchMedia.watchMediaDo.UseModel(&model.WatchMedia{})

	tableName := _watchMedia.watchMediaDo.TableName()
	_watchMedia.ALL = field.NewAsterisk(tableName)
	_watchMedia.ID = field.NewInt64(tableName, "id")
	_watchMedia.CreatedAt = field.NewTime(tableName, "created_at")
	_watchMedia.UpdatedAt = field.NewTime(tableName, "updated_at")
	_watchMedia.DeletedAt = field.NewField(tableName, "deleted_at")
	_watchMedia.Code = field.NewString(tableName, "code")
	_watchMedia.Name = field.NewString(tableName, "name")

	_watchMedia.fillFieldMap()

	return _watchMedia
}

type watchMedia struct {
	watchMediaDo

	ALL       field.Asterisk
	ID        field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	Code      field.String
	Name      field.String

	fieldMap map[string]field.Expr
}

func (w watchMedia) Table(newTableName string) *watchMedia {
	w.watchMediaDo.UseTable(newTableName)
	return w.updateTableName(newTableName)
}

func (w watchMedia) As(alias string) *watchMedia {
	w.watchMediaDo.DO = *(w.watchMediaDo.As(alias).(*gen.DO))
	return w.updateTableName(alias)
}

func (w *watchMedia) updateTableName(table string) *watchMedia {
	w.ALL = field.NewAsterisk(table)
	w.ID = field.NewInt64(table, "id")
	w.CreatedAt = field.NewTime(table, "created_at")
	w.UpdatedAt = field.NewTime(table, "updated_at")
	w.DeletedAt = field.NewField(table, "deleted_at")
	w.Code = field.NewString(table, "code")
	w.Name = field.NewString(table, "name")

	w.fillFieldMap()

	return w
}

func (w *watchMedia) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := w.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (w *watchMedia) fillFieldMap() {
	w.fieldMap = make(map[string]field.Expr, 6)
	w.fieldMap["id"] = w.ID
	w.fieldMap["created_at"] = w.CreatedAt
	w.fieldMap["updated_at"] = w.UpdatedAt
	w.fieldMap["deleted_at"] = w.DeletedAt
	w.fieldMap["code"] = w.Code
	w.fieldMap["name"] = w.Name
}

func (w watchMedia) clone(db *gorm.DB) watchMedia {
	w.watchMediaDo.ReplaceConnPool(db.Statement.ConnPool)
	return w
}

func (w watchMedia) replaceDB(db *gorm.DB) watchMedia {
	w.watchMediaDo.ReplaceDB(db)
	return w
}

type watchMediaDo struct{ gen.DO }

func (w watchMediaDo) Debug() *watchMediaDo {
	return w.withDO(w.DO.Debug())
}

func (w watchMediaDo) WithContext(ctx context.Context) *watchMediaDo {
	return w.withDO(w.DO.WithContext(ctx))
}

func (w watchMediaDo) ReadDB() *watchMediaDo {
	return w.Clauses(dbresolver.Read)
}

func (w watchMediaDo) WriteDB() *watchMediaDo {
	return w.Clauses(dbresolver.Write)
}

func (w watchMediaDo) Session(config *gorm.Session) *watchMediaDo {
	return w.withDO(w.DO.Session(config))
}

func (w watchMediaDo) Clauses(conds ...clause.Expression) *watchMediaDo {
	return w.withDO(w.DO.Clauses(conds...))
}

func (w watchMediaDo) Returning(value interface{}, columns ...string) *watchMediaDo {
	return w.withDO(w.DO.Returning(value, columns...))
}

func (w watchMediaDo) Not(conds ...gen.Condition) *watchMediaDo {
	return w.withDO(w.DO.Not(conds...))
}

func (w watchMediaDo) Or(conds ...gen.Condition) *watchMediaDo {
	return w.withDO(w.DO.Or(conds...))
}

func (w watchMediaDo) Select(conds ...field.Expr) *watchMediaDo {
	return w.withDO(w.DO.Select(conds...))
}

func (w watchMediaDo) Where(conds ...gen.Condition) *watchMediaDo {
	return w.withDO(w.DO.Where(conds...))
}

func (w watchMediaDo) Order(conds ...field.Expr) *watchMediaDo {
	return w.withDO(w.DO.Order(conds...))
}

func (w watchMediaDo) Distinct(cols ...field.Expr) *watchMediaDo {
	return w.withDO(w.DO.Distinct(cols...))
}

func (w watchMediaDo) Omit(cols ...field.Expr) *watchMediaDo {
	return w.withDO(w.DO.Omit(cols...))
}

func (w watchMediaDo) Join(table schema.Tabler, on ...field.Expr) *watchMediaDo {
	return w.withDO(w.DO.Join(table, on...))
}

func (w watchMediaDo) LeftJoin(table schema.Tabler, on ...field.Expr) *watchMediaDo {
	return w.withDO(w.DO.LeftJoin(table, on...))
}

func (w watchMediaDo) RightJoin(table schema.Tabler, on ...field.Expr) *watchMediaDo {
	return w.withDO(w.DO.RightJoin(table, on...))
}

func (w watchMediaDo) Group(cols ...field.Expr) *watchMediaDo {
	return w.withDO(w.DO.Group(cols...))
}

func (w watchMediaDo) Having(conds ...gen.Condition) *watchMediaDo {
	return w.withDO(w.DO.Having(conds...))
}

func (w watchMediaDo) Limit(limit int) *watchMediaDo {
	return w.withDO(w.DO.Limit(limit))
}

func (w watchMediaDo) Offset(offset int) *watchMediaDo {
	return w.withDO(w.DO.Offset(offset))
}

func (w watchMediaDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *watchMediaDo {
	return w.withDO(w.DO.Scopes(funcs...))
}

func (w watchMediaDo) Unscoped() *watchMediaDo {
	return w.withDO(w.DO.Unscoped())
}

func (w watchMediaDo) Create(values ...*model.WatchMedia) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Create(values)
}

func (w watchMediaDo) CreateInBatches(values []*model.WatchMedia, batchSize int) error {
	return w.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (w watchMediaDo) Save(values ...*model.WatchMedia) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Save(values)
}

func (w watchMediaDo) First() (*model.WatchMedia, error) {
	if result, err := w.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.WatchMedia), nil
	}
}

func (w watchMediaDo) Take() (*model.WatchMedia, error) {
	if result, err := w.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.WatchMedia), nil
	}
}

func (w watchMediaDo) Last() (*model.WatchMedia, error) {
	if result, err := w.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.WatchMedia), nil
	}
}

func (w watchMediaDo) Find() ([]*model.WatchMedia, error) {
	result, err := w.DO.Find()
	return result.([]*model.WatchMedia), err
}

func (w watchMediaDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.WatchMedia, err error) {
	buf := make([]*model.WatchMedia, 0, batchSize)
	err = w.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (w watchMediaDo) FindInBatches(result *[]*model.WatchMedia, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return w.DO.FindInBatches(result, batchSize, fc)
}

func (w watchMediaDo) Attrs(attrs ...field.AssignExpr) *watchMediaDo {
	return w.withDO(w.DO.Attrs(attrs...))
}

func (w watchMediaDo) Assign(attrs ...field.AssignExpr) *watchMediaDo {
	return w.withDO(w.DO.Assign(attrs...))
}

func (w watchMediaDo) Joins(fields ...field.RelationField) *watchMediaDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Joins(_f))
	}
	return &w
}

func (w watchMediaDo) Preload(fields ...field.RelationField) *watchMediaDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Preload(_f))
	}
	return &w
}

func (w watchMediaDo) FirstOrInit() (*model.WatchMedia, error) {
	if result, err := w.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.WatchMedia), nil
	}
}

func (w watchMediaDo) FirstOrCreate() (*model.WatchMedia, error) {
	if result, err := w.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.WatchMedia), nil
	}
}

func (w watchMediaDo) FindByPage(offset int, limit int) (result []*model.WatchMedia, count int64, err error) {
	result, err = w.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = w.Offset(-1).Limit(-1).Count()
	return
}

func (w watchMediaDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = w.Count()
	if err != nil {
		return
	}

	err = w.Offset(offset).Limit(limit).Scan(result)
	return
}

func (w watchMediaDo) Scan(result interface{}) (err error) {
	return w.DO.Scan(result)
}

func (w watchMediaDo) Delete(models ...*model.WatchMedia) (result gen.ResultInfo, err error) {
	return w.DO.Delete(models)
}

func (w *watchMediaDo) withDO(do gen.Dao) *watchMediaDo {
	w.DO = *do.(*gen.DO)
	return w
}