package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wechat-gptbot/config"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/6/18 17:05
* @Package: 定时任务相关接口
 */

// GetWeatherSetting 天气预报定时任务相关配置
func GetWeatherSetting(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
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
	for i := range *users {
		if (*users)[i].Name == user.Name {
			// 更新用户的城市信息
			(*users)[i].City = user.City
			ctx.JSON(http.StatusOK, gin.H{"msg": "updated"})
			return
		}
	}

	// 添加新用户
	*users = append(*users, struct {
		Name string `json:"name"`
		City string `json:"city"`
	}{Name: user.Name, City: user.City})

	ctx.JSON(http.StatusOK, gin.H{"msg": "added"})
}
