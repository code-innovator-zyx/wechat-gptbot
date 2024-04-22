package handler

import "github.com/eatmoreapple/openwechat"

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/4/19 10:41
* @Package:
 */

func doAction(msg *openwechat.Message) {
	if msg.IsTickledMe() {
		msg.ReplyText("别拍啦，小屁股都开花啦")
		return
	}
	if msg.IsJoinGroup() {
		msg.ReplyText("欢迎欢迎，热烈欢迎")
		return
	}
}
