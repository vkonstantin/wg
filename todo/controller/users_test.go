package controller

import (
	"github.com/stretchr/testify/assert"
	"github.com/vkonstantin/wg/todo/message"
	"github.com/vkonstantin/wg/todo/storage"
	"github.com/vkonstantin/wg/todo/storage/memory"
	"testing"
	"time"
)

func TestAddUser(t *testing.T) {
	s := newUsersService(getStorage())
	// user 1
	req := &message.AddUserRequest{RequestID: "r1"}
	response, err := s.AddUser(req)
	resp := response.(*message.AddUserResponse)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), resp.UserID)
	assert.NotEmpty(t, resp.Token)

	// user 2
	req = &message.AddUserRequest{RequestID: "r2"}
	response, err = s.AddUser(req)
	resp = response.(*message.AddUserResponse)
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), resp.UserID)

	// user 1 duplicate
	req = &message.AddUserRequest{RequestID: "r1"}
	response, err = s.AddUser(req)
	assert.NotNil(t, err)
	assert.Nil(t, response)

	// Add some users
	s.AddUser(&message.AddUserRequest{RequestID: "r3"})
	s.AddUser(&message.AddUserRequest{RequestID: "r4"})
	s.AddUser(&message.AddUserRequest{RequestID: "r5"})
	s.AddUser(&message.AddUserRequest{RequestID: "r6"})
	time.Sleep(time.Second * 2) // wait for GC

	// check user 1 duplicate again
	req = &message.AddUserRequest{RequestID: "r1"}
	response, err = s.AddUser(req)
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func getStorage() storage.Storage {
	conf := memory.Config{
		GarbageCollectionInterval: time.Second,
		DeduplicationCapacitySize: 2,
		GarbageThreshold:          0.5,
	}
	mem := memory.New(conf)
	return mem
}
