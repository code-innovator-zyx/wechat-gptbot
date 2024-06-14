# 给孩子或宠物创建一个 GPT 机器人

> 项目地址: [https://github.com/code-innovator-zyx/wechat-gptbot](https://github.com/code-innovator-zyx/wechat-gptbot)

最近家里迎来了一个新的生命，为了和她沟通交流，我创建了一个微信机器人账号，方便家人提前与她互动。这个项目不仅可以帮助培养她的微信账号，以后还能直接将微信号过继给她。☺️

## 新功能
 ![天气.png](docs/天气.png)
![新闻.png](docs/新闻.png)
- 时间: 2024年 6月14日
- 功能描述：自定义插件(目前有 天气预报、每日热点新闻)
- 可定制化: 可以自定义插件，在目录`core/plugins`目录下根据模版实现几个接口注册进插件管理器就行了

例如：

```go
// 新增热点新闻插件
package news

import (
	"wechat-gptbot/core/plugins"
)

const NewsPluginName = "NewsPlugin"

type plugin struct {
}

func NewPlugin() plugins.PluginSvr {
	return &plugin{}
}

// 执行插件
func (p plugin) Do(i ...interface{}) string {
	return "插件返回的结果"
}

// 插件名称
func (p plugin) Name() string {
	return NewsPluginName
}

// 插件场景描述
func (p plugin) Scenes() string {
	return "每日热点新闻"
}

// 插件是否可用
func (p plugin) IsUseful() bool {
	return true
}

// 运行插件需要的参数解释 比如我执行天气预报插件的参数  return []interface{}{"要查询天气的城市"}
func (p plugin) Args() []interface{} {
	return nil
}
```
## 项目优势

- **部署简单**：使用 Golang 编译的二进制文件，避免了其他语言依赖众多库的麻烦，直接运行即可。
- **突破微信登录限制**：使用桌面版微信协议，感谢开源项目 [openwechat](https://github.com/eatmoreapple/openwechat)。

## 功能列表

- **文本对话**：接收私聊/群聊消息，使用 OpenAI 的 GPT-4-turbo 生成回复内容，自动回答问题。
- **用户级对话上下文顺序保证**：确保每个用户的对话按提问顺序生成上下文。
- **触发口令设置**：
    1. 私聊时无需额外触发口令，直接对话。
    2. 群聊中需@对方或使用指定口令开头触发对话。
- **连续对话**：支持私聊/群聊开启连续对话功能，默认记忆最近三组对话及最初提示词，保持角色设定。
- **图片生成**：根据描述生成图片，并自动回复在当前私聊/群聊中。
- **称谓识别**：根据提示词识别聊天对象并带上对方称谓回复。
- **会话隔离**：不同用户与机器人对话，系统管理不同的 session。
- **生成图片压缩**：压缩生成的图片以便传输。
- **聊天模型配置化**：可自定义聊天模型。
- **模型代理切换**：支持使用 OpenAI 代理地址。
- **微信朋友圈插件**：修改微信计步器，控制每天微信运动步数（待接入机器人）。
- **支持所有文本对话模型，包括最新的 GPT-4**：目前仅支持对话，暂不支持图片和语音。
- **大模型交互界面**：纯后端实现交互界面，待更新更多功能(选装，如果你有安装python环境，会自动起一个UI界面)。
- **实时天气预报查询**：可以随时查询全国各地的天气预报，只需要直接询问即可

## 待实现功能

- **群聊和私信消息隔离**：不同用户的群聊信息和私信信息上下文隔离。
- **GPT-4 语音对话**：通过微信与 GPT-4 进行语音对话。
- **更多功能待补充**：欢迎大家提供意见。

## (选装)集成 UI 界面展示

需要安装 Python 环境，或使用 Docker 构建，所有环境已打包成基础镜像，见 Dockerfile。

### 优点

交互无需适配手机端，支持公网访问，手机可直接访问和修改配置。

![登录](docs/登录.png)
![UI 界面](ui.png/img.png)
![界面预览](docs/img.png)

## 聊天效果预览

先看使用效果，之后再介绍如何部署及配置。下图展示了**群聊对话**、**私聊对话**和画图的一些例子：

| ![群聊1](docs/群聊1.jpg) | ![群聊](docs/群聊.jpg) |
|----|--------------------|
| ![私聊](docs/私聊.jpg) |                    |
| ![画图](docs/画图.jpg) |                    |


## 开始部署

### 一、环境准备

- 一台服务器或长期开启的 PC
- OpenAI 账号及生成的 `SECRET KEY`
- 一个微信账号

> 注：OpenAI 的域名 `https://api.openai.com` 在国内可能无法访问，需自行解决 API 访问问题。
> 如果遇到困难，可使用国内代理，推荐 [https://api.v3.cm/register?aff=5KuW](https://api.v3.cm/register?aff=5KuW)，价格更便宜。

### 二、配置

1. 克隆代码到本地

你会看到包里有三个文件：

- `/config/config.json.template` : 机器人的基础配置文件，运行前需要修改。
- `/config/prompt.conf.template`: 给 OpenAI 语言模型的提示语。
- `wechat-gptbot` ：自行编译的二进制文件。

配置步骤：

将 `config.json.template` 重命名为 `config.json`，然后用文本编辑器修改此文件：

```json
{
  "gpt": {
    "text_config": {
      "proxy_url": "",
      "auth_token": "openai token",
      "trigger_prefix": "嘎子"
    },
    "image_config": {
      "proxy_url": "",
      "auth_token": "openai token",
      "trigger_prefix": "二狗子"
    }
  },
  "context_status": true,
  "base_model": "gpt-3.5-turbo"
}
```

- `text_config` 和 `image_config` 代表文本对话和图片生成的配置：
    - `proxy_url`：访问 OpenAPI 接口的代理地址，若可直接访问外网无需填写。
    - `auth_token`：在 OpenAI 官网生成的 `SECRET KEY`。
    - `trigger_prefix`：微信对话中触发 AI 回复的前缀。
- `context_status`：开启上下文聊天。

将 `prompt.conf.template` 重命名为 `prompt.conf`，用编辑器修改此文件，确保信息正确：

```txt
你将扮演家里最小孩子的角色，请记住以下信息：
    1.你的名字是***，小名**，出生于2023年10月26日；
    2.你擅长烹饪、占卜、中医、穿搭、健康咨询、情感沟通、生活管理、教育辅导、紧急情况应对和心理支持；
    3.你最喜欢的人是爸爸妈妈，还有**哥哥；
    4.你的爸爸叫**，昵称是'**'；
    5.你的妈妈叫**，昵称是'**'；
    6.你的姨妈叫**，昵称是'**'；
    7.你的姨爹叫**，昵称是'**'；
    8.你的外婆叫**，昵称是'**'；
    9.你的外公叫**，昵称是'**'；
    10.在后续对话中，当我们向你提问时会带上我们的昵称，你回答时也需要带上对我们的称谓。
```

可以根据需要调整文件内容，描述机器人的特点。

### 三、运行

配置完成后，编译并执行二进制文件：

```shell
# 本地运行
./run.sh
```

```shell
# Docker 运行
./dockerRun.sh
```

首次执行时，会出现二维码提示登录微信，用机器人的微信账号扫码登录。

- 登录完成后，会生成一个 `token.json` 文件保存当前微信的登录状态，实现热登录，避免每次运行都需要扫码。

## 完成部署

至此，已完成机器人的部署，快去微信中找好友试试吧！

## 代理设置

如果在服务器上运行且服务器开启了代理，需要设置环境变量：

```bash
export http_proxy=http://127.0.0.1:xxxx
export https_proxy=https://127.0.0.1:xxxx
```

## 其他

微信登录时常在夜晚掉线，内置了保活功能，~~需要关注公众号
“跳跳是只cat”，有需要的可以自行关注，项目会对公众号进行请求心跳保活。~~

## 联系作者

- 项目地址：[https://github.com/code-innovator-zyx/wechat-gptbot](https://github.com/code-innovator-zyx/wechat-gptbot)
  ，欢迎 Star，提交 PR。
- 有问题可在项目下提 `Issues`。