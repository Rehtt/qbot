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
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	b := New("ws://127.0.0.1:8080")
	b.Start()
	fmt.Println(b.GetFriendsList())
	b.Event.OnGroupMessage(func(messageId int32, senderQid, groupId int64, message *EventMessage) {
		fmt.Println(senderQid, groupId, message.Messages, message.RawMessage)
	})
	b.OnPrivateMessage(func(messageId int32, userId int64, message *EventMessage) {
		fmt.Println(userId, message.Messages, message.RawMessage)
	})
	b.SendMsg(852122585, MessageArray(TextMessage("test")), Group)
	time.Sleep(10 * time.Minute)
}

func BenchmarkNew(b *testing.B) {
	b.StopTimer()
	n := New("ws://cqhttp.rehtt.com/ws/")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		n.GetFriendsList()
	}
}
