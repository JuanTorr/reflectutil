package structs

import (
	"fmt"
	"testing"
	"time"
)

type Src struct {
	Str       string
	Num       *int
	SameName  rune
	NilPtr    *int
	Nested    NestedSrc
	Time      time.Time
	Recursive *Src
	StringArr []string
	StructArr []NestedSrc
}

type NestedSrc struct {
	A    int
	B    bool
	Same string
}

type Dest struct {
	TestString  string `trans:"from:Str"`
	TestString2 string `trans:"from:Num"`
	SameName    string
	NilPtr      string
	Nested      NestedDest
	Time        string
	Ptr         *NestedDest `trans:"from:Nested"`
	StringArr   []string
	Recursive   *Dest
	StructArr   []NestedDest
}
type NestedDest struct {
	ValA string `trans:"from:A"`
	ValB string `trans:"from:B"`
	Same string
}

func TestTransformStruct(t *testing.T) {
	/* i := 5
	a := Src{
		Str:       "hello",
		Num:       &i,
		SameName:  'a',
		Nested:    NestedSrc{A: 4, B: true, Same: "hi"},
		Time:      time.Now(),
		StringArr: []string{"a", "b", "c"},
		StructArr: []NestedSrc{NestedSrc{Same: "x"}, NestedSrc{Same: "y"}},
	}
	i = 6
	a.Recursive = &Src{
		Str:       "world",
		Num:       &i,
		SameName:  'b',
		Nested:    NestedSrc{A: 5, B: true, Same: "ho"},
		Time:      time.Now(),
		StringArr: []string{"1", "2", "3"},
	}
	b := Dest{}

	err := TransformStruct(&b, a)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("\n\nRESULT: %+v\n", b)
		if b.Recursive != nil {
			fmt.Println(*b.Recursive)
		}
		if b.Ptr != nil {
			fmt.Println(*b.Ptr)
		}
		fmt.Println()
	}
	*/
}

func TestMarshalStruct(t *testing.T) {
	type Nested struct {
		Delta   bool
		Epsilon float64
	}
	a := struct {
		Alfa  string
		Beta  int
		Gamma Nested
	}{"Hello", 5, Nested{false, 64}}

	j, err := MarshallTransformStruct(struct {
		A string `trans:"from:Alfa" json:"tagAlfa"`
		B string `trans:"from:Beta" json:"tagBeta"`
		C struct {
			D string `trans:"from:Delta" json:"tagDelta"`
			E string `trans:"from:Epsilon" json:"tagEpsilon"`
		} `trans:"from:Gamma" json:"tagGamma"`
	}{}, a)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(j))

}

/* func BenchmarkTransformStructB(b *testing.B) {
	i := 5
	a := Src{Str: "hello", Num: &i, SameName: 'a', Nested: NestedSrc{A: 4, B: true, Same: "hi"}, Time: time.Now()}
	i = 6
	a.Recursive = &Src{Str: "world", Num: &i, SameName: 'b', Nested: NestedSrc{A: 5, B: true, Same: "ho"}, Time: time.Now()}
	for i := 0; i < b.N; i++ {
		b := Dest{}
		TransformStruct(&b, a)
	}
}
*/
