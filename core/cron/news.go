package cron

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/plugins"
	"wechat-gptbot/core/plugins/news"
)

/*
* @Author: zouyx
* @Email:1003941268@qq.com
* @Date:   2024/6/18 10:04
* @Package:
 */

type News struct {
	pluginSvr plugins.PluginSvr
	bot       *openwechat.Self
	cfg       *config.NewsCronConfig
}

func NewNewsCron() *News {
	return &News{
		pluginSvr: plugins.Manger.GetPluginSvr(news.NewsPluginName),
	}
}
func (n *News) WithBot(bot *openwechat.Self) *News {
	n.bot = bot
	return n
}

func (n *News) WithCfg(cfg *config.NewsCronConfig) *News {
	n.cfg = cfg
	return n
}

func (n *News) Run() {
	if len(n.cfg.Users) == 0 && len(n.cfg.Groups) == 0 {
		// 没有需要通知的用户
		return
	}
	newsInfo := n.pluginSvr.Do()
	n.sendFriends(newsInfo)
	n.sendGroups(newsInfo)

	groups := make([]*openwechat.Group, 0, len(n.cfg.Groups))

	for i := range n.cfg.Groups {
		groups = append(groups, &openwechat.Group{&openwechat.User{UserName: n.cfg.Groups[i]}})
	}

	if len(groups) != 0 {
		n.bot.SendTextToGroups(newsInfo[0], time.Second*2, groups...)
	}
}

// 点对点私发
func (n *News) sendFriends(data []string) {
	if len(n.cfg.Users) == 0 {
		return
	}
	friends, err := n.bot.Friends()
	if nil != err {
		logrus.Errorf("failed get friends iter ,err =%s", err.Error())
		return
	}
	// 这里通过城市进行分组，每组用户使用转发的形式进行发送，避免被风控
	users := make([]*openwechat.Friend, 0, len(n.cfg.Users))
	for i := range n.cfg.Users {
		result := friends.SearchByNickName(1, n.cfg.Users[i])
		if len(result) == 0 {
			// 使用备注查询
			result = friends.SearchByRemarkName(1, n.cfg.Users[i])
		}
		if len(result) == 0 {
			logrus.Errorf("没有找到昵称或者备注是 %s 的用户", n.cfg.Users[i])
			continue
		}
		users = append(users, result...)
	}
	if len(users) == 0 {
		return
	}
	// 进行群发消息
	for i := range data {
		// 发送给一个用户
		msg, err := n.bot.SendTextToFriend(users[0], data[i])
		if err != nil {
			fmt.Printf("发送给 %s失败 %s\n", users[0].User.NickName, err.Error())
			continue
		}
		// 发送成功以后，进行消息转发给其他用户
		if len(users[1:]) > 0 {
			n.bot.ForwardMessageToFriends(msg, time.Second*1, users[1:]...)
		}
	}
}

// 群发
func (n *News) sendGroups(data []string) {
	if len(n.cfg.Groups) == 0 {
		return
	}
	grouops, err := n.bot.Groups()
	if nil != err {
		logrus.Errorf("failed get grouops iter ,err =%s", err.Error())
		return
	}
	users := make([]*openwechat.Group, 0, len(n.cfg.Users))
	for i := range n.cfg.Groups {
		result := grouops.SearchByNickName(1, n.cfg.Groups[i])
		if len(result) == 0 {
			logrus.Errorf("没有找到昵称是 %s 的群组", n.cfg.Groups[i])
			continue
		}
		users = append(users, result...)
	}
	if len(users) == 0 {
		return
	}
	// 进行群发消息
	for i := range data {
		// 发送给一个用户
		msg, err := n.bot.SendTextToGroup(users[0], data[i])
		if err != nil {
			fmt.Printf("发送给 %s失败 %s\n", users[0].User.NickName, err.Error())
			continue
		}
		// 发送成功以后，进行消息转发给其他群组
		if len(users[1:]) > 0 {
			n.bot.ForwardMessageToGroups(msg, time.Second*1, users[1:]...)
		}
	}
}
