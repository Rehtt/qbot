package cq

import (
	"fmt"
	"testing"
)

func TestMa(t *testing.T) {
	a := CQImage{
		File: "qwe",
		Url:  "ad",
		Type: "fsdf",
	}
	out, _ := MarshalToString(&a)
	fmt.Println(out)
}
