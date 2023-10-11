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

package cq

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	r_strings "github.com/Rehtt/Kit/strings"
)

const CQHeadr = "[CQ:"
const (
	cqImage = "image"
)

type CQImage struct {
	File string `cq:"file,omitempty"`
	Url  string `cq:"url"`
	Type string `cq:"type"`
}

func (CQImage) Name() string {
	return cqImage
}

func (c *CQImage) parse(data []string) error {
	return nil
}

type CQTest struct{}

type CQStruct interface {
	// CQImage | *CQImage
	// CQTest | *CQTest
	Name() string
	parse(data []string) error
}

func Marshal(v CQStruct) ([]byte, error) {
	value := reflect.ValueOf(v)
	ty := value.Type()
	for ty.Kind() == reflect.Ptr {
		value = value.Elem()
		ty = value.Type()
	}
	tmp := make([]string, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		subValue := value.Field(i)
		subType := ty.Field(i)

		if tag, ok := subType.Tag.Lookup("cq"); ok {
			split := strings.Split(tag, ",")
			if !subValue.IsValid() && strings.HasSuffix(tag, "omitempty") {
				tmp = append(tmp, split[0])
				continue
			}

			tmp = append(tmp, fmt.Sprintf("%s=%v", split[0], subValue.Interface()))
		}
	}

	return r_strings.ToBytes(fmt.Sprintf("%s%s,%s]", CQHeadr, v.Name(), strings.Join(tmp, ","))), nil
}

func MarshalToString(v CQStruct) (string, error) {
	out, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return r_strings.ToString(out), nil
}

func Parse(data []byte) (CQStruct, error) {
	if !bytes.HasPrefix(data, r_strings.ToBytes(CQHeadr)) {
		return nil, ErrorNotCQCode
	}
	split := strings.Split(r_strings.ToString(data), ",")
	switch split[0] {
	case cqImage:
		image := new(CQImage)
		err := image.parse(split[1:])
		return image, err
	}

	return nil, ErrorUnknownCode
}
