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

type Friend struct {
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
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
	UserId      int64            `json:"user_id"`
	GroupId     int64            `json:"group_id"`
	Message     string           `json:"message"`
	AutoEscape  *bool            `json:"auto_escape,omitempty"` // 消息内容是否作为纯文本发送 ( 即不解析 CQ 码 ) , 只在 message 字段是字符串时有效
}

type EventMessage struct {
	RawMessage string
	Messages   []Message
}
type MessageType uint8

const (
	TEXT = MessageType(iota)
	IMAGE
)

type Message struct {
	Type MessageType
	Url  string
	File string
	Text string
}
