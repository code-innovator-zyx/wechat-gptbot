package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"wechat-gptbot/config"
	"wechat-gptbot/core/handler"
	news2 "wechat-gptbot/core/plugins/news"
	weather2 "wechat-gptbot/core/plugins/weather"
	"wechat-gptbot/core/plugins/wechatMovement"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/6/18 17:05
* @Package: 定时任务相关接口
 */

// GetWeatherSetting 天气预报定时任务相关配置
func GetWeatherSetting(ctx *gin.Context) {
	if config.C.CronConfig.WeatherConfig.Desc == "" {
		config.C.CronConfig.WeatherConfig.Desc = handler.Context.Session.DescribeQuartzCron(config.C.CronConfig.WeatherConfig.Spec)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"cron":  config.C.CronConfig.WeatherConfig.Desc,
		"users": config.C.CronConfig.WeatherConfig.Users,
	})
}

// DeleteWeatherReceiver 删除天气预报接受者
func DeleteWeatherReceiver(ctx *gin.Context) {
	name, ok := ctx.GetQuery("name")
	if !ok {
		// 如果查询参数中没有提供name，则返回Bad Request
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Missing 'name' query parameter"})
		return
	}
	config.C.Lock()
	defer config.C.Unlock()
	users := &config.C.CronConfig.WeatherConfig.Users
	for i, user := range *users {
		if user.Name == name {
			// 删除用户
			(*users)[i] = (*users)[len(*users)-1]
			*users = (*users)[:len(*users)-1]
			ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
			return
		}
	}
	// 如果没有找到用户，返回Not Found
	ctx.JSON(http.StatusNotFound, gin.H{"msg": "User not found"})
}

type BindUser struct {
	Name string `json:"name" binding:"required"`
	City string `json:"city" binding:"required"`
}

// AddWeatherReceiver 添加天气预报接受者
func AddWeatherReceiver(ctx *gin.Context) {
	var user BindUser
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON input"})
		return
	}

	config.C.Lock()
	defer config.C.Unlock()
	users := &config.C.CronConfig.WeatherConfig.Users

	// 添加新用户
	*users = append(*users, struct {
		Name string `json:"name"`
		City string `json:"city"`
	}{Name: user.Name, City: user.City})

	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

// GetNewsSetting 热点新闻定时任务相关配置
func GetNewsSetting(ctx *gin.Context) {
	if config.C.CronConfig.NewsConfig.Desc == "" {
		config.C.CronConfig.NewsConfig.Desc = handler.Context.Session.DescribeQuartzCron(config.C.CronConfig.NewsConfig.Spec)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"users":  config.C.CronConfig.NewsConfig.Users,
		"groups": config.C.CronConfig.NewsConfig.Groups,
		"cron":   config.C.CronConfig.NewsConfig.Desc,
	})
}

type NewsUser struct {
	Users  []string `json:"users"`
	Groups []string `json:"groups" `
}

// ResetNewsReceiver 重置热点接收
func ResetNewsReceiver(ctx *gin.Context) {
	var user NewsUser
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON input"})
		return
	}
	config.C.Lock()
	defer config.C.Unlock()
	newsConf := &config.C.CronConfig.NewsConfig
	newsConf.Users = user.Users
	newsConf.Groups = user.Groups

	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

// GetSportSetting 获取微信运动配置
func GetSportSetting(ctx *gin.Context) {
	if config.C.CronConfig.SportConfig.Desc == "" && config.C.CronConfig.SportConfig.Spec != "" {
		config.C.CronConfig.SportConfig.Desc = handler.Context.Session.DescribeQuartzCron(config.C.CronConfig.SportConfig.Spec)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"users": config.C.CronConfig.SportConfig.Users,
		"cron":  config.C.CronConfig.SportConfig.Desc})
}

func AddSportReceiver(ctx *gin.Context) {
	var user config.SportAccount
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON input"})
		return
	}

	config.C.Lock()
	defer config.C.Unlock()
	users := &config.C.CronConfig.SportConfig.Users

	// 添加新用户
	*users = append(*users, user)

	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

// DeleteSportReceiver 删除天气预报接受者
func DeleteSportReceiver(ctx *gin.Context) {
	name, ok := ctx.GetQuery("name")
	if !ok {
		// 如果查询参数中没有提供name，则返回Bad Request
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Missing 'name' query parameter"})
		return
	}
	config.C.Lock()
	defer config.C.Unlock()
	users := &config.C.CronConfig.SportConfig.Users
	for i, user := range *users {
		if user.Name == name {
			// 删除用户
			(*users)[i] = (*users)[len(*users)-1]
			*users = (*users)[:len(*users)-1]
			ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
			return
		}
	}
	// 如果没有找到用户，返回Not Found
	ctx.JSON(http.StatusNotFound, gin.H{"msg": "User not found"})
}

type ResetCronRequest struct {
	PluginName string `json:"plugin_name"`
	Desc       string `json:"desc"` // 时间描述
}

func ResetCron(c *gin.Context) {
	var req ResetCronRequest
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON input"})
		return
	}
	var spec *string
	// 修改配置
	switch req.PluginName {
	case weather2.WeatherPluginName:
		config.C.CronConfig.WeatherConfig.Desc = req.Desc
		spec = &config.C.CronConfig.WeatherConfig.Spec
	case news2.NewsPluginName:
		config.C.CronConfig.NewsConfig.Desc = req.Desc
		spec = &config.C.CronConfig.NewsConfig.Spec

	case wechatMovement.StepPluginName:
		config.C.CronConfig.SportConfig.Desc = req.Desc
		spec = &config.C.CronConfig.SportConfig.Spec

	default:
		c.JSON(http.StatusBadRequest, gin.H{"msg": "unknown pluginName"})
		return
	}
	// 根据描述生成 表达式
	cron := handler.Context.Session.GenerateQuartzCron(req.Desc)
	logrus.Infof("%s   生成的 cron 表达式 %s", req.Desc, cron)
	*spec = cron
	// 重置 定时器
	handler.Context.CronServer.ResetPluginCron(req.PluginName, *spec)
	c.JSON(http.StatusBadRequest, gin.H{"msg": "ok"})
}
