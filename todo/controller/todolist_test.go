package controller

import (
	"github.com/stretchr/testify/assert"
	"github.com/vkonstantin/wg/todo/common/auth"
	"github.com/vkonstantin/wg/todo/message"
	"github.com/vkonstantin/wg/todo/model"
	"testing"
)

func TestAddTODO(t *testing.T) {
	s := newTodoListService(getStorage())
	u := &auth.User{ID: 1}
	response, err := s.AddTODO(u, &message.AddTodoRequest{RequestID: "r1", Text: "text1"})
	resp1 := response.(*message.AddTodoResponse)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), resp1.ID)

	resp, err := s.ListOfTODOs(u, nil)
	list := resp.([]*model.Item)
	assert.Equal(t, 1, len(list))

	item2resp, _ := s.AddTODO(u, &message.AddTodoRequest{RequestID: "r2", Text: "text2"})
	item2 := item2resp.(*message.AddTodoResponse)
	s.AddTODO(u, &message.AddTodoRequest{RequestID: "r3", Text: "text3"})

	resp, err = s.ListOfTODOs(u, nil)
	list = resp.([]*model.Item)
	assert.Equal(t, 3, len(list))

	s.ResolveTODO(u, &message.ResolveTodoRequest{RequestID: "r4", ID: item2.ID})
	resp, err = s.ListOfTODOs(u, nil)
	list = resp.([]*model.Item)
	assert.Equal(t, 2, len(list))
}
