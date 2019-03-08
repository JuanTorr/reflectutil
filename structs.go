package structs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//ErrInvalidTag defines the error when the tags are invalids
var ErrInvalidTag = errors.New("Invalid tag")

const tagName = "trans"

type tag struct {
	From string
}

func newTag(f reflect.StructField) (t tag, err error) {
	s := f.Tag.Get(tagName)
	if s == "" {
		return tag{From: f.Name}, nil
	}
	fields := strings.Split(s, ",")
	for _, f := range fields {
		fv := strings.Split(f, ":") //field-value
		if len(fv) != 2 {
			return t, ErrInvalidTag
		}
		switch fv[0] {
		case "from":
			t.From = fv[1]
		}
	}
	return t, nil
}

func TransformStruct(dest, src interface{}) error {
	return transformStruct(reflect.Indirect(reflect.ValueOf(dest)), reflect.ValueOf(src))
}

func transformStruct(destv, srcv reflect.Value) error {
	destType := destv.Type()
	for i := 0; i < destv.NumField(); i++ {
		tag, err := newTag(destType.Field(i))
		if err != nil {
			return err
		}
		fmt.Print(" From ", tag.From, " :")
		v := srcv.FieldByName(tag.From)
		if v.IsValid() {
			fmt.Print("v:", v, " ")
			setValue(destv.Field(i), v)
		}
		fmt.Println()
	}
	return nil
}

func setValue(dest, src reflect.Value) {
	fmt.Print("(", dest.Kind(), " ", ") ")
	switch dest.Kind() {
	case reflect.String:
		dest.SetString(fmt.Sprint(src))
	case reflect.Struct:
	case reflect.Ptr:
	}
}
