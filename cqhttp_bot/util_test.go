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
)

func TestGenCode(t *testing.T) {
	t.Log(GenCode([]byte("123")))
	if GenCode([]byte("123")) == GenCode([]byte("123")) {
		t.Error("重复")
	} else {
		t.Log("通过")
	}
}

func TestJ(t *testing.T) {
	a := "123[CQ:image,file=70c441608c8497b9eb883a5cc0c4cf8a.image,subType=0,url=https://gchat.qpic.cn/gchatpic_new/271304801/860088749-2988409058-70C441608C8497B9EB883A5CC0C4CF8A/0?term=3&amp;is_origin=1]131"
	fmt.Println(ParseMessage(a))
}
