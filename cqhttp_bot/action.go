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
	"bytes"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type Action struct {
	ws          *websocket.Conn
	actionMap   sync.Map
	actionIndex atomic.Int64
	sendLock    sync.Mutex
	options     *Options
	cache       sync.Map
}

func (b *Action) Send(action string, data any) (jsoniter.Any, error) {
	return b.action(action, data)
}

func (b *Action) action(action string, data any) (jsoniter.Any, error) {
	request := actionRequest{
		Action: action,
		Params: data,
	}
	var tmp bytes.Buffer
	_ = jsoniter.NewEncoder(&tmp).Encode(request)
	tmp.WriteString(strconv.FormatInt(b.actionIndex.Add(1), 10))
	request.Echo = GenCode(tmp.Bytes())
	out, _ := jsoniter.Marshal(request)
	b.options.Log().Debug("qbot action", "msg", string(out))
	// 不允许并发发送
	b.sendLock.Lock()
	err := b.ws.WriteMessage(websocket.TextMessage, out)
	b.sendLock.Unlock()
	if err != nil {
		return nil, err
	}
	for {
		time.Sleep(50 * time.Millisecond)
		if d, ok := b.actionMap.LoadAndDelete(request.Echo); ok {
			response := d.(jsoniter.Any)
			if response.Get("status").ToString() != "ok" {
				return nil, errors.New(response.Get("wording").ToString())
			}
			return response.Get("data"), nil
		}
	}
}

// GetFriendsList 获取好友列表
func (b *Action) GetFriendsList() ([]Friend, error) {
	data, err := b.action("get_friend_list", nil)
	if err != nil {
		return nil, err
	}
	var fs []Friend
	data.ToVal(&fs)
	return fs, nil
}

// SendMsg 发送消息
//
// @param qid	目标，私聊为对方qq号，群聊为群号
// @param msg	消息
// @param ty	消息类型，私聊为Private，群聊为Group
//
// @return id	消息id
// @return error
func (b *Action) SendMsg(qid any, msg Messages, ty EventMessageType, autoEscape ...bool) (int32, error) {
	m := Msg{
		MessageType: ty,
		Message:     msg.RawMessage(),
	}
	switch ty {
	case Private:
		m.UserId = qid
	case Group:
		m.GroupId = qid
	}
	if len(autoEscape) > 0 {
		m.AutoEscape = &autoEscape[0]
	}
	response, err := b.action("send_msg", m)
	if err != nil {
		return 0, err
	}
	return response.Get("message_id").ToInt32(), nil
}

func (b *Action) SendPrivateMsg(userId int64, msg string) (int32, error) {
	return b.SendMsg(userId, MessageArray(TextMessage(msg)), Private)
}

// GetMsg 根据message_id获取消息
func (b *Action) GetMsg(messageId string) (message *Message, err error) {
	data, err := b.action("get_msg", map[string]any{
		"message_id": messageId,
	})
	if err != nil {
		return nil, err
	}
	return &ParseMessage(data.Get("message").ToString())[0], nil
}

func (b *Action) SetFriendAddRequest(flag string, approve bool, remark ...string) error {
	tmp := struct {
		Flag    string `json:"flag"`
		Approve bool   `json:"approve"`
		Remark  string `json:"remark,omitempty"` // 备注
	}{
		Flag:    flag,
		Approve: approve,
	}
	if len(remark) != 0 {
		tmp.Remark = remark[0]
	}
	_, err := b.action("set_friend_add_request", tmp)
	return err
}

func (b *Action) SetGroupAddRequest(flag string, subType GroupRequestEventSubType, approve bool, reason ...string) error {
	tmp := struct {
		Flag    string                   `json:"flag"`
		Approve bool                     `json:"approve"`
		Reason  string                   `json:"reason,omitempty"` // 拒绝理由
		SubType GroupRequestEventSubType `json:"sub_type"`
	}{
		Flag:    flag,
		Approve: approve,
		SubType: subType,
	}
	if len(reason) != 0 {
		tmp.Reason = reason[0]
	}
	_, err := b.action("set_group_add_request", tmp)
	return err
}

func (b *Action) GetSelfInfo() (*SelfInfo, error) {
	userId, ok := b.cache.Load("user_id")
	if !ok {
		loginInfo, err := b.action("get_login_info", nil)
		if err != nil {
			return nil, err
		}
		userId = loginInfo.Get("user_id").ToInt64()
	}
	info, err := b.action("set_qq_profile", nil)
	if err != nil {
		return nil, err
	}
	return &SelfInfo{
		UserId:       userId.(int64),
		Nickname:     info.Get("nick_name").ToString(),
		Company:      info.Get("company").ToString(),
		Email:        info.Get("email").ToString(),
		College:      info.Get("college").ToString(),
		PersonalNote: info.Get("personal_note").ToString(),
	}, nil
}
