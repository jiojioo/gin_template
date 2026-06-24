package mysql

import "github.com/jiojioo/gin_template/internal/model"

func AutoMigrate() error {
	return Client.AutoMigrate(&model.User{})
}
