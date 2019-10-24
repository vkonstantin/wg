package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/vkonstantin/wg/todo/common/auth"
	"github.com/vkonstantin/wg/todo/controller"
	"github.com/vkonstantin/wg/todo/message"
	"github.com/vkonstantin/wg/todo/server"
	"log"
	"net/http"
)

const (
	authTokenHeader = "token"
)

func New(listen string, controller controller.MainService) server.Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	s := rest{
		listen:     listen,
		engine:     engine,
		controller: controller,
	}
	s.initHandlers()
	return &s
}

type rest struct {
	listen     string
	engine     *gin.Engine
	controller controller.MainService
}

func (r rest) Run() error {
	log.Printf("start listen on %s", r.listen)
	return r.engine.Run(r.listen)
}

func (r rest) handlerNoAuth(req interface{}, fnc controller.ActionNoAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.BindJSON(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, message.Error{Message: err.Error()})
			return
		}

		resp, er := fnc(req)
		if er != nil {
			c.JSON(er.HttpCode, message.Error{Message: er.Message})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func (r rest) handler(req interface{}, fnc controller.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader(authTokenHeader)
		token, err := auth.NewTokenFromString(tokenStr)
		if err != nil {
			log.Printf("invalid auth token: %s", tokenStr)
			c.Data(http.StatusUnauthorized, "", nil)
			return
		}

		if req != nil {
			err = c.BindJSON(req)
			if err != nil {
				c.JSON(http.StatusBadRequest, message.Error{Message: err.Error()})
				return
			}
		}

		user := token.User()
		resp, er := fnc(user, req)
		if er != nil {
			c.JSON(er.HttpCode, message.Error{Message: er.Message})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
