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
	"time"

	jsoniter "github.com/json-iterator/go"
)

func (b *Event) event(postType string, data jsoniter.Any) {
	switch postType {
	case "message", "message_sent":
		if len(b.onMessage) == 0 && len(b.onPrivateMessages) == 0 && len(b.onGroupMessages) == 0 {
			break
		}
		b.MessageEvent.eventMessage(data)
	case "meta_event":
		// 暂时忽略
	case "request":
		b.RequestEvent.eventRequest(data)
	case "notice":
		b.NoticeEvent.eventNotice(data)
	}

	if b.callback != nil {
		for _, event := range []string{postType, ""} {
			for _, f := range b.callback[event] {
				f(postType, data)
			}
		}
	}
}

// 所有事件回调
func (e *Event) AllEventCallback(f EventCallback) {
	e.eventCallback("", f)
}

// 指定事件回调
func (e *Event) SpecifyEventCallback(event string, f EventCallback) {
	e.eventCallback(event, f)
}

func (e *Event) eventCallback(event string, f EventCallback) {
	if e.callback == nil {
		e.callback = make(map[string][]EventCallback)
	}
	e.callback[event] = append(e.callback[event], f)
}

func (b *NoticeEvent) eventNotice(data jsoniter.Any) {
	// todo 事件
	// todo 消息撤回
	// {"post_type":"notice","notice_type":"group_recall","time":1699759315,"self_id":1033853263,"group_id":852122585,"user_id":1033853263,"operator_id":1033853263,"message_id":1258115355}}
}

func (b *MessageEvent) ParseEventMessage(data jsoniter.Any) *EventMessageContext {
	var (
		senderQid = data.Get("user_id").ToInt64()
		messageId = data.Get("message_id").ToInt32()
		m         = NewEventMessage()
		ctx       = NewEventMessageContext()
	)
	defer m.Close()
	defer ctx.Close()
	m.RawMessage = data.Get("raw_message").ToString()
	m.Messages = ParseMessage(m.RawMessage)

	ctx.MessageId = messageId
	ctx.SenderId = senderQid
	ctx.Sender = Sender{
		Nickname: data.Get("sender", "nickname").ToString(),
		Card:     data.Get("sender", "card").ToString(),
	}
	ctx.Time = time.Unix(data.Get("time").ToInt64(), 0)
	ctx.Message = m
	ctx.GroupId = 0

	switch data.Get("message_type").ToString() {
	case "group":
		groupId := data.Get("group_id").ToInt64()
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
	return ctx
}

func (b *MessageEvent) eventMessage(data jsoniter.Any) {
	ctx := b.ParseEventMessage(data)
	for i := range b.onMessage {
		b.onMessage[i](ctx)
	}
}

// OnMessage 接收所有消息
func (b *MessageEvent) OnMessage(f OnMessageFunc) {
	b.onMessage = append(b.onMessage, f)
}

// OnGroupMessage 接收群消息
func (b *MessageEvent) OnGroupMessage(f OnGroupMessageFunc) {
	b.onGroupMessages = append(b.onGroupMessages, f)
}

// OnPrivateMessage 接收私人消息
func (b *MessageEvent) OnPrivateMessage(f OnPrivateMessageFunc) {
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

func (r *RequestEvent) OnFriendRequest(f OnFriendRequestFunc) {
	r.onFriendRequests = append(r.onFriendRequests, f)
}

func (r *RequestEvent) OnGroupRequest(f OnGroupRequestFunc) {
	r.onGroupRequests = append(r.onGroupRequests, f)
}
