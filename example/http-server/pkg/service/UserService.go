package service

import (
	"errors"
	"fmt"
	"github.com/zhangliang-zl/reskit/example/http-server/pkg/model"
)

type UserService struct {
}

func (*UserService) Info(id string) (model.UserModel, error) {
	var user model.UserModel
	if id != "0" {
		user = model.UserModel{
			Id:   id,
			Name: fmt.Sprintf("name_%s", id),
		}
		return user, nil
	}

	return user, errors.New("not found this record")
}

func NewUserService() *UserService {
	return &UserService{}
}
