package plugins

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"sync"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2024/5/31 11:09
* @Package:
 */

// 插件

const (
	PluginTypeUnknown = iota
	PluginTypeGpt     // gpt 插件 通过gpt的回答进行插件选择并执行插件
	PluginTypeSystem  // 系统插件 定期执行的插件
)

type PluginSvr interface {
	// Do 执行插件
	Do(...interface{}) string
	// Name 获取插件名称
	Name() string
	// Scenes 使用场景
	Scenes() string
	// IsUseful 插件是否可用
	IsUseful() bool
	Args() []interface{}
}

// PluginManger 插件管理器
type PluginManger struct {
	*sync.Map // 存放插件
}

func NewPluginRegistry() *PluginManger {
	return &PluginManger{new(sync.Map)}
}

var Manger *PluginManger

func init() {
	Manger = NewPluginRegistry()
}

func (m *PluginManger) DoPlugin(msg string) (resetMsg string, ok bool) {
	var pp pluginPrompt
	err := jsoniter.UnmarshalFromString(msg, &pp)
	if err != nil {
		return msg, false
	}
	m.Range(func(key, value any) bool {
		if key.(string) == pp.Name {
			plugin := value.(PluginSvr)
			if !plugin.IsUseful() {
				resetMsg = fmt.Sprintf("%s 不可用状态", plugin.Name())
				return false
			}
			// 执行插件
			resetMsg = value.(PluginSvr).Do(pp.Args...)
			ok = true
			return false
		}
		return true
	})
	return
}

// Register 注册插件
func (m *PluginManger) Register(svr ...PluginSvr) {
	for i := range svr {
		m.Store(svr[i].Name(), svr[i])
	}
}

// ResetPlugin 修改插件
func (m *PluginManger) ResetPlugin(name string, fun func(svr PluginSvr) PluginSvr) error {
	value, exist := m.Load(name)
	if !exist {
		return errors.New("插件不存在")
	}
	resetSvr := fun(value.(PluginSvr))
	m.Store(resetSvr.Name(), resetSvr)
	return nil
}

// RemovePlugin 移除插件
func (m *PluginManger) RemovePlugin(name string) {
	m.Delete(name)
}

func (m *PluginManger) GetPlugin(name string) (PluginSvr, error) {
	value, exist := m.Load(name)
	if !exist {
		return nil, errors.New("插件不存在")
	}
	return value.(PluginSvr), nil
}

type pluginPrompt struct {
	Name   string        `json:"name"`
	Scenes string        `json:"scenes"`
	Args   []interface{} `json:"args"`
}

// PluginPrompt 构建插件提示词
func (m *PluginManger) PluginPrompt() string {
	var prompts []pluginPrompt
	m.Range(func(_, value any) bool {
		if svr, ok := value.(PluginSvr); ok {
			prompts = append(prompts, pluginPrompt{
				svr.Name(),
				svr.Scenes(),
				svr.Args(),
			})
		}
		return true
	})
	toString, _ := jsoniter.MarshalToString(prompts)

	return toString
}
