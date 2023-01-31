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

package cqhttp_bot

import (
	"fmt"
	"strings"
)

func ParseMessage(raw string) (m Messages) {
	var index int
	for {
		index = strings.Index(raw, "[CQ:")
		if index > 0 {
			m = append(m, TextMessage(raw[:index]))
			raw = raw[index:]
		} else if index == 0 {
			feet := strings.Index(raw, "]")
			if feet == -1 {
				m = append(m, TextMessage(raw))
				break
			}
			data := parse(raw[:feet])
			cq := strings.Split(raw[:feet], ",")
			switch cq[0] {
			case "[CQ:image": // 图片
				m = append(m, ImageMessage(data["file"], data["url"], data["type"] == "flash"))
			case "[CQ:face": // 表情
			case "[CQ:reply": // 回复
				m = append(m, ReplyMessage(data["id"]))
			case "[CQ:at": // @
				m = append(m, AtMessage(data["qq"]))
			}
			raw = raw[feet+1:]
		} else if index == -1 {
			if raw != "" {
				m = append(m, TextMessage(raw))
			}
			break
		}
	}
	return
}
func MessageArray(msg ...Message) Messages {
	return msg
}
func TextMessage(text string) Message {
	return Message{
		Type: TEXT,
		Text: text,
	}
}

func AtMessage(qid any, name ...string) Message {
	m := Message{
		Type: At,
		At: &messageAt{
			Qid: qid,
		},
	}
	if len(name) != 0 && name[0] != "" {
		m.At.Name = name[0]
	}
	return m
}

func ImageMessage(file, url string, flash bool) Message {
	return Message{
		Type: IMAGE,
		Image: &messageImage{
			Url:   url,
			File:  file,
			Flash: flash,
		},
	}
}
func ReplyMessage(messageId string) Message {
	return Message{
		Type: Reply,
		Text: messageId,
	}
}

func (m *Messages) Add(msg Message) Messages {
	*m = append(*m, msg)
	return *m
}
func (m *Messages) TextMessage(text string) Messages {
	return m.Add(TextMessage(text))
}
func (m *Messages) AtMessage(qid any, name ...string) Messages {
	return m.Add(AtMessage(qid, name...))
}

func (m *Messages) ImageMessage(file, url string, flash bool) Messages {
	return m.Add(ImageMessage(file, url, flash))
}
func (m *Messages) ReplyMessage(messageId string) Messages {
	return m.Add(ReplyMessage(messageId))
}
func (m *Messages) RawMessage() string {
	var out strings.Builder
	for _, msg := range *m {
		switch msg.Type {
		case Reply:
			out.WriteString(fmt.Sprintf("[CQ:reply,id=%s]", msg.Text))
		case IMAGE:
			if msg.Image == nil {
				continue
			}
			var imageType = "show"
			if msg.Image.Flash {
				imageType = "flash"
			}
			var file = msg.Image.File
			if file == "" && msg.Image.Url != "" {
				file = msg.Image.Url
			}
			out.WriteString(fmt.Sprintf("[CQ:image,file=%s,type=%s]", file, imageType))
		case TEXT:
			out.WriteString(msg.Text)
		case At:
			if msg.At == nil {
				continue
			}
			out.WriteString(fmt.Sprintf("[CQ:at,qq=%v", msg.At.Qid))
			if msg.At.Name != "" {
				out.WriteString(fmt.Sprintf(",name=%s", msg.At.Name))
			}
			out.WriteString("]")
		}
	}
	return out.String()
}

// QuickReplyText 对消息快速回复
// @param at @发送人，仅当消息类型为Group时有效
func (e *EventMessageContext) QuickReplyText(bot *Bot, msg string, at ...bool) (int32, error) {
	return e.QuickReply(bot, MessageArray(TextMessage(msg)), at...)
}

// QuickReply 对消息快速回复
// @param at @发送人，仅当消息类型为Group时有效
func (e *EventMessageContext) QuickReply(bot *Bot, msg Messages, at ...bool) (int32, error) {
	qid := e.SenderId
	if e.MessageType == Group {
		qid = e.GroupId
		if len(at) != 0 && at[0] {
			msg.TextMessage("\n")
			msg.AtMessage(e.SenderId)
		}
	}
	return bot.SendMsg(qid, msg, e.MessageType)
}
