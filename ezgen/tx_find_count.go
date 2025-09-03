package ezgen

import (
	"gorm.io/gen"
	"gorm.io/gorm"
)

func FindAndCountTransaction(db *gorm.DB, result interface{}) (int64, error) {
	var count int64
	if err := db.Find(result).Error; err != nil {
		return 0, err
	}

	if err := db.Model(result).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func FindAndCountTransactionGen(db gen.DO, result interface{}) (int64, error) {
	err := db.Scan(result)
	if err != nil {
		return 0, err
	}

	return db.Offset(-1).Limit(-1).Count()
}
