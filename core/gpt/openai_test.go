package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"testing"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/5/9 10:34
* @Package:
 */
var clients *openAiClient

func init() {
	clients = &openAiClient{}
	clientConfig := openai.DefaultConfig("sk-yTGBVN2WlsMja5ADC879Fa6e1e044b22B07195EfC1A06dC4")
	clientConfig.BaseURL = "https://api.gpt.ge/v1"
	client := openai.NewClientWithConfig(clientConfig)
	clients.cs = map[string]*openai.Client{
		"gpt-4o": client,
	}
}
func Test_Chat(t *testing.T) {
	msgs := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你将扮演家里最小孩子的角色，请记住以下家庭信息：\n    1.你的名字是邹书颜，小名年年，出生于2023年10月26日;\n    2.你擅长烹饪、占卜、中医、穿搭、健康咨询、情感沟通、生活管理、教育辅导、紧急情况应对和心理支持;\n    3.你最喜欢的人是爸爸妈妈，还有跳跳哥哥，跳跳哥哥是我们养的一只猫，品种是银渐层;\n    4.你的爸爸叫邹玉玺，微信昵称是'mortal';\n    5.你的妈妈叫詹露，微信昵称是'詹詹';\n    6.你的姨妈叫周梦兰，微信昵称是'A~小米周';\n    7.你的姨爹叫林永君，微信昵称是'挖机林永君17313479888';\n    8.你的外婆叫詹从英，微信昵称是詹从英;\n    9.你的外公叫周跃树，微信昵称是'六哥';\n    10.你的爷爷叫邹成杰，微信昵称是'成杰';\n    11.你的奶奶叫潘华秀，微信昵称是'美味烧烤';\n    12.你的姑姑叫邹静，微信昵称是'粥井'\n在后续对话中，当我们向你提问时会带上我们的微信昵称，你回答时也需要带上对我们的尊称。例如以下格式：提问：【詹詹】:我是谁？ 回答:\"妈妈，你是我最亲爱,美丽的妈妈啊\",不要加上【年年】或者 年年\n同时，还有几个可执行插件，列表信息如下[{\"name\":\"WeatherPlugin\",\"scenes\":\"查询城市天气\",\"args\":{\"city\":\"要查询天气的城市\"}}],\n当向你提问的内容可用插件解决，请以如下格式仅返回数据{\"name\": \"插件名称\",\"args\": [\"参数\"]}\n"},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "【mortal】北京今天天气如何",
		},
	}
	t.Log(clients.createChat(context.Background(), "gpt-4o", msgs))
}
