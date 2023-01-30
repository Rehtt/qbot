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

package main

import (
	"fmt"
	"github.com/Rehtt/qbot/cqhttp_bot"
	"github.com/Rehtt/qbot/cqhttp_bot/comment"
	"time"
)

func main() {
	bot := cqhttp_bot.New("ws://127.0.0.1:8060")
	bot.Start()
	fmt.Println(bot.GetFriendsList())
	bot.Event.OnPrivateMessage(func(messageId int32, userId int64, message *cqhttp_bot.EventMessage) {
		fmt.Println(userId, message.Messages, message.RawMessage)
	})
	bot.Event.OnGroupMessage(func(messageId int32, senderQid, groupId int64, message *cqhttp_bot.EventMessage) {
		fmt.Println(senderQid, groupId, message.Messages, message.RawMessage)
	})
	bot.SendMsg(852122585, cqhttp_bot.MessageArray(cqhttp_bot.TextMessage("test")), cqhttp_bot.Group)

	// 开启命令模式
	c := comment.New(bot)
	test := &comment.Comment{
		Name:  "test",
		Usage: "test1",
		Run: func(paramete string, flag comment.Flag, bot *cqhttp_bot.Bot, messageType cqhttp_bot.EventMessageType, messageId int32, senderQid, groupId int64, message *cqhttp_bot.EventMessage) {
			fmt.Println("test1", flag, paramete)
		},
	}
	test.Flag().Var("f", "123", "可选参数")
	c.AddComment(test)
	time.Sleep(10 * time.Minute)
}
