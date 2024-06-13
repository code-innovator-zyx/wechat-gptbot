package weather

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"wechat-gptbot/core/plugins"
)

const WeatherPluginName = "WeatherPlugin"

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/6/12 18:08
* @Package:
 */

type WeatherPlugin struct {
	url string
}

func NewWeatherPlugin() plugins.PluginSvr {
	return &WeatherPlugin{"http://139.9.115.47:80/wechat-helper/weather"}
}
func (s WeatherPlugin) Do(args ...interface{}) string {
	fmt.Printf("查询 %s 天气 \n", args[0])
	if len(args) <= 0 {
		return "请输入查询的地址"
	}
	uri := fmt.Sprintf("%s?city=%s", s.url, args[0])
	res, err := http.Get(uri)
	if err != nil {
		logrus.Errorf("failed call weatherSvr %s", err.Error())
		return ""
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusTooManyRequests {
		return "当天天气查询调用次数过多，每天最多查询10次"
	}
	b, _ := io.ReadAll(res.Body)
	data := make(map[string]string, 1)
	json.Unmarshal(b, &data)
	fmt.Println("data", data["msg"])
	return data["msg"]
}
func (s WeatherPlugin) IsUseful() bool {
	return true
}

func (s WeatherPlugin) Name() string {
	return WeatherPluginName
}

func (s WeatherPlugin) Scenes() string {
	return "查询城市天气"
}

func (s WeatherPlugin) Args() []interface{} {
	return []interface{}{"要查询天气的城市"}
}
