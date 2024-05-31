package controller

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2024/5/28 17:31
* @Package:
 */

type CheckLoginResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		QrUrl     string `json:"qr_url"`
		UserName  string `json:"user_name"`
		AvatarUrl string `json:"avatar_url"`
	} `json:"data"`
}

// CheckLogin 校验微信登录状态
func CheckLogin(c *gin.Context, bot *openwechat.Bot) {
	// 获取当前用户
	user, err := bot.GetCurrentUser()
	response := CheckLoginResponse{
		Code: http.StatusOK,
		Msg:  "ok",
	}
	if nil != err {
		response.Code = http.StatusNetworkAuthenticationRequired
		response.Msg = err.Error()
		response.Data.QrUrl = openwechat.GetQrcodeUrl(bot.UUID())
	} else {
		// 已登录状态 返回用户名
		response.Data.UserName = user.NickName
		response.Data.AvatarUrl = bot.Caller.Client.Domain.BaseHost() + user.HeadImgUrl
	}
	c.JSON(http.StatusOK, response)
}
