package routes

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"wechat-gptbot/server/controller"
)

func InitRoute(bot *openwechat.Bot) *gin.Engine {
	router := gin.New()
	r := router.Group("/wechat-gptbot")
	{
		r.GET("/checklogin", func(context *gin.Context) {
			controller.CheckLogin(context, bot)
		})
		r.GET("/current-model", controller.CurrentModel)
		r.POST("/reset-model", controller.ResetModel)
	}
	return router
}
