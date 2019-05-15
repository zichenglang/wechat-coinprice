package autoreply

import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
)
// register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(autoReply), "text-replier")
	if err := session.HandlerRegister.Add(wxweb.MSG_IMG, wxweb.Handler(autoReply), "img-replier"); err != nil {
		logs.Error(err)
	}

	if err := session.HandlerRegister.EnableByName("text-replier"); err != nil {
		logs.Error(err)
	}

	if err := session.HandlerRegister.EnableByName("img-replier"); err != nil {
		logs.Error(err)
	}

}

func msgFilter(in string) string {
	return ""
}
func autoReply(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	if !msg.IsGroup {
		session.SendText("暂时不在，稍后回复", session.Bot.UserName, msg.FromUserName)
	}
}
