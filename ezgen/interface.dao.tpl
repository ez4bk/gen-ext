package dao

import (
	"context"
	"reflect"

	"{{.ModelPackage}}/internal/dao/model"
	"{{.ModelPackage}}/internal/dao/query"

	"{{.ModelPackage}}/cmd/gen/ezgen"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	{{range .ImportPkgPaths}}{{.}} ` + "\n" + `{{end}}
)

type i{{.ModelName}}Dao interface {
	Add(ctx context.Context, data ...*model.{{.ModelName}}) (err error)
	Get(ctx context.Context, id {{.PKType}}, opts ...GetOption) (result *model.{{.ModelName}}, err error)
	// List returns the specified models from database by params
	List(ctx context.Context, params *List{{.ModelName}}Params) (list []*model.{{.ModelName}}, total int64, err error)
	Update(ctx context.Context, data *model.{{.ModelName}}) (err error)
	// Delete soft deletes data
	Delete(ctx context.Context, id {{.PKType}}) (err error)
	// Destroy hard deletes data
	Destroy(ctx context.Context, id {{.PKType}}) (err error)

	// Custom methods goes here
	Custom(ctx context.Context, data *model.{{.ModelName}}) (err error)

}

type {{.DaoName}} struct {
	db *gorm.DB
}

func (dao *{{.DaoName}}) Custom(ctx context.Context, data *model.{{.ModelName}}) (err error) {
	// q := query.Use(dao.db).{{.ModelName}}
    return
}
