package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
)

//ErrInvalidTag defines the error when the tags are invalids
var ErrInvalidTag = errors.New("Invalid tag")

//TransformStructArr Trnasformforms the struct array from src type to dest type
func TransformStructArr(dest, src interface{}) {
	//TODO:
}

//MarshallTransformStruct realiza el marshall de la estructura transformada
func MarshallTransformStruct(dest, src interface{}) ([]byte, error) {

	switch reflect.TypeOf(dest).Kind() {
	case reflect.Ptr:
		TransformStruct(dest, src)
		return json.Marshal(dest)
	case reflect.Struct:
		destv := reflect.New(reflect.TypeOf(dest))
		transformStruct(destv.Elem(), reflect.ValueOf(src))
		return json.Marshal(destv.Interface())
	default:
		panic("not pointer or struct")
	}
}

//TransformStruct Transforms the src struc to the dest struct
func TransformStruct(dest, src interface{}) error {
	return transformStruct(reflect.ValueOf(dest).Elem(), reflect.ValueOf(src))
}

func transformStructArr(destv, srcv reflect.Value) {
	//TODO:
}

func transformStruct(destv, srcv reflect.Value) error {
	destType := destv.Type()
	for i := 0; i < destv.NumField(); i++ {
		tag, err := newTag(destType.Field(i))
		if err != nil {
			return err
		}
		log.Print(" From ", tag.From, " ")
		v := srcv.FieldByName(tag.From)
		setValue(destv.Field(i), v)
		log.Println()
	}
	return nil
}

func setValue(destv, srcv reflect.Value) {
	if !srcv.IsValid() {
		return
	}
	switch destv.Kind() {
	case reflect.String:
		setString(destv, srcv)
	case reflect.Struct:
		setStruct(destv, srcv)
	case reflect.Ptr:
		setPtrValue(destv, srcv)
	case reflect.Slice:
		setSliceValue(destv, srcv)
	}
}

func setString(destv, srcv reflect.Value) {
	switch srcv.Kind() {
	case reflect.Ptr:
		setValue(destv, srcv.Elem())
		return
	}

	switch v := srcv.Interface().(type) {
	case time.Time:
		log.Printf(" = setting_time[%v] ", srcv)
		destv.SetString(v.Format(time.RFC3339))
		return
	}

	log.Printf(" = setting_%v[%v] ", srcv.Kind(), srcv)
	destv.SetString(fmt.Sprint(srcv))
}

func setStruct(destv, srcv reflect.Value) {
	switch k := srcv.Kind(); k {
	case reflect.Struct:
		transformStruct(destv, srcv)
	case reflect.Ptr:
		setValue(destv, srcv.Elem())
	default:
		panic(fmt.Sprintf("Can't set kind %v to struct", k))
	}
}

func setPtrValue(destv, srcv reflect.Value) {
	switch srcv.Kind() {
	case reflect.Ptr:
		setValue(destv, srcv.Elem())
		return
	}
	if srcv.IsValid() {
		ptr := reflect.New(destv.Type().Elem())
		setValue(ptr.Elem(), srcv)
		destv.Set(ptr)
	}
}

func setSliceValue(destv, srcv reflect.Value) {
	switch destv.Type().Elem().Kind() {
	case reflect.Struct:
		setSliceStruct(destv, srcv)
	default:
		setSliceString(destv, srcv)
	}
}

func setSliceString(destv, srcv reflect.Value) {
	log.Print("setSliceString ")
	l := srcv.Len()
	if l == 0 {
		return
	}
	s := make([]string, l)
	for i := 0; i < l; i++ {
		s[i] = fmt.Sprint(srcv.Index(i))
	}
	destv.Set(reflect.ValueOf(s))
}
func setSliceStruct(destv, srcv reflect.Value) {
	log.Print("setSliceStruct ")
	l := srcv.Len()
	if l == 0 {
		return
	}
	s := reflect.MakeSlice(reflect.SliceOf(destv.Type().Elem()), l, l)
	for i := 0; i < l; i++ {
		transformStruct(s.Index(i), srcv.Index(i))
	}
	destv.Set(s)
}
