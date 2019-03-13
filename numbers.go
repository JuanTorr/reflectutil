package reflectutil

import (
	"fmt"
	"strconv"
	"strings"
)

//IToBool convert the interface value to bool
func IToBool(i interface{}) (res bool, err error) {
	switch v := i.(type) {
	case bool:
		return v, nil
	case string:
		return strings.ToLower(v) == "true", nil
	default:
		f, err := IToFloat64(v)
		return f != 0, err
	}
}

//IToUint64 convert the interface value to uint64
func IToUint64(i interface{}) (res uint64, err error) {
	switch v := i.(type) {
	case uint64:
		return v, nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		f, err := strconv.ParseFloat(fmt.Sprint(i), 64)
		return uint64(f), err
	}
}

//IToInt64 convert the interface value to int64
func IToInt64(i interface{}) (res int64, err error) {
	switch v := i.(type) {
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		f, err := strconv.ParseFloat(fmt.Sprint(i), 64)
		return int64(f), err
	}
}

//IToFloat64 convert the interface value to float64
func IToFloat64(i interface{}) (res float64, err error) {
	switch v := i.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return strconv.ParseFloat(fmt.Sprint(i), 64)
	}
}
