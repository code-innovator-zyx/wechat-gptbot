package cron

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/plugins"
	"wechat-gptbot/core/plugins/weather"
)

/*
* @Author: zouyx
* @Email:1003941268@qq.com
* @Date:   2024/6/18 10:04
* @Package:
 */

type Weather struct {
	pluginSvr plugins.PluginSvr
	bot       *openwechat.Self
	cfg       *config.WeatherCronConfig
}

func NewWeatherCron() *Weather {
	return &Weather{
		pluginSvr: plugins.Manger.GetPluginSvr(weather.WeatherPluginName),
	}
}
func (w *Weather) WithBot(bot *openwechat.Self) *Weather {
	w.bot = bot
	return w
}
func (w *Weather) WithCfg(cfg *config.WeatherCronConfig) *Weather {
	w.cfg = cfg
	return w
}
func (w *Weather) Run() {
	if len(w.cfg.Users) == 0 {
		// 没有需要通知的用户
		return
	}

	friends, err := w.bot.Friends()
	if nil != err {
		logrus.Errorf("failed get friends iter ,err =%s", err.Error())
		return
	}
	// 这里通过城市进行分组，每组用户使用转发的形式进行发送，避免被风控
	groups := make(map[string]openwechat.Friends)
	for i := range w.cfg.Users {
		result := friends.SearchByNickName(1, w.cfg.Users[i].Name)
		if len(result) == 0 {
			// 使用备注查询
			result = friends.SearchByRemarkName(1, w.cfg.Users[i].Name)
		}
		if len(result) == 0 {
			logrus.Errorf("没有找到昵称或者备注是 %s 的用户", w.cfg.Users[i].Name)
			continue
		}
		groups[w.cfg.Users[i].City] = append(groups[w.cfg.Users[i].City], result...)
	}
	// 进行群发消息
	for city, users := range groups {
		// 发送消息
		msg, err := w.bot.SendTextToFriend(users[0], w.pluginSvr.Do(city)[0])
		if err != nil {
			fmt.Printf("发送给 %s失败 %s\n", users[0].User.NickName, err.Error())
			continue
		}
		// 发送成功以后，进行消息转发给其他用户
		if len(users[1:]) > 0 {
			w.bot.ForwardMessageToFriends(msg, time.Second*1, users[1:]...)
		}
	}
}
