package controller

import (
	"github.com/vkonstantin/wg/todo/common/auth"
	"github.com/vkonstantin/wg/todo/message"
	"github.com/vkonstantin/wg/todo/storage"
	"log"
	"net/http"
)

// ActionNoAuth is a handler function of action that no need authorization
type ActionNoAuth func(req interface{}) (resp interface{}, err *message.Error)

// UsersService interface of actions of user controller
type UsersService interface {
	AddUser(req interface{}) (resp interface{}, err *message.Error)
}

type usersService struct {
	storage storage.Storage
}

func newUsersService(s storage.Storage) *usersService {
	us := usersService{storage: s}
	return &us
}

// AddUser create new user
func (s *usersService) AddUser(req interface{}) (interface{}, *message.Error) {
	r := req.(*message.AddUserRequest)
	// that is safe to cast like this in this case because req structure was created exactly in this type in previous function
	if s.storage.IsDuplicate(r.RequestID) {
		return nil, &message.Error{HttpCode: http.StatusBadRequest, Message: "Duplicate"}
	}

	user, err := s.storage.AddUser()
	if err != nil {
		log.Printf("error: %s, on storage.AddUser", err)
		return nil, &message.Error{HttpCode: http.StatusInternalServerError}
	}
	token := auth.NewToken(user)
	tokenStr, err := token.String()
	if err != nil {
		log.Printf("error: %s, on token.String", err)
		return nil, &message.Error{HttpCode: http.StatusInternalServerError}
	}
	return &message.AddUserResponse{UserID: user.ID, Token: tokenStr}, nil
}
