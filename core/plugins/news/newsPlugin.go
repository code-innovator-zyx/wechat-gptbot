package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"wechat-gptbot/core/plugins"
)

const NewsPluginName = "NewsPlugin"

type plugin struct {
	url string
}

func NewPlugin() plugins.PluginSvr {
	return &plugin{"https://i.news.qq.com/gw/event/pc_hot_ranking_list?ids_hash=&offset=0&page_size=50&appver=15.5_qqnews_7.1.60&rank_id=hot"}
}

func (p plugin) Do(i ...interface{}) string {
	newsRes, err := fetchNews(p.url)
	if err != nil {
		log.Fatalf("error fetching news: %v", err)
		return "热点新闻获取失败"
	}
	builder := strings.Builder{}
	builder.WriteString("-----今日热点-----\n")

	for _, news := range newsRes.List[0].NewsList[1:] {
		builder.WriteString(fmt.Sprintf("\n%d ℹ️%s\n⏰ %s\n🔗 %s\n ",
			news.HotEvent.Ranking, news.HotEvent.Title, news.Time, news.Url))
	}
	return builder.String()
}

func (p plugin) Name() string {
	return NewsPluginName
}

func (p plugin) Scenes() string {
	return "每日热点新闻"
}

func (p plugin) IsUseful() bool {
	return true
}

func (p plugin) Args() []interface{} {
	return nil
}

type newsResponse struct {
	List []struct {
		NewsList []struct {
			Url      string `json:"url"` // 原文链接
			Time     string `json:"time"`
			HotEvent struct {
				Ranking int    `json:"ranking"` // 热点序号
				Title   string `json:"title"`   // 标题
			} `json:"hotEvent"` // 热点事件
		} `json:"newslist"`
	} `json:"idlist"`
}

func fetchNews(apiURL string) (*newsResponse, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var news newsResponse
	err = json.Unmarshal(body, &news)
	if err != nil {
		return nil, err
	}

	return &news, nil
}
