package rest

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/vkonstantin/wg/todo/message"
	"net/http"
	"testing"
)

func TestAddUser(t *testing.T) {
	s := newTestService()
	code, body, _ := s.request("POST", "/user", &message.AddUserRequest{RequestID: randRequestID()})
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, `{"userID":1,"token":"{\"id\":1}"}`, body)
}

func TestTodo(t *testing.T) {
	s := newTestService()
	_, body, _ := s.request("POST", "/user", &message.AddUserRequest{RequestID: randRequestID()})
	u := new(message.AddUserResponse)
	err := json.Unmarshal([]byte(body), u)
	assert.Nil(t, err)
	assert.NotEmpty(t, u.Token)

	// try to add item
	req := &message.AddTodoRequest{RequestID: randRequestID(), Text: "TODO_1"}
	code, body, _ := s.request("POST", "/todo", req)
	assert.Equal(t, http.StatusUnauthorized, code)

	// try with token
	auth := []string{authTokenHeader, u.Token}
	code, body, _ = s.request("POST", "/todo", req, auth...)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, `{"id":1}`, body)

	// list
	code, body, _ = s.request("GET", "/todo", req, auth...)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, `[{"id":1,"text":"TODO_1"}]`, body)

	// remove
	reqRem := &message.ResolveTodoRequest{RequestID: randRequestID(), ID: 1}
	code, body, _ = s.request("POST", "/todo/resolve", reqRem, auth...)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, `{"status":"Done"}`, body)

	// list
	code, body, _ = s.request("GET", "/todo", req, auth...)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, `[]`, body)
}
