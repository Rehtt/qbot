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
	jsoniter "github.com/json-iterator/go"
)

type Event struct {
	onGroupMessages   []onGroupMessage
	onPrivateMessages []onPrivateMessage
}
type onGroupMessage func(senderQid, groupId int64, message *EventMessage)
type onPrivateMessage func(userId int64, message *EventMessage)

func (b *Event) event(postType string, data jsoniter.Any) {
	switch postType {
	case "message":
		b.eventMessage(data)
	case "meta_event":
		// 暂时忽略
		return
	}
}

func (b *Event) eventMessage(data jsoniter.Any) {
	var (
		senderQid  = data.Get("user_id").ToInt64()
		rawMessage = data.Get("raw_message").ToString()
		m          = &EventMessage{
			RawMessage: rawMessage,
			Messages:   ParseMessage(rawMessage),
		}
	)
	switch data.Get("message_type").ToString() {
	case "group":
		var (
			groupId = data.Get("group_id").ToInt64()
		)
		for i := range b.onGroupMessages {
			b.onGroupMessages[i](senderQid, groupId, m)
		}
	case "private":
		for i := range b.onPrivateMessages {
			b.onPrivateMessages[i](senderQid, m)
		}
	}
}

func (b *Event) OnGroupMessage(f onGroupMessage) {
	b.onGroupMessages = append(b.onGroupMessages, f)
}
func (b *Event) OnPrivateMessage(f onPrivateMessage) {
	b.onPrivateMessages = append(b.onPrivateMessages, f)
}
