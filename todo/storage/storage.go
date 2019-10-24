package storage

import (
	"github.com/vkonstantin/wg/todo/common/auth"
	"github.com/vkonstantin/wg/todo/model"
)

// Storage interface
type Storage interface {
	IsDuplicate(requestID string) bool
	AddUser() (*auth.User, error)
	AddTODO(userID uint64, text string) (*model.Item, error)
	ListTODOs() ([]*model.Item, error)
	RemoveTODO(itemID uint64) error
}
