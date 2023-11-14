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
 * @Date: 2023/1/1 12:27
 */

package cqhttp_bot

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// GenCode 生成唯一编码
func GenCode(data []byte) string {
	s := sha256.New()
	s.Write(data)
	s.Write([]byte(time.Now().String()))

	tmp := make([]byte, 20)
	_, _ = rand.Read(tmp)
	s.Write(tmp)

	return base64.StdEncoding.EncodeToString(s.Sum(nil))
}

func parse(raw string) (out map[string]string) {
	out = make(map[string]string)
	for _, r := range strings.Split(raw, ",") {
		s := strings.SplitN(r, "=", 2)
		if len(s) == 2 {
			out[s[0]] = s[1]
		}
	}
	return
}

func DeepCopy(dst, src any) {
	o, _ := jsoniter.Marshal(src)
	_ = jsoniter.Unmarshal(o, dst)
}

func Unescape(str string) string {
	return strings.NewReplacer(
		"&#91;", "[",
		"&#93;", "]",
		"&amp;", "&",
		"&#44;", ",",
	).Replace(str)
}
