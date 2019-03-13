package reflectutil

import (
	"reflect"
	"strings"
)

type tag struct {
	From string
}

func newTag(f reflect.StructField) (t tag, err error) {
	s := f.Tag.Get("trans")
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
