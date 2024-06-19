package controller

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/6/18 16:08
* @Package: 微信好友
 */

type FriendsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Users  []string `json:"users"`
		Groups []string `json:"groups"`
	} `json:"data"`
}

// GetFriends 获取当前微信所有朋友群组关系
func GetFriends(c *gin.Context, bot *openwechat.Bot) {
	response := FriendsResponse{
		Code: http.StatusOK,
		Msg:  "ok",
	}
	// 获取当前用户
	user, err := bot.GetCurrentUser()
	// 还没有登录
	if nil != err {
		// 如果未登录，返回错误信息
		c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{
			"code": http.StatusNetworkAuthenticationRequired,
			"msg":  "User not authenticated",
		})
		return
	}
	// 获取当前用户朋友
	friends, err := user.Friends(true)
	if err != nil {
		// 获取朋友列表失败，返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed to get friends",
		})
		return
	}
	// 收集昵称
	for i := range friends {
		response.Data.Users = append(response.Data.Users, friends[i].NickName)
	}
	// 获取群聊
	groups, err := user.Groups(true)
	if err != nil {
		// 获取群聊列表失败，返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed to get groups",
		})
		return
	}
	for i := range groups {
		response.Data.Groups = append(response.Data.Groups, groups[i].NickName)
	}

	c.JSON(http.StatusOK, response)

}
