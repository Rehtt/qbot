/**
 * Copyright (c) 2023 Rehtt
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
			_, _ = ctx.QuickReplyText(bot, help, true)
			//switch ctx.MessageType {
			//case cqhttp_bot.Private:
			//	bot.SendPrivateMsg(ctx.SenderQid, help)
			//case cqhttp_bot.Group:
			//	bot.SendMsg(groupId, cqhttp_bot.MessageArray(cqhttp_bot.TextMessage(help)), messageType)
			//}
		},
	}
}
