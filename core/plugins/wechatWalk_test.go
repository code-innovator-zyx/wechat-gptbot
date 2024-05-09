package plugins

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/5/9 14:31
* @Package:
 */

func TestWechatWalk(t *testing.T) {
	phone := "786618102@qq.com"
	password := "4f4ezha!"
	// 测试示例
	login(phone, password)
}

func login(user, pwd string) {
	//url1 := "https://api-user.huami.com/registrations/" + user + "/tokens"
	login_headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
		"User-Agent":   "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2",
	}
	//data1 := url.Values{}
	//data1.Set("client_id", "HuaMi")
	//data1.Set("password", pwd)
	//data1.Set("redirect_uri", "https://s3-us-west-2.amazonaws.com/hm-registration/successsignin.html")
	//data1.Set("token", "access")
	//
	//fmt.Println(url1)
	//req, _ := http.NewRequest(http.MethodPost, url1, strings.NewReader(data1.Encode()))
	//for key, value := range login_headers {
	//	req.Header.Set(key, value)
	//}
	//fmt.Println(req.Header)
	//
	//client := &http.Client{
	//	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	//		return http.ErrUseLastResponse
	//	},
	//}
	//
	//resp, _ := client.Do(req)
	//if resp.StatusCode != http.StatusSeeOther {
	//	fmt.Printf("登录异常1，status: %d\n", resp.StatusCode)
	//	return
	//}
	//location := resp.Header.Get("Location")
	//location := "https://s3-us-west-2.amazonaws.com/hm-registration/successsignin.html?region=us-west-2&access=ZQVBQDZOQmJaR0YyajYmWnJoBAgAAAAAAYT1aTzdzaW5HUk1Kbm1MaHdDUmxEcWZWd0FBQVk5Y2JDcHcmcj0xMiZ0PWh1YW1pJnRpPTc4NjYxODEwMkBxcS5jb20maD0xNzE1MjQyOTE3NTA4Jmk9ODY0MDAwJnVzZXJuYW1lPTc4NjYxODEwMl7QyjErjcPCXgCFRFk9pZk&country_code=CN&expiration=1716106917"
	//code := getAccessToken(location)
	code := "ZQVBQDZOQmJaR0YyajYmWnJoBAgAAAAAAYT1aTzdzaW5HUk1Kbm1MaHdDUmxEcWZWd0FBQVk5Y2JDcHcmcj0xMiZ0PWh1YW1pJnRpPTc4NjYxODEwMkBxcS5jb20maD0xNzE1MjQyOTE3NTA4Jmk9ODY0MDAwJnVzZXJuYW1lPTc4NjYxODEwMl7QyjErjcPCXgCFRFk9pZk"

	url2 := "https://account.huami.com/v2/client/login"

	data2 := url.Values{}
	data2.Set("allow_registration=", "false")
	data2.Set("app_name", "com.xiaomi.hm.health")
	data2.Set("app_version", "https://s3-us-west-2.amazonaws.com/hm-registration/successsignin.html")
	data2.Set("code", code)
	data2.Set("country_code", "CN")
	data2.Set("device_id", "2C8B4939-0CCD-4E94-8CBA-CB8EA6E613A1")
	data2.Set("device_model", "phone")
	data2.Set("dn", "api-user.huami.com%2Capi-mifit.huami.com%2Capp-analytics.huami.com")
	data2.Set("grant_type", "access_token")
	data2.Set("lang", "zh_CN")
	data2.Set("os_version", "1.5.0")
	data2.Set("source", "com.xiaomi.hm.health")
	data2.Set("third_name", "email")

	req2, _ := http.NewRequest(http.MethodPost, url2, strings.NewReader(data2.Encode()))
	for key, value := range login_headers {
		req2.Header.Set(key, value)
	}
	resp2, _ := http.DefaultClient.Do(req2)
	if resp2.StatusCode != 200 {
		fmt.Printf("登录异常，status: %d\n", resp2.StatusCode)
		return
	}
	respBody, _ := io.ReadAll(resp2.Body)
	fmt.Println("========")
	fmt.Println(string(respBody))
}

func getAccessToken(location string) string {
	// 定义正则表达式模式，匹配整个access=...部分，并捕获其中的内容
	codePattern := regexp.MustCompile(`access=([^&$]*)`)

	// 在给定的location中查找匹配的部分
	match := codePattern.FindStringSubmatch(location)

	// 如果未找到匹配项，则返回空字符串
	if len(match) < 2 {
		return ""
	}
	// 返回捕获的内容
	return match[1]
}
