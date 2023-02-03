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
 * @Date: 2023/1/1 11:43
 */

package cqhttp_bot

import (
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/panjf2000/ants/v2"
	"log"
)

type Bot struct {
	ws *websocket.Conn
	*Options
	Event
	Action
}

// New 实例化一个Bot对象
func New(addr string, options ...Option) (b *Bot) {
	b = new(Bot)
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatalln("链接websocket错误：", err)
	}
	b.ws = conn
	b.Action.ws = conn
	b.Options = loadOptions(options...)
	b.Action.options = b.Options
	return
}
func (b *Bot) handle() {
	ws := b.ws
	if b.handleThreadNum == 0 {
		b.handleThreadNum = 200
	}
	// 使用goroutine池处理
	h, _ := ants.NewPoolWithFunc(b.handleThreadNum, func(i interface{}) {
		message := i.(jsoniter.Any)
		if ty := message.Get("post_type").ToString(); ty != "" {
			b.event(ty, message)
		} else if m := message.Get("echo").ToString(); m != "" {
			b.actionMap.Store(m, message)
		}
	})
	defer h.Release()
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("信息错误：", err)
		}
		h.Invoke(jsoniter.Get(msg))
	}
}

// Start 已非阻塞的方式运行
func (b *Bot) Start() {
	go b.handle()
}

// Run 已阻塞的方式运行
func (b *Bot) Run() {
	b.handle()
}
