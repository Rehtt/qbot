package comment

import (
	"github.com/Rehtt/qbot/cqhttp_bot"
)

func DefaultComment(c *Cmd) map[string]*Comment {
	return map[string]*Comment{
		"help": {
			Name: "help",
			Help: "显示帮助",
			Run: func(thisComment *Comment, paramete string, bot *cqhttp_bot.Bot, messageType cqhttp_bot.EventMessageType, messageId int32, senderQid, groupId int64, message *cqhttp_bot.EventMessage) {
				switch messageType {
				case cqhttp_bot.Private:
					bot.SendPrivateMsg(senderQid, c.help())
				case cqhttp_bot.Group:
					bot.SendMsg(groupId, cqhttp_bot.MessageArray(cqhttp_bot.TextMessage(c.help())), messageType)
				}
			},
		},
	}
}
