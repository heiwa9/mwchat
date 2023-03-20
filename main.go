package main

import (
	"flag"
	"fmt"
    "strings"

	chatgpt "github.com/DanPlayer/chatgpt-sdk/v1"
    "github.com/skip2/go-qrcode"
	"github.com/eatmoreapple/openwechat"
)

func main() {
	var SecretKey = flag.String("key", "", "Input Your OpenAi SecretKey")
	var proxy = flag.String("proxy", "", "Input Your proxy url")
	flag.Parse()

	var opt chatgpt.ChatGptOption
	if *proxy == "" {
		opt = chatgpt.ChatGptOption{SecretKey: *SecretKey}
	} else {
		opt = chatgpt.ChatGptOption{SecretKey: *SecretKey, HasProxy: true, ProxyUrl: *proxy}
	}

	var ChatGpt = chatgpt.Client(opt)

	var cmc = NewChatMsgCtx()
	cmc.ClearCtxTask()

	// qbot := openqq.NewClient(692455239, "something.core")

	//WeChat
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if !msg.IsText() {
			return
		}
		content, found := strings.CutPrefix(msg.Content, "/chat ")
		if !found {
			return
		}
		if content == "清除会话" {
			cmc.FindCtx(msg.FromUserName).Clear()
			msg.ReplyText("上下文已清空")
			return
		}
		mc := cmc.FindCtx(msg.FromUserName)
		if mc.IsLock() {
			msg.ReplyText("请等待上条内容回复")
		}
		mc.Lock()
		mc.Add(chatgpt.ChatMessage{
			Role:    "user",
			Content: content,
		})

		response, err := ChatGpt.CreateChatCompletion(chatgpt.CreateChatCompletionRequest{
			Model:    chatgpt.GPT3Dot5Turbo0301,
			Messages: mc.HistoryMsgs,
		})
		if err != nil {
			msg.ReplyText(fmt.Sprintf("completions error: %s", err.Error()))
			mc.UnLock()
			return
		}
		var res string
		for _, c := range response.Choices {
			res += c.Message.Content
		}
		msg.ReplyText(res)
		mc.UnLock()
	}
	// 注册登陆二维码回调
	bot.UUIDCallback = func(uuid string) {
        qrurl := "https://login.weixin.qq.com/l/"+uuid
        fmt.Println(qrurl)
        q, _ := qrcode.New(qrurl, qrcode.Low)
        fmt.Println(q.ToString(true))
	}

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()

	// 登陆
	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	//    self, err := bot.GetCurrentUser()
	//    if err != nil {
	//        fmt.Println(err)
	//        return
	//    }

	// 获取所有的好友
	//    friends, err := self.Friends()
	//    fmt.Println(friends, err)

	// 获取所有的群组
	//    groups, err := self.Groups()
	//    fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
