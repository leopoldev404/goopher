package services

import "awesome/types"

type userService struct {
	db types.Db
}

func NewUserService(db types.Db) *userService {
	return &userService{db: db}
}

func (us *userService) Save(user *types.User) {

}
