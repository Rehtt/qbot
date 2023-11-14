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

/**
 * @Author: Rehtt <dsreshiram@gmail.com>
 * @Date: 2023/1/1 11:44
 */

package cqhttp_bot

import (
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Friend struct {
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
}

type SelfInfo struct {
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`

	Company      string `json:"company"`
	Email        string `json:"email"`
	College      string `json:"college"`
	PersonalNote string `json:"personal_note"`
}

type actionRequest struct {
	Action string `json:"action"`
	Params any    `json:"params,omitempty"`
	Echo   string `json:"echo"`
}

type EventMessageType string

const (
	Private = EventMessageType("private")
	Group   = EventMessageType("group")
)

type Msg struct {
	MessageType EventMessageType `json:"message_type"`
	UserId      any              `json:"user_id"`
	GroupId     any              `json:"group_id"`
	Message     string           `json:"message"`
	AutoEscape  *bool            `json:"auto_escape,omitempty"` // 消息内容是否作为纯文本发送 ( 即不解析 CQ 码 ) , 只在 message 字段是字符串时有效
}

type EventMessageContext struct {
	MessageId   int32
	MessageType EventMessageType
	Time        time.Time
	SenderId    int64
	Sender      Sender
	GroupId     int64
	Message     *EventMessage
}
type Sender struct {
	Nickname string
	UserId   int64
	Card     string
}
type EventMessage struct {
	RawMessage string
	Messages   []Message
}
type MessageType uint8

const (
	TEXT = MessageType(iota) + 1
	IMAGE
	REPLY
	AT
	VIDEO
)

type (
	Messages []Message
	Message  struct {
		Type  MessageType
		Image *messageImage
		At    *messageAt
		Text  string
		Video *messageVideo
	}
)

type messageVideo struct {
	File string
	Url  string
}

type messageImage struct {
	Url   string
	File  string
	Src   []byte
	Flash bool
}
type messageAt struct {
	Qid  any
	Name string // 当在群中找不到此QQ号的名称时才会生效
}

type Event struct {
	MessageEvent
	RequestEvent
	NoticeEvent

	callback map[string][]EventCallback
}

type (
	EventCallback func(event string, data jsoniter.Any)
)

type MessageEvent struct {
	onGroupMessages   []OnGroupMessageFunc
	onPrivateMessages []OnPrivateMessageFunc
	onMessage         []OnMessageFunc
}
type RequestEvent struct {
	onFriendRequests []OnFriendRequestFunc
	onGroupRequests  []OnGroupRequestFunc
}
type NoticeEvent struct{}
type (
	OnGroupMessageFunc   func(messageId int32, senderQid, groupId int64, message *EventMessage)
	OnPrivateMessageFunc func(messageId int32, userId int64, message *EventMessage)
	OnMessageFunc        func(ctx *EventMessageContext)
	OnFriendRequestFunc  func(userId int64, comment string, flag string)
	OnGroupRequestFunc   func(userId, groupId int64, requestType GroupRequestEventSubType, comment, flag string)
)

type GroupRequestEventSubType string

const (
	Add    = GroupRequestEventSubType("add")    // 加入
	Invite = GroupRequestEventSubType("invite") // 邀请
)
