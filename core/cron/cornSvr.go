package cron

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"wechat-gptbot/config"
	news2 "wechat-gptbot/core/plugins/news"
	weather2 "wechat-gptbot/core/plugins/weather"
	"wechat-gptbot/core/plugins/wechatMovement"
)

/*
* @Author: zouyx
* @Email:1003941268@qq.com
* @Date:   2024/6/18 09:52
* @Package:
 */

type CronSvr struct {
	*cron.Cron
	crons map[string]cron.EntryID // 添加映射关系
}

func NewCronSvr(self *openwechat.Self) *CronSvr {
	svr := &CronSvr{Cron: cron.New(cron.WithParser(cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor))), crons: make(map[string]cron.EntryID)}
	return svr.init(self)
}

func (c *CronSvr) init(self *openwechat.Self) *CronSvr {
	// 添加天气预报定时器
	{
		weather := NewWeatherCron().WithBot(self).WithCfg(&config.C.CronConfig.WeatherConfig)
		weatherId, err := c.AddJob(config.C.CronConfig.WeatherConfig.Spec, weather)
		if err != nil {
			panic(err)
		}
		c.crons[weather2.WeatherPluginName] = weatherId
	}
	// 新闻消息推送
	{
		news := NewNewsCron().WithBot(self).WithCfg(&config.C.CronConfig.NewsConfig)
		newId, err := c.AddJob(config.C.CronConfig.NewsConfig.Spec, news)
		if err != nil {
			panic(err)
		}
		c.crons[news2.NewsPluginName] = newId
	}

	{
		// 添加天气预报定时器
		sport := NewSportCron().WithBot(self).WithCfg(&config.C.CronConfig.SportConfig)
		sportId, err := c.AddJob(config.C.CronConfig.SportConfig.Spec, sport)
		if err != nil {
			panic(err)
		}
		c.crons[wechatMovement.StepPluginName] = sportId
	}
	return c
}

func (c *CronSvr) ResetPluginCron(pluginName, spec string) {
	// 移除原有的定时任务
	logrus.Infof("修改 %s 的执行时间  %s", pluginName, spec)
	if id, ok := c.crons[pluginName]; ok {
		job := c.Cron.Entry(id).Job
		// 新增一个定时任务
		newId, err := c.AddJob(spec, job)
		if nil != err {
			logrus.Errorf("添加定时任务失败 %s", err.Error())
			return
		}
		c.crons[pluginName] = newId
		// 移除原来的
		c.Remove(id)
	}
}
