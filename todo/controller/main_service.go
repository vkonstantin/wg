package controller

import (
	"github.com/vkonstantin/wg/todo/storage"
)

// MainService is a composition of all controllers of this service
type MainService interface {
	UsersService
	TodoListService
}

type service struct {
	*usersService
	*todoListService
}

// NewMainService constructor
func NewMainService(s storage.Storage) MainService {
	srv := service{
		usersService:    newUsersService(s),
		todoListService: newTodoListService(s),
	}

	return &srv
}
