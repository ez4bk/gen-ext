package ezgen

import (
	"unsafe"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func FindTableMeta(tables []any, tableName string) *QueryStructMeta {
	for _, elem := range tables {
		meta := ToQueryStructMeta(elem)
		if meta.TableName == tableName {
			return meta
		}
	}

	return nil
}

func ToQueryStructMeta(v any) *QueryStructMeta {
	return *(**QueryStructMeta)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + unsafe.Sizeof(uintptr(0)))) // skip eface.itab
}

// QueryStructMeta struct info in generated code
type QueryStructMeta struct {
	db *gorm.DB

	Generated       bool   // whether to generate db model
	FileName        string // generated file name
	S               string // the first letter(lower case)of simple Name (receiver)
	QueryStructName string // internal query struct name
	ModelStructName string // origin/model struct name
	TableName       string // table name in db server
	TableComment    string // table comment in db server
	StructInfo      ParserParam
	Fields          []*ModelField
	Source          SourceCode
	ImportPkgPaths  []string
	ModelMethods    []*ParserMethod // user custom method bind to db base struct

	interfaceMode bool
}

// ParserParam parameters in method
type ParserParam struct { // (user model.User)
	PkgPath   string // package's path: internal/model
	Package   string // package's name: model
	Name      string // param's name: user
	Type      string // param's type: User
	IsArray   bool   // is array or not
	IsPointer bool   // is pointer or not
}

// ModelField user input structures
type ModelField struct {
	Name             string
	Type             string
	ColumnName       string
	ColumnComment    string
	MultilineComment bool
	Tag              field.Tag
	GORMTag          field.GormTag
	CustomGenType    string
	Relation         *field.Relation
}

// SourceCode source code
type SourceCode int

// ParserMethod Apply to query struct and base struct custom method
type ParserMethod struct {
	Receiver   ParserParam
	MethodName string
	Doc        string
	Params     []ParserParam
	Result     []ParserParam
	Body       string
}
