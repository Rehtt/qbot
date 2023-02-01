package command

import (
	"fmt"
	"github.com/Rehtt/qbot/cqhttp_bot"
)

func DefaultCommandHelp(c *Cmd) *Command {
	return &Command{
		Name:  "help",
		Usage: "显示帮助",
		Run: func(paramete string, flag Flag, bot *cqhttp_bot.Bot, ctx *cqhttp_bot.EventMessageContext) {
			var help string
			if paramete != "" {
				com, _, _ := c.Commands.Parse(paramete)
				if com == nil {
					help = fmt.Sprintf("找不到相关命令 `%s`", paramete)
				} else {
					help = fmt.Sprintf("命令：\n%s\n说明：\n%s\n", paramete, com.Usage)
					help += com.Help()
				}

			} else {
				help = c.Help()
			}
			ctx.QuickReplyText(bot, help, true)
			//switch ctx.MessageType {
			//case cqhttp_bot.Private:
			//	bot.SendPrivateMsg(ctx.SenderQid, help)
			//case cqhttp_bot.Group:
			//	bot.SendMsg(groupId, cqhttp_bot.MessageArray(cqhttp_bot.TextMessage(help)), messageType)
			//}
		},
	}
}
