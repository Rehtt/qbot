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

func (b *Event) event(postType string, data jsoniter.Any) {
	switch postType {
	case "message":
		b.MessageEvent.eventMessage(data)
	case "meta_event":
		// 暂时忽略
	case "request":
		b.RequestEvent.eventRequest(data)
	}
}

func (b *MessageEvent) eventMessage(data jsoniter.Any) {
	var (
		senderQid  = data.Get("user_id").ToInt64()
		rawMessage = data.Get("raw_message").ToString()
		messageId  = data.Get("message_id").ToInt32()
		m          = EventMessage{
			RawMessage: rawMessage,
			Messages:   ParseMessage(rawMessage),
		}
	)
	ctx := EventMessageContext{
		MessageId: messageId,
		SenderId:  senderQid,
		Message:   m,
	}
	switch data.Get("message_type").ToString() {
	case "group":
		var (
			groupId = data.Get("group_id").ToInt64()
		)
		for i := range b.onGroupMessages {
			b.onGroupMessages[i](messageId, senderQid, groupId, m)
		}
		ctx.MessageType = Group
		ctx.GroupId = groupId
	case "private":
		for i := range b.onPrivateMessages {
			b.onPrivateMessages[i](messageId, senderQid, m)
		}
		ctx.MessageType = Private
	}
	for i := range b.onMessage {
		b.onMessage[i](ctx)
	}
}

// OnMessage 接收所有消息
func (b *MessageEvent) OnMessage(f onMessage) {
	b.onMessage = append(b.onMessage, f)
}

// OnGroupMessage 接收群消息
func (b *MessageEvent) OnGroupMessage(f onGroupMessage) {
	b.onGroupMessages = append(b.onGroupMessages, f)
}

// OnPrivateMessage 接收私人消息
func (b *MessageEvent) OnPrivateMessage(f onPrivateMessage) {
	b.onPrivateMessages = append(b.onPrivateMessages, f)
}

func (r *RequestEvent) eventRequest(data jsoniter.Any) {
	var (
		userId  = data.Get("user_id").ToInt64()
		comment = data.Get("comment").ToString()
		flag    = data.Get("flag").ToString()
	)

	switch data.Get("request_type").ToString() {
	case "friend":
		for i := range r.onFriendRequests {
			r.onFriendRequests[i](userId, comment, flag)
		}
	case "group":
		groupId := data.Get("group_id").ToInt64()
		subType := data.Get("sub_type").ToString()
		for i := range r.onGroupRequests {
			r.onGroupRequests[i](userId, groupId, GroupRequestEventSubType(subType), comment, flag)
		}
	}
}
func (r *RequestEvent) OnFriendRequest(f onFriendRequest) {
	r.onFriendRequests = append(r.onFriendRequests, f)
}
func (r *RequestEvent) OnGroupRequest(f onGroupRequest) {
	r.onGroupRequests = append(r.onGroupRequests, f)
}
