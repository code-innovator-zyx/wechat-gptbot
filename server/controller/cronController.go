package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"wechat-gptbot/config"
	"wechat-gptbot/core/handler"
	"wechat-gptbot/core/plugins"
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
	if config.C.Cron.WeatherConfig.Desc == "" && config.C.Cron.WeatherConfig.Spec != "" {
		config.C.ResetCron(func(cfg *config.CronConfig) {
			cfg.WeatherConfig.Desc = handler.Context.Session.DescribeQuartzCron(config.C.Cron.WeatherConfig.Spec)
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"cron":  config.C.Cron.WeatherConfig.Desc,
		"users": config.C.Cron.WeatherConfig.Users,
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
	config.C.ResetCron(func(cfg *config.CronConfig) {
		for i, user := range cfg.WeatherConfig.Users {
			if user.Name == name {
				// 删除用户
				(cfg.WeatherConfig.Users)[i] = (cfg.WeatherConfig.Users)[len(cfg.WeatherConfig.Users)-1]
				cfg.WeatherConfig.Users = (cfg.WeatherConfig.Users)[:len(cfg.WeatherConfig.Users)-1]
				return
			}
		}
	})
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
	config.C.ResetCron(func(cfg *config.CronConfig) {
		cfg.WeatherConfig.Users = append(cfg.WeatherConfig.Users, struct {
			Name string `json:"name"`
			City string `json:"city"`
		}{Name: user.Name, City: user.City})
	})
	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

// GetNewsSetting 热点新闻定时任务相关配置
func GetNewsSetting(ctx *gin.Context) {
	if config.C.Cron.NewsConfig.Desc == "" && config.C.Cron.NewsConfig.Spec != "" {
		config.C.ResetCron(func(cfg *config.CronConfig) {
			cfg.NewsConfig.Desc = handler.Context.Session.DescribeQuartzCron(config.C.Cron.NewsConfig.Spec)
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"users":      config.C.Cron.NewsConfig.Users,
		"groups":     config.C.Cron.NewsConfig.Groups,
		"cron":       config.C.Cron.NewsConfig.Desc,
		"rss_source": config.C.Cron.NewsConfig.RssSource,
		"top_n":      config.C.Cron.NewsConfig.TopN,
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
	config.C.ResetCron(func(cfg *config.CronConfig) {
		cfg.NewsConfig.Users = user.Users
		cfg.NewsConfig.Groups = user.Groups
	})

	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

// GetSportSetting 获取微信运动配置
func GetSportSetting(ctx *gin.Context) {
	if config.C.Cron.SportConfig.Desc == "" && config.C.Cron.SportConfig.Spec != "" {
		config.C.ResetCron(func(cfg *config.CronConfig) {
			cfg.SportConfig.Desc = handler.Context.Session.DescribeQuartzCron(config.C.Cron.SportConfig.Spec)
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"users": config.C.Cron.SportConfig.Users,
		"cron":  config.C.Cron.SportConfig.Desc})
}

func AddSportReceiver(ctx *gin.Context) {
	var user config.SportAccount
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON input"})
		return
	}

	config.C.ResetCron(func(cfg *config.CronConfig) {
		cfg.SportConfig.Users = append(cfg.SportConfig.Users, user)
	})

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
	config.C.ResetCron(func(cfg *config.CronConfig) {
		for i, user := range cfg.SportConfig.Users {
			if user.Name == name {
				// 删除用户
				(cfg.SportConfig.Users)[i] = (cfg.SportConfig.Users)[len(cfg.SportConfig.Users)-1]
				cfg.SportConfig.Users = (cfg.SportConfig.Users)[:len(cfg.SportConfig.Users)-1]
				return
			}
		}
	})
	// 如果没有找到用户，返回Not Found
	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
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
	var fun func(cfg *config.CronConfig)
	switch req.PluginName {
	case weather2.WeatherPluginName:
		fun = func(cfg *config.CronConfig) {
			cfg.WeatherConfig.Desc = req.Desc
		}
		spec = &config.C.Cron.WeatherConfig.Spec
	case news2.NewsPluginName:
		fun = func(cfg *config.CronConfig) {
			cfg.NewsConfig.Desc = req.Desc

		}
		spec = &config.C.Cron.NewsConfig.Spec
	case wechatMovement.StepPluginName:
		fun = func(cfg *config.CronConfig) {
			cfg.SportConfig.Desc = req.Desc
		}
		spec = &config.C.Cron.SportConfig.Spec

	default:
		c.JSON(http.StatusBadRequest, gin.H{"msg": "unknown pluginName"})
		return
	}
	// 修改配置
	config.C.ResetCron(fun)
	// 根据描述生成 表达式
	cron := handler.Context.Session.GenerateQuartzCron(req.Desc)
	logrus.Infof("%s   生成的 cron 表达式 %s", req.Desc, cron)
	*spec = cron
	// 重置 定时器
	handler.Context.CronServer.ResetPluginCron(req.PluginName, *spec)
	c.JSON(http.StatusBadRequest, gin.H{"msg": "ok"})
}

type ResetRssRequest struct {
	Source string `json:"source"`
	TopN   int    `json:"top_n"` // 时间描述
}

func ResetRss(c *gin.Context) {
	var request ResetRssRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		// 如果查询参数中没有提供name，则返回Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing 'source' query parameter"})
		return
	}
	request.Source = strings.TrimSpace(request.Source)
	// 关闭 rss 无需校验
	if request.Source != "" {
		fp := gofeed.NewParser()
		_, err = fp.ParseURL(request.Source)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"msg": "无效的RSS源"})
			return
		}
	}
	config.C.ResetCron(func(cfg *config.CronConfig) {
		cfg.NewsConfig.RssSource = request.Source
		cfg.NewsConfig.TopN = request.TopN
	})

	// 重置插件
	err = plugins.Manger.ResetPlugin(news2.NewsPluginName, func(_ plugins.PluginSvr) plugins.PluginSvr {
		return news2.NewPlugin(news2.SetTopN(request.TopN), news2.SetRssSource(request.Source))
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"msg": "设置失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
