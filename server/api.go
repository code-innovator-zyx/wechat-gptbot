package server

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"net/http"
	routes "wechat-gptbot/server/routers"
)

/*
* @Author: zouyx
* @Email: 开放的api接口
* @Date:   2024/5/28 16:03
* @Package:
 */
type ApiServer struct {
	port int
	bot  *openwechat.Bot
}

// options 熔断器配置参数.
type option struct {
	port int
}

type Option func(*option)

func NewApiServer(self *openwechat.Bot, opts ...Option) ApiServer {
	opt := option{
		port: 8502,
	}
	for i := range opts {
		opts[i](&opt)
	}

	return ApiServer{
		port: opt.port,
		bot:  self,
	}
}

func WithPort(p int) Option {
	return func(c *option) {
		c.port = p
	}
}

func (server ApiServer) Run() {
	// 关闭debug模式
	gin.SetMode(gin.ReleaseMode)
	router := routes.InitRoute(server.bot)
	http.ListenAndServe(fmt.Sprintf(":%d", server.port), router)
}
