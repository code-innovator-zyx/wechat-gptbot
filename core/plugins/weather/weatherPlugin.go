package weather

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
	"time"
	"wechat-gptbot/core/plugins"
)

const WeatherPluginName = "WeatherPlugin"

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/6/12 18:08
* @Package:
 */

type Plugin struct {
	url string
	*sync.Map
}

// 同一区域的天气情况缓存起来，防止重复查询
type cache struct {
	city           string // 城市
	data           string // 数据
	cacheTimeStamp int64  // 缓存时间
}

func NewPlugin() plugins.PluginSvr {
	return &Plugin{"http://139.9.115.47:80/wechat-helper/weather", new(sync.Map)}
}

// isExpired 检查缓存数据是否过期
func isExpired(timestamp int64) bool {
	// 缓存有效期为1小时
	const cacheDuration = time.Hour
	// 当前时间
	now := time.Now().Unix()
	// 比较时间戳，判断是否超过缓存有效期
	return now-timestamp > int64(cacheDuration.Seconds())
}

func (s *Plugin) loadData(city string) (cache, bool) {
	var data cache
	if v, ok := s.Load(city); ok {
		data = v.(cache)
		// 检查数据是否过期
		if !isExpired(data.cacheTimeStamp) {
			return data, true
		}
	}
	return data, false
}
func (s *Plugin) addData(city, data string) {
	s.Store(city, cache{
		city:           city,
		data:           data,
		cacheTimeStamp: time.Now().Unix(),
	})
}
func (s *Plugin) Do(args ...interface{}) []string {
	if len(args) <= 0 {
		return []string{"请输入查询的地址"}
	}
	city, ok := args[0].(string)
	if !ok {
		city = "成都"
	}
	// 获取缓存
	if data, ok := s.loadData(city); ok {
		return []string{data.data}
	}
	uri := fmt.Sprintf("%s?city=%s", s.url, city)
	res, err := http.Get(uri)
	if err != nil {
		logrus.Errorf("failed call weatherSvr %s", err.Error())
		return []string{err.Error()}
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusTooManyRequests {
		return []string{"当天天气查询调用次数过多，每天最多查询10次"}
	}
	b, _ := io.ReadAll(res.Body)
	data := make(map[string]string, 1)
	json.Unmarshal(b, &data)
	// 添加缓存
	s.addData(city, data["msg"])
	return []string{data["msg"]}
}
func (s *Plugin) IsUseful() bool {
	return true
}

func (s *Plugin) Name() string {
	return WeatherPluginName
}

func (s *Plugin) Scenes() string {
	return "查询城市天气"
}

func (s *Plugin) Args() []interface{} {
	return []interface{}{"要查询天气的城市"}
}
