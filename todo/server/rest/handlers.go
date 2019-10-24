package rest

import "github.com/vkonstantin/wg/todo/message"

func (r rest) initHandlers() {
	r.engine.POST("/user", r.handlerNoAuth(new(message.AddUserRequest), r.controller.AddUser))

	todoGroup := r.engine.Group("/todo")
	todoGroup.POST("", r.handler(new(message.AddTodoRequest), r.controller.AddTODO))
	todoGroup.POST("/resolve", r.handler(new(message.ResolveTodoRequest), r.controller.ResolveTODO))
	todoGroup.GET("", r.handler(nil, r.controller.ListOfTODOs))
}
