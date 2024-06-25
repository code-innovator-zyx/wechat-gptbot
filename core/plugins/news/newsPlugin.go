package news

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"wechat-gptbot/core/plugins"
)

const (
	NewsPluginName = "NewsPlugin"
	TencentSource  = "https://i.news.qq.com/gw/event/pc_hot_ranking_list?ids_hash=&offset=0&page_size=50&appver=15.5_qqnews_7.1.60&rank_id=hot"
)

type plugin struct {
	url     string
	topN    int
	fromRss bool
}

// Option 熔断器配置
type Option func(*plugin)

// SetTopN 设置TopN
func SetTopN(topN int) Option {
	return func(p *plugin) {
		p.topN = topN
	}
}

// SetRssSource 设置rss 源
func SetRssSource(source string) Option {
	return func(p *plugin) {
		p.url = strings.TrimSpace(source)
		p.fromRss = true
	}
}

func NewPlugin(opts ...Option) plugins.PluginSvr {
	p := plugin{}
	for _, o := range opts {
		o(&p)
	}
	if p.topN == 0 {
		p.topN = 10
	}
	if p.url == "" {
		p.url = TencentSource
		p.fromRss = false
	}
	return p
}

func (p plugin) Do(...interface{}) []string {
	fun := tencent
	if p.fromRss {
		fun = rss
	}
	news, err := fun(p.url, p.topN)
	if err != nil {
		logrus.Errorf("error fetching news: %v", err)
		return []string{"暂时无法帮你查看，请检查配置"}
	}
	return news
}

func (p plugin) Name() string {
	return NewsPluginName
}

func (p plugin) Scenes() string {
	return "查看实时热点或者订阅消息"
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

// 腾讯新闻
func tencent(apiURL string, topN int) ([]string, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var newsResp newsResponse
	err = json.Unmarshal(body, &newsResp)
	if err != nil {
		return nil, err
	}
	builder := strings.Builder{}
	// 为了减少触发腾讯的内容风控(组合关键词)，这里进行批处理返回
	batchSize := 5
	reply := make([]string, 0, (topN/batchSize)+1)
	for index, news := range newsResp.List[0].NewsList[1:] {
		if index%batchSize == 0 {
			builder.WriteString(fmt.Sprintf("🔥腾讯实时热点 【%d】\n\n", len(reply)+1))
		}
		builder.WriteString(fmt.Sprintf("\n%d ℹ️%s\n⏰ %s\n🔗 %s\n ",
			news.HotEvent.Ranking, news.HotEvent.Title, news.Time, news.Url))
		if index%batchSize == batchSize-1 {
			reply = append(reply, builder.String())
			builder.Reset()
		}
		if index+1 >= topN {
			break
		}
	}
	if builder.Len() != 0 {
		reply = append(reply, builder.String())
	}
	return reply, nil
}

// 自定义rss资源
func rss(source string, topN int) ([]string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(source)
	if err != nil {
		return nil, err
	}
	batchSize := 5
	reply := make([]string, 0, (topN/batchSize)+1)
	builder := strings.Builder{}
	for index, item := range feed.Items {
		if index%batchSize == 0 {
			builder.WriteString(fmt.Sprintf("🆙 %s 【%d】\n\n", feed.Title, len(reply)+1))
		}
		builder.WriteString(fmt.Sprintf("\n❤️%s\n🔗 %s\n",
			item.Title, item.Link))
		if index%batchSize == batchSize-1 {
			reply = append(reply, builder.String())
			builder.Reset()
		}
		if index+1 >= topN {
			break
		}
	}
	if builder.Len() != 0 {
		reply = append(reply, builder.String())
	}
	return reply, nil
}
