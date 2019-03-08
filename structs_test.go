package structs

import (
	"fmt"
	"testing"
)

type Src struct {
	Str       string
	Num       int
	SameName  rune
	Recursive *Src
}

type Dest struct {
	TestString  string `trans:"from:Str"`
	TestString2 string `trans:"from:Num"`
	SameName    string
	Recursive   *Dest
}

func TestTransformStruct(t *testing.T) {
	a := Src{Str: "hello", Num: 1, SameName: 'R'}
	b := Dest{}
	err := TransformStruct(&b, a)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("\n\nRESULT: %+v\n", b)
	}
}
