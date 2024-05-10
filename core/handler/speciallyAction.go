package handler

import (
	"github.com/eatmoreapple/openwechat"
	"time"
	"wechat-gptbot/config"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/4/19 10:41
* @Package:
 */

func KeepAlive(bot *openwechat.Self) {

	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			heartBeat(bot)
		}
	}
}
func heartBeat(bot *openwechat.Self) {
	// 获取公众号
	if mps, _ := bot.Mps(false); mps != nil {
		for i := range mps {
			if mps[i].NickName == config.C.KeepaliveRobot {
				mps[i].SendText("ping")
				return
			}
		}
	}
}
