package dao

import (
	"gorm.io/gorm"
)

type Table interface {
	TableName(svc string) string
	PgCreate(db *gorm.DB, svc string) error
}
type TableTooler interface {
	NewTable(db *gorm.DB, t Table, svc string) error
	CheckTableExist(db *gorm.DB, t Table, svc string) bool
}
