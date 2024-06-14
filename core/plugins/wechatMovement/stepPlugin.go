package wechatMovement

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
	"wechat-gptbot/core/plugins"
	"wechat-gptbot/core/plugins/wechatMovement/zeepLife"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2024/5/31 15:29
* @Package: 微信步数插件
 */

type StepPlugin struct {
	*rand.Rand
	account  string
	password string
	max      int
	min      int
}

const StepPluginName = "StepPlugin"

func NewStepPlugin(account, pwd string, min, max int) plugins.PluginSvr {
	source := rand.NewSource(time.Now().UnixNano())
	return &StepPlugin{account: account, password: pwd, min: min, max: max, Rand: rand.New(source)}
}
func (s StepPlugin) Do(args ...interface{}) string {
	app := zeepLife.NewZeppLife(s.account, s.password)
	step := s.Rand.Intn(s.max-s.min+1) + s.min
	err := app.SetSteps(step)
	if err != nil {
		logrus.Errorf("failed to set step")
	}
	return fmt.Sprintf("已经成功帮你设置微信步数: %d", step)
}
func (s StepPlugin) IsUseful() bool {
	return s.account != "" && s.password != "" && s.min != 0 && s.max != 0
}

func (s StepPlugin) Name() string {
	return StepPluginName
}

func (s StepPlugin) Scenes() string {
	return "设置修改微信运动步数"
}

func (s StepPlugin) Args() []interface{} {
	return nil
}
