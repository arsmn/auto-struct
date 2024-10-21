package autostruct

import (
	"fmt"
	"reflect"
	"strconv"
)

const tagName = "default"

type SetterFunc func(reflect.Value, string) error

func getSetterFunc(k reflect.Kind) SetterFunc {
	switch k {
	case reflect.Bool:
		return BoolSetter
	case reflect.String:
		return StringSetter
	case reflect.Struct:
		return StructSetter
	case reflect.Int:
		return IntSetter
	case reflect.Int8:
		return Int8Setter
	case reflect.Int16:
		return Int16Setter
	case reflect.Int32:
		return Int32Setter
	case reflect.Int64:
		return Int64Setter
	case reflect.Uint:
		return UintSetter
	case reflect.Uint8:
		return Uint8Setter
	case reflect.Uint16:
		return Uint16Setter
	case reflect.Uint32:
		return Uint32Setter
	case reflect.Uint64:
		return Uint64Setter
	case reflect.Float32:
		return Float32Setter
	case reflect.Float64:
		return Float64Setter
	case reflect.Complex64:
		return Complex64Setter
	case reflect.Complex128:
		return Complex128Setter
	case reflect.Pointer:
		return PointerSetter
	default:
		return nil
	}
}

func BoolSetter(v reflect.Value, tag string) error {
	if v.Kind() != reflect.Bool {
		return fmt.Errorf("BoolSetter does not support this type: %s", v.Kind())
	}

	b, err := strconv.ParseBool(tag)
	if err != nil {
		return err
	}

	v.SetBool(b)

	return nil
}

func StringSetter(v reflect.Value, tag string) error {
	if v.Kind() != reflect.String {
		return fmt.Errorf("StringSetter does not support this type: %s", v.Kind())
	}

	v.SetString(tag)

	return nil
}

func IntSetter(v reflect.Value, tag string) error {
	return intSetter(strconv.IntSize, v, tag)
}

func Int8Setter(v reflect.Value, tag string) error {
	return intSetter(8, v, tag)
}

func Int16Setter(v reflect.Value, tag string) error {
	return intSetter(16, v, tag)
}

func Int32Setter(v reflect.Value, tag string) error {
	return intSetter(32, v, tag)
}

func Int64Setter(v reflect.Value, tag string) error {
	return intSetter(64, v, tag)
}

func intSetter(bitSize int, v reflect.Value, tag string) error {
	if !v.CanInt() {
		return fmt.Errorf("Int%dSetter does not support this type: %s", bitSize, v.Kind())
	}

	i, err := strconv.ParseInt(tag, 10, bitSize)
	if err != nil {
		return err
	}

	v.SetInt(i)

	return nil
}

func UintSetter(v reflect.Value, tag string) error {
	return uintSetter(strconv.IntSize, v, tag)
}

func Uint8Setter(v reflect.Value, tag string) error {
	return uintSetter(8, v, tag)
}

func Uint16Setter(v reflect.Value, tag string) error {
	return uintSetter(16, v, tag)
}

func Uint32Setter(v reflect.Value, tag string) error {
	return uintSetter(32, v, tag)
}

func Uint64Setter(v reflect.Value, tag string) error {
	return uintSetter(64, v, tag)
}

func uintSetter(bitSize int, v reflect.Value, tag string) error {
	if !v.CanUint() {
		return fmt.Errorf("Uint%dSetter does not support this type: %s", bitSize, v.Kind())
	}

	i, err := strconv.ParseUint(tag, 10, bitSize)
	if err != nil {
		return err
	}

	v.SetUint(i)

	return nil
}

func Float32Setter(v reflect.Value, tag string) error {
	return floatSetter(32, v, tag)
}

func Float64Setter(v reflect.Value, tag string) error {
	return floatSetter(64, v, tag)
}

func floatSetter(bitSize int, v reflect.Value, tag string) error {
	if !v.CanFloat() {
		return fmt.Errorf("Float%dSetter does not support this type: %s", bitSize, v.Kind())
	}

	f, err := strconv.ParseFloat(tag, bitSize)
	if err != nil {
		return err
	}

	v.SetFloat(f)

	return nil
}

func Complex64Setter(v reflect.Value, tag string) error {
	return complexSetter(64, v, tag)
}

func Complex128Setter(v reflect.Value, tag string) error {
	return complexSetter(128, v, tag)
}

func complexSetter(bitSize int, v reflect.Value, tag string) error {
	if !v.CanComplex() {
		return fmt.Errorf("Complex%dSetter does not support this type: %s", bitSize, v.Kind())
	}

	c, err := strconv.ParseComplex(tag, bitSize)
	if err != nil {
		return err
	}

	v.SetComplex(c)

	return nil
}

func PointerSetter(v reflect.Value, tag string) error {
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("PointerSetter does not support this type: %s", v.Kind())
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	return setFieldValue(v.Elem(), tag)
}

func StructSetter(v reflect.Value, tag string) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("StructSetter does not support this type: %s", v.Kind())
	}

	if tag != "nested" {
		return nil
	}

	return structSetter(v)
}

func structSetter(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if err := setFieldValue(v.Field(i), t.Field(i).Tag.Get(tagName)); err != nil {
			return err
		}
	}

	return nil
}

func setFieldValue(v reflect.Value, tag string) error {
	if tag == "" {
		return nil
	}

	if !v.CanSet() {
		return fmt.Errorf("[%s] is not exported", v)
	}

	fn := getSetterFunc(v.Kind())
	if fn == nil {
		return fmt.Errorf("[%s] type is not supported", v.Kind())
	}

	return fn(v, tag)
}

func Set(obj any) error {
	rv := dereference(obj)

	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("[%s] type is not supported. must be struct", rv.Kind())
	}

	return structSetter(rv)
}

func dereference(obj any) reflect.Value {
	rv := reflect.ValueOf(obj)

	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	return rv
}

func MustSet(v any) {
	if err := Set(v); err != nil {
		panic(err)
	}
}

func New[T any]() T {
	var v T
	MustSet(&v)
	return v
}
