package dao

import (
	"context"

	"{{.ModelPackage}}/internal/dao/model"
	"{{.ModelPackage}}/internal/dao/query"
)

func ({{.S}} *{{.DaoName}}) Example() (result *model.{{.ModelName}}, err error) {
	// example code
	ctx := context.Background()
	q := query.Use(dao.db).{{.ModelName}}
	return q.WithContext(ctx).First()
}
