package ezgen

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DefaultModelOpt = []gen.ModelOpt{
	// 乐观锁
	gen.FieldType("version", "optimisticlock.Version"),
	gen.FieldType("is_deleted", "soft_delete.DeletedAt"),
	gen.FieldModify(func(f gen.Field) gen.Field {
		if f.ColumnName == "is_deleted" {
			// 添加删除标志
			f.GORMTag.Set("softDelete", "flag")
		}
		return f
	}),
}

func GeneratorForeignKey(g *gen.Generator, db *gorm.DB, tableName string) []gen.ModelOpt {
	var opts []gen.ModelOpt
	foreignKeyInfos, _ := GetReferencingForeignKeys(db, tableName)
	for _, info := range foreignKeyInfos {
		fmt.Printf("%#v\n", info)
		relationship, suffix, err := GetRelationship(info.ConstraintName)
		if err != nil {
			fmt.Println(err)
			continue
		}

		opts = append(opts, gen.FieldRelate(relationship, SnakeToPascalCase(info.ReferencedTable)+suffix, g.GenerateModel(info.ReferencedTable),
			&field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": append([]string{}, info.ColumnName), "references": append([]string{}, info.ReferencedColumn)},
			}))

	}
	return opts
}

func GetRelationship(constraintName string) (field.RelationshipType, string, error) {
	if strings.Contains(constraintName, string(schema.HasOne)) {
		return field.HasOne, "", nil
	}
	if strings.Contains(constraintName, string(schema.HasMany)) {
		return field.HasMany, "List", nil
	}
	// if strings.Contains(constraintName, string(schema.BelongsTo)) {
	// 	return field.BelongsTo, nil
	// }
	// if strings.Contains(constraintName, string(schema.Many2Many)) {
	// 	return field.Many2Many, nil
	// }
	return field.HasOne, "", errors.New("unknown constraint name: " + constraintName)
}

// ForeignKeyInfo 表示外键约束信息
type ForeignKeyInfo struct {
	ConstraintName   string `gorm:"column:CONSTRAINT_NAME"`
	TableName        string `gorm:"column:TABLE_NAME"`
	ColumnName       string `gorm:"column:COLUMN_NAME"`
	ReferencedTable  string `gorm:"column:REFERENCED_TABLE_NAME"`
	ReferencedColumn string `gorm:"column:REFERENCED_COLUMN_NAME"`
	UpdateRule       string `gorm:"column:UPDATE_RULE"`
	DeleteRule       string `gorm:"column:DELETE_RULE"`
}

// GetReferencingForeignKeys 查询引用特定表的外键约束
func GetReferencingForeignKeys(db *gorm.DB, tableName string) ([]ForeignKeyInfo, error) {
	databaseName, err := GetDatabaseName(db)
	if err != nil {
		return nil, err
	}

	var referencingForeignKeys []ForeignKeyInfo

	result := db.Table("INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS k").
		Select(`
			k.CONSTRAINT_NAME,
			k.TABLE_NAME,
			k.COLUMN_NAME,
			k.REFERENCED_TABLE_NAME,
			k.REFERENCED_COLUMN_NAME,
			r.UPDATE_RULE,
			r.DELETE_RULE
		`).
		Joins("JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS AS r ON k.CONSTRAINT_NAME = r.CONSTRAINT_NAME AND k.TABLE_SCHEMA = r.CONSTRAINT_SCHEMA").
		Where("k.TABLE_SCHEMA = ? AND k.TABLE_NAME = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL", databaseName, tableName).
		Order("k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION").
		Find(&referencingForeignKeys)

	if result.Error != nil {
		err := fmt.Errorf("查询引用外键失败: %v", result.Error)
		return nil, err
	}

	return referencingForeignKeys, nil
}

// GetDatabaseName 获取当前连接的数据库名称
func GetDatabaseName(db *gorm.DB) (string, error) {
	var dbName string
	result := db.Raw("SELECT DATABASE()").Scan(&dbName)
	return dbName, result.Error
}
