package weather

import (
	"testing"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/6/13 10:02
* @Package:
 */

func TestWeatherPlugin_Do(t *testing.T) {
	plugin := NewWeatherPlugin()
	t.Log(plugin.Do("成都"))
}
