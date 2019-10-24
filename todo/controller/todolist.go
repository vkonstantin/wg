package controller

import (
	"github.com/vkonstantin/wg/todo/common/auth"
	"github.com/vkonstantin/wg/todo/message"
	"github.com/vkonstantin/wg/todo/storage"
	"log"
	"net/http"
)

// Action is a handler function of action that required authorization
type Action func(u *auth.User, req interface{}) (resp interface{}, err *message.Error)

// TodoListService interface of actions of TodoList controller
type TodoListService interface {
	AddTODO(u *auth.User, req interface{}) (resp interface{}, err *message.Error)
	ResolveTODO(u *auth.User, req interface{}) (resp interface{}, err *message.Error)
	ListOfTODOs(u *auth.User, req interface{}) (interface{}, *message.Error)
}

type todoListService struct {
	storage storage.Storage
}

func newTodoListService(s storage.Storage) *todoListService {
	ls := todoListService{storage: s}
	return &ls
}

func (s *todoListService) ListOfTODOs(u *auth.User, req interface{}) (interface{}, *message.Error) {
	list, err := s.storage.ListTODOs()
	if err != nil {
		log.Printf("error: %s, on storage.ListTODOs", err)
		return nil, &message.Error{HttpCode: http.StatusInternalServerError}
	}

	return list, nil
}

// AddTODO add todo item to the list
func (s *todoListService) AddTODO(u *auth.User, req interface{}) (interface{}, *message.Error) {
	r := req.(*message.AddTodoRequest)
	// that is safe to cast like this in this case because req structure was created exactly in this type in previous function
	if s.storage.IsDuplicate(r.RequestID) {
		return nil, &message.Error{HttpCode: http.StatusBadRequest, Message: "Duplicate"}
	}

	item, err := s.storage.AddTODO(u.ID, r.Text)
	if err != nil {
		log.Printf("error: %s, on storage.AddTODO", err)
		return nil, &message.Error{HttpCode: http.StatusInternalServerError}
	}

	resp := &message.AddTodoResponse{ID: item.ID}
	return resp, nil
}

// ResolveTODO remove todo item from the list
func (s *todoListService) ResolveTODO(u *auth.User, req interface{}) (interface{}, *message.Error) {
	r := req.(*message.ResolveTodoRequest)
	if s.storage.IsDuplicate(r.RequestID) {
		return nil, &message.Error{HttpCode: http.StatusBadRequest, Message: "Duplicate"}
	}

	err := s.storage.RemoveTODO(r.ID)
	if err != nil {
		log.Printf("error: %s, on storage.RemoveTODO", err)
		return nil, &message.Error{HttpCode: http.StatusInternalServerError}
	}

	resp := &message.ResolveTodoResponse{Status: "Done"}
	return resp, nil
}
