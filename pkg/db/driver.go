package db

import (
	"gorm.io/gorm"
)

type Driver interface {
	Name() string
	Initialize(dsn string) error
	Connect(dsn string) (*gorm.DB, error)
}
