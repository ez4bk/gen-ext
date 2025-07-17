package dao

import (
	"context"

	"{{.ModelPackage}}/internal/dao/model"
)

func ({{.S}} *{{.DaoName}}) Example(ctx context.Context) (result *model.{{.ModelName}}, err error) {
	// example code
	return {{.S}}.WithContext(ctx).First()
}
