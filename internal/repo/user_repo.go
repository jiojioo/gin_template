package repo

import (
	"context"

	"github.com/jiojioo/gin_template/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(context.Context, uint64) (*model.User, error)
	FindByUsername(context.Context, string) (*model.User, error)
	Create(context.Context, *model.User) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
