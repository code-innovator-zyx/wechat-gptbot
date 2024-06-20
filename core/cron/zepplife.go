package cron

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/plugins"
	"wechat-gptbot/core/plugins/wechatMovement"
	"wechat-gptbot/core/plugins/wechatMovement/zeepLife"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/6/19 17:04
* @Package:
 */

type Sport struct {
	pluginSvr plugins.PluginSvr
	bot       *openwechat.Self
	cfg       *config.WechatSportCronConfig
}

func NewSportCron() *Sport {
	return &Sport{
		pluginSvr: plugins.Manger.GetPluginSvr(wechatMovement.StepPluginName),
	}
}
func (w *Sport) WithBot(bot *openwechat.Self) *Sport {
	w.bot = bot
	return w
}
func (w *Sport) WithCfg(cfg *config.WechatSportCronConfig) *Sport {
	w.cfg = cfg
	return w
}
func (w *Sport) Run() {
	if len(w.cfg.Users) == 0 {
		// 没有需要设置的用
		return
	}

	friends, err := w.bot.Friends()
	if nil != err {
		logrus.Errorf("failed get friends iter ,err =%s", err.Error())
		return
	}
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
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		if len(result) != 0 {
			// 设置微信步数
			app := zeepLife.NewZeppLife(w.cfg.Users[i].Account, w.cfg.Users[i].Pwd)
			step := r.Intn(w.cfg.Users[i].Max-w.cfg.Users[i].Min+1) + w.cfg.Users[i].Min
			err = app.SetSteps(step)
			msg := fmt.Sprintf("【定时任务】已经成功帮你设置微信步数: %d", step)
			if err != nil {
				logrus.Errorf("failed to set step")
				msg = fmt.Sprintf("【定时任务】设置微信步数失败,原因  【%s】", err.Error())
			}
			w.bot.SendTextToFriend(result[0], msg)
		}
	}
}
