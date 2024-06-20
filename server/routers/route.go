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
		r.GET("/friends", func(context *gin.Context) {
			controller.GetFriends(context, bot)
		})
		r.POST("/cron-reset", controller.ResetCron)
	}
	{
		weather := r.Group("weather")
		weather.GET("/cron-setting", controller.GetWeatherSetting)
		weather.DELETE("/receiver", controller.DeleteWeatherReceiver)
		weather.POST("/receiver", controller.AddWeatherReceiver)
	}
	{
		news := r.Group("news")
		news.GET("/cron-setting", controller.GetNewsSetting)
		news.POST("/receiver", controller.ResetNewsReceiver)
	}
	{
		sport := r.Group("sport")
		sport.GET("/cron-setting", controller.GetSportSetting)
		sport.DELETE("/receiver", controller.DeleteSportReceiver)
		sport.POST("/receiver", controller.AddSportReceiver)
	}
	return router
}
