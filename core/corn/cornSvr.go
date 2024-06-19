package corn

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron"
	"wechat-gptbot/config"
)

/*
* @Author: zouyx
* @Email:1003941268@qq.com
* @Date:   2024/6/18 09:52
* @Package:
 */

type CronSvr struct {
	*cron.Cron
}

func NewCronSvr(self *openwechat.Self) *CronSvr {
	return &CronSvr{Cron: newCron(self)}
}

func newCron(self *openwechat.Self) *cron.Cron {
	// 初始化定时器
	cr := cron.New()

	// 添加天气预报定时器
	weather := NewWeatherCron().WithBot(self).WithCfg(&config.C.CronConfig.WeatherConfig)
	if err := cr.AddJob(config.C.CronConfig.WeatherConfig.Spec, weather); err != nil {
		panic(err)
	}
	news := NewNewsCron().WithBot(self).WithCfg(&config.C.CronConfig.NewsConfig)
	if err := cr.AddJob(config.C.CronConfig.NewsConfig.Spec, news); err != nil {
		panic(err)
	}
	return cr
}
