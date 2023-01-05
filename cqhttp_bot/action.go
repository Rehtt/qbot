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
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Action struct {
	ws          *websocket.Conn
	actionMap   sync.Map
	actionIndex atomic.Int64
}

func (b *Action) action(action string, data any) (jsoniter.Any, error) {
	request := actionRequest{
		Action: action,
		Params: data,
	}
	var tmp bytes.Buffer
	jsoniter.NewEncoder(&tmp).Encode(request)
	tmp.WriteString(strconv.FormatInt(b.actionIndex.Add(1), 10))
	request.Echo = GenCode(tmp.Bytes())
	err := b.ws.WriteJSON(request)
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
func (b *Action) GetFriendsList() ([]*Friend, error) {
	data, err := b.action("get_friend_list", nil)
	if err != nil {
		return nil, err
	}
	var fs []*Friend
	data.ToVal(&fs)
	return fs, nil
}

// SendMsg 发送消息
func (b *Action) SendMsg(qid int64, msg string, ty EventMessageType, autoEscape ...bool) error {
	m := Msg{
		MessageType: ty,
		Message:     msg,
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
	_, err := b.action("send_msg", m)
	return err
}

func (b *Action) SendPrivateMsg(userId int64, msg string) error {
	return b.SendMsg(userId, msg, Private)
}
