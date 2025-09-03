package ezgen

import (
	"reflect"

	"gorm.io/gen"
	"gorm.io/gorm"
)

func WithDeletedList(withDeleted []bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(withDeleted) == 0 || !withDeleted[0] {
			return db
		} else {
			return db.Unscoped()
		}
	}
}

func WithDeleted(withDeleted bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if withDeleted {
			return db.Unscoped()
		} else {
			return db
		}
	}
}

func Paginate(p Pager) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p == nil || reflect.ValueOf(p).IsNil() || p.GetPageSize() <= 0 || p.GetPageIndex() <= 0 {
			return db
		}
		return db.Offset(int((p.GetPageIndex() - 1) * p.GetPageSize())).Limit(int(p.GetPageSize()))
	}
}

func Cond(cond bool, query any, args ...any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if cond {
			return db.Where(query, args...)
		}
		return db
	}
}

func Nullable(cond bool, query any, f func() any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if cond {
			return db.Where(query, f())
		}
		return db
	}
}

func WithDeletedListGen(withDeleted []bool) func(db gen.Dao) gen.Dao {
	return func(db gen.Dao) gen.Dao {
		if len(withDeleted) == 0 || !withDeleted[0] {
			return db
		} else {
			return db.Unscoped()
		}
	}
}

func WithDeletedGen(withDeleted bool) func(db gen.Dao) gen.Dao {
	return func(db gen.Dao) gen.Dao {
		if withDeleted {
			return db.Unscoped()
		} else {
			return db
		}
	}
}

func PaginateGen(p Pager) func(db gen.Dao) gen.Dao {
	return func(db gen.Dao) gen.Dao {
		if p == nil || reflect.ValueOf(p).IsNil() || p.GetPageSize() <= 0 || p.GetPageIndex() <= 0 {
			return db
		}
		return db.Offset(int((p.GetPageIndex() - 1) * p.GetPageSize())).Limit(int(p.GetPageSize()))
	}
}

func CondGen(cond bool, conds ...gen.Condition) func(db gen.Dao) gen.Dao {
	return func(db gen.Dao) gen.Dao {
		if cond {
			return db.Where(conds...)
		}
		return db
	}
}

func NullableGen(cond bool, f func() []gen.Condition) func(db gen.Dao) gen.Dao {
	return func(db gen.Dao) gen.Dao {
		if cond {
			return db.Where(f()...)
		}
		return db
	}
}
