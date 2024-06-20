package svc

import (
	"wechat-gptbot/core/cron"
	"wechat-gptbot/core/gpt"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/4/12 11:09
* @Package:
 */

type ServiceContext struct {
	Session    gpt.Session
	CronServer *cron.CronSvr
}

func NewServiceContext() *ServiceContext {
	return &ServiceContext{
		Session: gpt.NewSession(),
	}
}
