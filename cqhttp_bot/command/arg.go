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

package command

import (
	"strconv"
)

type Arg struct {
	name  string
	usage string
	value string
}
type Flag struct {
	args map[string]*Arg
}

func (f *Flag) Var(name, defaultValue, usage string) {
	if f.args == nil {
		f.args = make(map[string]*Arg)
	}
	f.args[name] = &Arg{
		name:  name,
		usage: usage,
		value: defaultValue,
	}
}

func (f Flag) Get(name string) (string, bool) {
	if a, ok := f.args[name]; ok {
		return a.value, true
	} else {
		return "", false
	}
}

func (f Flag) GetInt(name string) (int, bool) {
	if s, ok := f.Get(name); !ok {
		return 0, ok
	} else {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, false
		}
		return n, true
	}
}

func (f Flag) Range(fun func(name, value, usage string)) {
	for _, v := range f.args {
		fun(v.name, v.value, v.usage)
	}
}

func (f *Flag) set(name, value string) {
	if a, ok := f.args[name]; ok {
		a.value = value
	} else {
		f.Var(name, value, "")
	}
}
