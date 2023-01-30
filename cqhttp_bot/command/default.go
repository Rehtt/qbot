package command

import (
	"fmt"
	"github.com/Rehtt/qbot/cqhttp_bot"
)

func DefaultCommand(c *Cmd) map[string]*Command {
	return map[string]*Command{
		"help": {
			Name:  "help",
			Usage: "显示帮助",
			Run: func(paramete string, flag Flag, bot *cqhttp_bot.Bot, messageType cqhttp_bot.EventMessageType, messageId int32, senderQid, groupId int64, message *cqhttp_bot.EventMessage) {
				var help string
				if paramete != "" {
					com, _, _ := c.Commands.Parse(paramete)
					if com == nil {
						help = fmt.Sprintf("找不到相关命令 `%s`", paramete)
					} else {
						help = fmt.Sprintf("命令：\n%s\n说明：\n%s\n", com.Usage, paramete)
						help += com.Help()
					}

				} else {
					help = c.Help()
				}
				switch messageType {
				case cqhttp_bot.Private:
					bot.SendPrivateMsg(senderQid, help)
				case cqhttp_bot.Group:
					bot.SendMsg(groupId, cqhttp_bot.MessageArray(cqhttp_bot.TextMessage(help)), messageType)
				}
			},
		},
	}
}
