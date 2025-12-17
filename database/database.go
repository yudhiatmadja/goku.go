package database

import (
	"gorm.io/gorm"
)

// DB holds the global database connection pool.
// It is initialized in the bootstrap package.
var DB *gorm.DB