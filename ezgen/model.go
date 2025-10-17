package ezgen

import (
	"fmt"
	"strings"

	"gorm.io/gen"
	"gorm.io/gorm"
)

func TypeNullable(columnType gorm.ColumnType, dataType string) string {
	if n, ok := columnType.Nullable(); ok && n {
		return fmt.Sprintf("*%s", dataType)
	} else {
		return dataType
	}
}

func GetDataMapMySQL(cfg *gen.Config, customMap map[string]func(gorm.ColumnType) (dataType string)) map[string]func(gorm.ColumnType) (dataType string) {
	dataMap := map[string]func(gorm.ColumnType) (dataType string){
		"numeric": func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "int32") },
		"integer": func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "int32") },
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			ct, _ := columnType.ColumnType()
			if strings.HasPrefix(ct, "tinyint(1)") {
				if columnType.Name() == "is_deleted" {
					return "soft_delete.DeletedAt"
				}
				return getDataType(cfg, columnType, "bool")
			}
			return getDataType(cfg, columnType, "int8")
		},
		"smallint": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "int16")
		},
		"mediumint": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "int32")
		},
		"bigint": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "int64")
		},
		"int": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "int")
		},
		"float": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "float64")
		},
		"real":       func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "float64") },
		"double":     func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "float64") },
		"decimal":    func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "float64") },
		"char":       func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"varchar":    func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"tinytext":   func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"mediumtext": func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"longtext":   func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"binary":     func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "[]byte") },
		"varbinary":  func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "[]byte") },
		"tinyblob":   func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "[]byte") },
		"blob":       func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "[]byte") },
		"mediumblob": func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "[]byte") },
		"longblob":   func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "[]byte") },
		"text":       func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"json": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "datatypes.JSON")
		},
		"enum": func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"time": func(columnType gorm.ColumnType) (dataType string) {
			if columnType.Name() == "deleted_at" {
				return getDataType(cfg, columnType, "gorm.DeletedAt")
			}
			return getDataType(cfg, columnType, "time.Time")
		},
		"date": func(columnType gorm.ColumnType) (dataType string) {
			if columnType.Name() == "deleted_at" {
				return getDataType(cfg, columnType, "gorm.DeletedAt")
			}
			return getDataType(cfg, columnType, "time.Time")
		},
		"datetime": func(columnType gorm.ColumnType) (dataType string) {
			if columnType.Name() == "deleted_at" {
				return getDataType(cfg, columnType, "gorm.DeletedAt")
			}
			return getDataType(cfg, columnType, "time.Time")
		},
		"timestamp": func(columnType gorm.ColumnType) (dataType string) {
			if columnType.Name() == "deleted_at" {
				return getDataType(cfg, columnType, "gorm.DeletedAt")
			}
			return getDataType(cfg, columnType, "time.Time")
		},
		"year":    func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "int32") },
		"bit":     func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "[]uint8") },
		"boolean": func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "bool") },
	}
	for k, v := range customMap {
		dataMap[k] = v
	}
	return dataMap
}

func GetDataMapPostgreSQL(cfg *gen.Config, customMap map[string]func(gorm.ColumnType) (
	dataType string)) map[string]func(gorm.ColumnType) (dataType string) {
	dataMap := map[string]func(gorm.ColumnType) (dataType string){
		"int2": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "int16")
		},
		"int4": func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "int32") },
		"int8": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "int64")
		},
		"float4": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "float64")
		},
		"float8": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType,
				"float64")
		},
		"numeric": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType,
				"float64")
		},
		"bpchar":  func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"varchar": func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"text":    func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"bytea": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType,
				"[]byte")
		},
		"json": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "datatypes.JSON")
		},
		"jsonb": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "datatypes.JSON")
		},
		"enum": func(columnType gorm.ColumnType) (dataType string) { return getDataType(cfg, columnType, "string") },
		"time": func(columnType gorm.ColumnType) (dataType string) {
			if columnType.Name() == "deleted_at" {
				return "gorm.DeletedAt"
			}
			return getDataType(cfg, columnType, "time.Time")
		},
		"timez": func(columnType gorm.ColumnType) (dataType string) {
			if columnType.Name() == "deleted_at" {
				return "gorm.DeletedAt"
			}
			return getDataType(cfg, columnType, "time.Time")
		},
		"timestamp": func(columnType gorm.ColumnType) (dataType string) {
			if columnType.Name() == "deleted_at" {
				return "gorm.DeletedAt"
			}
			return getDataType(cfg, columnType, "time.Time")
		},
		"timestampz": func(columnType gorm.ColumnType) (dataType string) {
			if columnType.Name() == "deleted_at" {
				return "gorm.DeletedAt"
			}
			return getDataType(cfg, columnType, "time.Time")
		},
		"date": func(columnType gorm.ColumnType) (dataType string) {
			if columnType.Name() == "deleted_at" {
				return "gorm.DeletedAt"
			}
			return getDataType(cfg, columnType, "time.Time")
		},
		"bool": func(columnType gorm.ColumnType) (dataType string) {
			if columnType.Name() == "is_deleted" {
				return "soft_delete.DeletedAt"
			}
			return getDataType(cfg, columnType, "bool")
		},
		"uuid": func(columnType gorm.ColumnType) (dataType string) {
			return getDataType(cfg, columnType, "string")
		},
	}
	for k, v := range customMap {
		dataMap[k] = v
	}
	return dataMap
}

func getDataType(cfg *gen.Config, columnType gorm.ColumnType, targetType string) (dataType string) {
	ct, _ := columnType.ColumnType()
	if cfg.FieldSignable && strings.HasPrefix(targetType, "int") && strings.HasSuffix(ct, "unsigned") {
		targetType = strings.Replace(targetType, "int", "uint", 1)
	}
	if cfg.FieldNullable {
		return TypeNullable(columnType, targetType)
	}
	return targetType
}
