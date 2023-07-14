package orm

import (
	"go-train/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	tagColumn = "column"
)

type Registry interface {
	Get(val any) error
	Register(val any) (*Model, error)
}

type Model struct {
	tableName string
	// 字段名到字段的映射
	fieldsMap map[string]*Field
	// 列名到字段定义
	columnMap map[string]*Field
}

type ModelOption func(m *Model) error

type Field struct {
	goName  string       // 字段名
	colName string       // 列名
	typ     reflect.Type // 字段类型
}

type registry struct {
	models sync.Map
}

func newRegistry() *registry {
	return &registry{
		models: sync.Map{},
	}
}

// Get double check lock
func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	var err error
	m, err = r.Register(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*Model), nil
}

// Register 解析模型
func (r *registry) Register(entity any, opts ...ModelOption) (*Model, error) {
	typ := reflect.TypeOf(entity)
	// 限制只能用一级指针
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typElem := typ.Elem()
	numFields := typElem.NumField()
	fieldMap := make(map[string]*Field, numFields)
	columnMap := make(map[string]*Field, numFields)
	for i := 0; i < numFields; i++ {
		fd := typElem.Field(i)
		pair, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, nil
		}
		columnName := pair[tagColumn]
		if columnName == "" {
			columnName = underscoreName(fd.Name)
		}
		fdMeta := &Field{
			goName:  fd.Name,
			colName: columnName,
			typ:     fd.Type,
		}
		fieldMap[fd.Name] = fdMeta
		columnMap[columnName] = fdMeta
	}
	var tableName string
	// 断言，看看是否实现了 TableName 接口
	if tbl, ok := entity.(TableName); ok {
		tableName = tbl.TableName()
	} else {
		tableName = underscoreName(typElem.Name())
	}
	res := &Model{
		tableName: tableName,
		fieldsMap: fieldMap,
		columnMap: columnMap,
	}

	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, res)

	return res, nil
}

func ModelWithTableName(tableName string) ModelOption {
	return func(m *Model) error {
		m.tableName = tableName
		return nil
	}
}

func ModelWithColumnName(field string, colName string) ModelOption {
	return func(m *Model) error {
		fd, ok := m.fieldsMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.colName = colName
		return nil
	}
}

// Register 解析模型
func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}
	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		key := segs[0]
		val := segs[1]
		res[key] = val
	}
	return res, nil
}

// 大小写转换
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)

}
