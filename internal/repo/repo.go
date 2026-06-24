// Package repo defines data access interfaces and containers.
package repo

import "gorm.io/gorm"

type Repository struct {
	User UserRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{User: NewUserRepo(db)}
}
