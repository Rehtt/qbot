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
	"log/slog"
	"net/http"
)

type Options struct {
	handleThreadNum int
	log             *slog.Logger
	requestHead     http.Header
}
type Option func(options *Options)

// WithHandleThreadNum 默认为200，如果数值太小会导致处理出现延时
func WithHandleThreadNum(n int) Option {
	return func(options *Options) {
		options.handleThreadNum = n
	}
}

func WithLogger(l *slog.Logger) Option {
	return func(options *Options) {
		options.log = l
	}
}

func WithRequestHead(h http.Header) Option {
	return func(options *Options) {
		options.requestHead = h.Clone()
	}
}

func loadOptions(options ...Option) *Options {
	o := new(Options)
	for _, opt := range options {
		opt(o)
	}
	return o
}

func (o *Options) Log() *slog.Logger {
	if o.log == nil {
		o.log = slog.Default()
	}
	return o.log
}
