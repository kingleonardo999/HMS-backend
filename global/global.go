package global

import (
	"gorm.io/gorm"
)

var (
	Db *gorm.DB // Global variable for the database connection
)
