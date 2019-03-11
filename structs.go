package structs

import (
	"encoding/json"
	"errors"
	"fmt"

	"reflect"
	"time"
)

//ErrInvalidTag defines the error when the tags are invalids
var ErrInvalidTag = errors.New("Invalid tag")

//TransStructArr Trnasformforms the struct array from struct array  type to dest struct array
func TransStructArr(dest, src interface{}) error {
	destv := reflect.ValueOf(dest)
	if destv.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer %v", destv.Type())
	}
	// get the value that the pointer v points to.
	destv = destv.Elem()
	if destv.Kind() != reflect.Slice {
		return fmt.Errorf("can't fill non-slice value")
	}

	return setSliceValue(destv, reflect.ValueOf(src))
}

//TransStruct Transforms the src struc to the dest struct
func TransStruct(dest, src interface{}) error {
	destv := reflect.ValueOf(dest)
	if destv.Kind() != reflect.Ptr {
		return fmt.Errorf("%v not a pointer", destv.Type())
	}
	// get the value that the pointer v points to.
	destv = destv.Elem()
	if destv.Kind() != reflect.Struct {
		return fmt.Errorf("can't fill non-struct value")
	}
	return transformStruct(destv, reflect.ValueOf(src))
}

//MarshalTransStructArr Marsharlls and transforms the struct array from struct array  type to dest struct array
func MarshalTransStructArr(dest, src interface{}) ([]byte, error) {
	destv := reflect.ValueOf(dest)
	if destv.Kind() == reflect.Ptr {
		destv = destv.Elem()
	} else if destv.Kind() == reflect.Slice {
		destv = reflect.New(destv.Type()).Elem()
	} else {
		return []byte{}, fmt.Errorf("%s not a pointer nor slice", destv.Type())
	}
	err := setSliceValue(destv, reflect.ValueOf(src))
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(destv.Interface())
}

//MarshallTransStruct realiza el marshall de la estructura transformada
func MarshallTransStruct(dest, src interface{}) ([]byte, error) {
	destv := reflect.ValueOf(dest)
	switch reflect.TypeOf(dest).Kind() {
	case reflect.Ptr:
	case reflect.Struct:
		destv = reflect.New(reflect.TypeOf(dest))
	default:
		return []byte{}, fmt.Errorf("%s not a struct nor slice", destv.Type())
	}
	err := transformStruct(destv.Elem(), reflect.ValueOf(src))
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(dest)
}

func transformStruct(destv, srcv reflect.Value) error {
	destType := destv.Type()
	for i := 0; i < destv.NumField(); i++ {
		tag, err := newTag(destType.Field(i))
		if err != nil {
			return err
		}
		//log.Print(" From ", tag.From, " ")
		v := srcv.FieldByName(tag.From)
		err = setValue(destv.Field(i), v)
		if err != nil {
			return err
		}
		//log.Println()
	}
	return nil
}

func setValue(destv, srcv reflect.Value) error {
	if !srcv.IsValid() {
		return nil
	}
	switch destv.Kind() {
	case reflect.String:
		return setString(destv, srcv)
	case reflect.Struct:
		return setStruct(destv, srcv)
	case reflect.Ptr:
		return setPtrValue(destv, srcv)
	case reflect.Slice:
		return setSliceValue(destv, srcv)
	}
	return fmt.Errorf("%s type not implemented", destv.Kind())
}

func setString(destv, srcv reflect.Value) error {
	switch srcv.Kind() {
	case reflect.Ptr:
		return setValue(destv, srcv.Elem())
	}

	switch v := srcv.Interface().(type) {
	case time.Time:
		//log.Printf(" = setting_time[%v] ", srcv)
		destv.SetString(v.Format(time.RFC3339))
		return nil
	}

	//log.Printf(" = setting_%v[%v] ", srcv.Kind(), srcv)
	destv.SetString(fmt.Sprint(srcv))
	return nil
}

func setStruct(destv, srcv reflect.Value) error {
	switch k := srcv.Kind(); k {
	case reflect.Struct:
		return transformStruct(destv, srcv)
	case reflect.Ptr:
		return setValue(destv, srcv.Elem())
	default:
		return fmt.Errorf("Can't set kind %v to struct", k)
	}
}

func setPtrValue(destv, srcv reflect.Value) error {
	switch srcv.Kind() {
	case reflect.Ptr:
		return setValue(destv, srcv.Elem())

	}
	if srcv.IsValid() {
		ptr := reflect.New(destv.Type().Elem())
		err := setValue(ptr.Elem(), srcv)
		if err != nil {
			return err
		}
		destv.Set(ptr)
		return nil
	}
	return nil
}

func setSliceValue(destv, srcv reflect.Value) error {
	switch destv.Type().Elem().Kind() {
	case reflect.Struct:
		return setSliceStruct(destv, srcv)
	default:
		return setSliceString(destv, srcv)
	}
}

func setSliceString(destv, srcv reflect.Value) error {
	//log.Print("setSliceString ")
	l := srcv.Len()
	if l == 0 {
		return nil
	}
	s := make([]string, l)
	for i := 0; i < l; i++ {
		s[i] = fmt.Sprint(srcv.Index(i))
	}
	destv.Set(reflect.ValueOf(s))
	return nil
}
func setSliceStruct(destv, srcv reflect.Value) error {
	//log.Print("setSliceStruct ")
	l := srcv.Len()
	if l == 0 {
		return nil
	}
	s := reflect.MakeSlice(reflect.SliceOf(destv.Type().Elem()), l, l)
	for i := 0; i < l; i++ {
		err := transformStruct(s.Index(i), srcv.Index(i))
		if err != nil {
			return err
		}
	}
	destv.Set(s)
	return nil
}
