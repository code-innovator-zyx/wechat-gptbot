package handler

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"time"
	"wechat-gptbot/utils"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/4/19 10:41
* @Package:
 */

func KeepAlive(bot *openwechat.Self) {

	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			heartBeat(bot)
		}
	}
}
func heartBeat(bot *openwechat.Self) {
	// 向文件传输助手发送消息，不要再关注公众号了
	// 生成要发送的消息
	outMessage := fmt.Sprintf("防微信自动退出登录[%d]", utils.GetRandInt64(1000))
	bot.SendTextToFriend(openwechat.NewFriendHelper(bot), outMessage)
}
