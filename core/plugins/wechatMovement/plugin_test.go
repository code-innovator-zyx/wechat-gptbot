package wechatMovement

import (
	"fmt"
	"testing"
	"wechat-gptbot/core/plugins"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/6/3 16:01
* @Package:
 */
var Registry *plugins.PluginManger

func init() {
	fmt.Println("init Registry")
	Registry = plugins.NewPluginRegistry()
}
func TestPluginManger(t *testing.T) {

	t.Run("AddPlugin", func(t *testing.T) {

	})
	t.Run("ResetPlugin", func(t *testing.T) {

	})
	t.Run("DoPlugin", func(t *testing.T) {

	})

	t.Run("PluginPrompt", func(t *testing.T) {
		Registry.Register(NewStepPlugin("", "", 0, 0))
		fmt.Println(Registry.PluginPrompt())
	})
}
