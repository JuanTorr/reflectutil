package structs

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func randString() string {
	n := rand.Intn(65)
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

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
	NotStr    string `trans:"from:Str"`
	NotNum    string `trans:"from:Num"`
	SameName  string
	NilPtr    string
	Nested    NestedDest
	Time      string
	Ptr       *NestedDest `trans:"from:Nested"`
	StringArr []string
	Recursive *Dest
	StructArr []NestedDest
}
type NestedDest struct {
	ValA string `trans:"from:A"`
	ValB string `trans:"from:B"`
	Same string
}

func getSrc() Src {
	i := rand.Int()
	a := Src{
		Str:       randString(),
		SameName:  rune(rand.Int()),
		Nested:    NestedSrc{A: rand.Int(), B: rand.Intn(2) == 1, Same: randString()},
		Time:      time.Now(),
		StringArr: []string{randString(), randString(), randString()},
		StructArr: []NestedSrc{NestedSrc{Same: randString()}, NestedSrc{Same: randString()}},
	}
	if rand.Intn(2) == 1 {
		a.Num = &i
	}
	if rand.Intn(2) == 1 {
		i = rand.Int()
		a.Recursive = &Src{
			Str:       randString(),
			SameName:  rune(rand.Int()),
			Nested:    NestedSrc{A: rand.Int(), B: rand.Intn(2) == 1, Same: randString()},
			Time:      time.Now(),
			StringArr: []string{randString(), randString(), randString()},
		}
		if rand.Intn(2) == 1 {
			a.Recursive.Num = &i
		}
	}
	i = rand.Intn(2)
	if i == 0 {
		i = rand.Intn(5)
		a.StructArr = make([]NestedSrc, i)
		for j := 0; j < i; j++ {
			a.StructArr[j] = NestedSrc{A: rand.Int(), B: rand.Intn(2) == 1, Same: randString()}
		}
	}

	return a
}

func compareStructValues(dest Dest, src Src) (where string, ok bool) {
	ok = true
	compareok := func(currentOK bool, a string) {
		if !currentOK {
			where += fmt.Sprint(":", a)
			ok = false
		}
	}
	compareok(dest.NotStr == src.Str, `dest.NotStr == src.Str`)
	compareok(dest.SameName == fmt.Sprint(src.SameName), `dest.SameName == fmt.Sprint(src.SameName)`)
	compareok(compareNested(dest.Nested, src.Nested), `compareNested(dest.Nested,src.Nested)`)
	compareok(dest.Time == src.Time.Format(time.RFC3339), `dest.Time == src.Time.Format(time.RFC3339)`)
	if src.Num == nil {
		compareok(dest.NotNum == "", `dest.NotNum == ""`)
	} else {
		compareok(dest.NotNum == fmt.Sprint(*src.Num), `dest.NotNum == fmt.Sprint(src.Num)`)
	}
	if src.NilPtr == nil {
		compareok(dest.NilPtr == "", `dest.NilPtr == ""`)
	} else {
		compareok(dest.NilPtr == fmt.Sprint(*src.NilPtr), `dest.NilPtr == fmt.Sprint(src.NilPtr)`)
	}
	if dest.Recursive != nil {
		if src.Recursive == nil {
			compareok(false, `src.Recursive == nil`)
			return
		}
		where, curOk := compareStructValues(*dest.Recursive, *src.Recursive)
		compareok(curOk, "\n"+where)
	}
	compareok(compareArrString(dest.StringArr, src.StringArr), `compareArrString(dest.StringArr, src.StringArr)`)
	if len(dest.StructArr) != len(src.StructArr) {
		return where + ":" + "len struct Arr", false
	}
	for i, v := range dest.StructArr {
		compareok(compareNested(v, src.StructArr[i]), `compareNested(dest.Nested,src.Nested)`)
	}
	return
}

func compareNested(dest NestedDest, src NestedSrc) bool {
	return dest.ValA == fmt.Sprint(src.A) && dest.ValB == fmt.Sprint(src.B) && dest.Same == src.Same
}

func compareArrString(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
func TestTransStructArr(t *testing.T) {
	for i := 0; i < 100; i++ {
		var d Dest
		s := getSrc()
		err := TransStruct(&d, s)
		if err != nil {
			t.Error(err)
			return
		}
		where, ok := compareStructValues(d, s)
		if !ok {
			t.Errorf("Incorrect transformation:%s\n\n%+v\n%+v", where, s, d)
			return
		}
	}
}
func TestTransStruct(t *testing.T) {
	i := rand.Intn(100) + 100
	a := make([]Src, i, i)
	for i := range a {
		a[i] = getSrc()
	}
	var b []Dest
	TransStructArr(&b, a)
	if len(b) != len(a) {
		t.Error("trans struct len not equal", len(b), len(a))
	}
	for i, v := range b {
		where, ok := compareStructValues(v, a[i])
		if !ok {
			t.Errorf("Incorrect transformation:%s\n\n%+v\n%+v", where, a, b)
			return
		}
	}
}
func TestMarshalTransStructArr(t *testing.T) {
	i := rand.Intn(100) + 100
	a := make([]Src, i, i)
	for i := range a {
		a[i] = getSrc()
	}
	_, err := MarshalTransStructArr([]Dest{}, a)

	if err != nil {
		t.Error(err)
		return
	}
}
func TestMarshalTransStructArrPtr(t *testing.T) {
	i := rand.Intn(100) + 100
	a := make([]Src, i, i)
	for i := range a {
		a[i] = getSrc()
	}
	var b []Dest
	_, err := MarshalTransStructArr(&b, a)

	if err != nil {
		t.Error(err)
		return
	}
}
func TestMarshallTransStruct(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := getSrc()
		_, err := MarshallTransStruct(Dest{}, s)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
func TestMarshallTransStructPtr(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := getSrc()
		_, err := MarshallTransStruct(Dest{}, s)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestEmptyArr(t *testing.T) {
	var src []Src
	var dest []Dest
	err := TransStructArr(&dest, src)
	if err != nil {
		t.Error(err)
		return
	}
	if dest != nil {
		t.Error("nil array should result in nil arrar")
		return
	}

	src = []Src{}
	err = TransStructArr(&dest, src)
	if err != nil {
		t.Error(err)
		return
	}
	if dest == nil || len(dest) != 0 {
		t.Error("Empty array should result in empty arrar")
		return
	}
}
