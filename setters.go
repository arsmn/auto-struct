package autostruct

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type setterFunc func(*config, reflect.Value, string) error

func getSetterFunc(k reflect.Kind) setterFunc {
	switch k {
	case reflect.Bool:
		return boolSetter
	case reflect.String:
		return stringSetter
	case reflect.Struct:
		return structSetter
	case reflect.Int:
		return int0Setter
	case reflect.Int8:
		return int8Setter
	case reflect.Int16:
		return int16Setter
	case reflect.Int32:
		return int32Setter
	case reflect.Int64:
		return int64Setter
	case reflect.Uint:
		return uint0Setter
	case reflect.Uint8:
		return uint8Setter
	case reflect.Uint16:
		return uint16Setter
	case reflect.Uint32:
		return uint32Setter
	case reflect.Uint64:
		return uint64Setter
	case reflect.Float32:
		return float32Setter
	case reflect.Float64:
		return float64Setter
	case reflect.Complex64:
		return complex64Setter
	case reflect.Complex128:
		return complex128Setter
	case reflect.Pointer:
		return pointerSetter
	case reflect.Array:
		return arraySetter
	default:
		return nil
	}
}

func boolSetter(cfg *config, v reflect.Value, tag string) error {
	if v.Kind() != reflect.Bool {
		return fmt.Errorf("BoolSetter does not support [%s]", v.Kind())
	}

	b, err := strconv.ParseBool(tag)
	if err != nil {
		return err
	}

	v.SetBool(b)

	return nil
}

func stringSetter(cfg *config, v reflect.Value, tag string) error {
	if v.Kind() != reflect.String {
		return fmt.Errorf("StringSetter does not support [%s]", v.Kind())
	}

	v.SetString(tag)

	return nil
}

func int0Setter(cfg *config, v reflect.Value, tag string) error {
	return intSetter(cfg, v, tag, 0)
}

func int8Setter(cfg *config, v reflect.Value, tag string) error {
	return intSetter(cfg, v, tag, 8)
}

func int16Setter(cfg *config, v reflect.Value, tag string) error {
	return intSetter(cfg, v, tag, 16)
}

func int32Setter(cfg *config, v reflect.Value, tag string) error {
	return intSetter(cfg, v, tag, 32)
}

func int64Setter(cfg *config, v reflect.Value, tag string) error {
	return intSetter(cfg, v, tag, 64)
}

func intSetter(_ *config, v reflect.Value, tag string, bitSize int) error {
	if !v.CanInt() {
		return fmt.Errorf("Int%dSetter does not support [%s]", bitSize, v.Kind())
	}

	i, err := strconv.ParseInt(tag, 10, bitSize)
	if err != nil {
		return err
	}

	v.SetInt(i)

	return nil
}

func uint0Setter(cfg *config, v reflect.Value, tag string) error {
	return uintSetter(cfg, v, tag, 0)
}

func uint8Setter(cfg *config, v reflect.Value, tag string) error {
	return uintSetter(cfg, v, tag, 8)
}

func uint16Setter(cfg *config, v reflect.Value, tag string) error {
	return uintSetter(cfg, v, tag, 16)
}

func uint32Setter(cfg *config, v reflect.Value, tag string) error {
	return uintSetter(cfg, v, tag, 32)
}

func uint64Setter(cfg *config, v reflect.Value, tag string) error {
	return uintSetter(cfg, v, tag, 64)
}

func uintSetter(_ *config, v reflect.Value, tag string, bitSize int) error {
	if !v.CanUint() {
		return fmt.Errorf("Uint%dSetter does not support [%s]", bitSize, v.Kind())
	}

	i, err := strconv.ParseUint(tag, 10, bitSize)
	if err != nil {
		return err
	}

	v.SetUint(i)

	return nil
}

func float32Setter(cfg *config, v reflect.Value, tag string) error {
	return floatSetter(cfg, v, tag, 32)
}

func float64Setter(cfg *config, v reflect.Value, tag string) error {
	return floatSetter(cfg, v, tag, 64)
}

func floatSetter(_ *config, v reflect.Value, tag string, bitSize int) error {
	if !v.CanFloat() {
		return fmt.Errorf("Float%dSetter does not support [%s]", bitSize, v.Kind())
	}

	f, err := strconv.ParseFloat(tag, bitSize)
	if err != nil {
		return err
	}

	v.SetFloat(f)

	return nil
}

func complex64Setter(cfg *config, v reflect.Value, tag string) error {
	return complexSetter(cfg, v, tag, 64)
}

func complex128Setter(cfg *config, v reflect.Value, tag string) error {
	return complexSetter(cfg, v, tag, 128)
}

func complexSetter(_ *config, v reflect.Value, tag string, bitSize int) error {
	if !v.CanComplex() {
		return fmt.Errorf("Complex%dSetter does not support [%s]", bitSize, v.Kind())
	}

	c, err := strconv.ParseComplex(tag, bitSize)
	if err != nil {
		return err
	}

	v.SetComplex(c)

	return nil
}

func pointerSetter(cfg *config, v reflect.Value, tag string) error {
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("PointerSetter does not support [%s]", v.Kind())
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	return valueSetter(cfg, v.Elem(), tag)
}

func structSetter(cfg *config, v reflect.Value, tag string) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("StructSetter does not support [%s]", v.Kind())
	}

	if tag != "nested" {
		return nil
	}

	return structFieldsSetter(cfg, v)
}

func arraySetter(cfg *config, v reflect.Value, tag string) error {
	if v.Kind() != reflect.Array {
		return fmt.Errorf("ArraySetter does not support [%s]", v.Kind())
	}

	if tag == "nested" {
		rv := reflect.New(v.Type().Elem()).Elem()

		if err := structFieldsSetter(cfg, rv); err != nil {
			return err
		}

		for i := 0; i < v.Len(); i++ {
			v.Index(i).Set(rv)
		}

		return nil
	}

	arr := reflect.New(reflect.ArrayOf(v.Len(), v.Type().Elem())).Elem()

	if err := json.Unmarshal([]byte(tag), arr.Addr().Interface()); err != nil {
		return err
	}

	for i := 0; i < arr.Len(); i++ {
		v.Index(i).Set(arr.Index(i))
	}

	return nil
}

func structFieldsSetter(cfg *config, v reflect.Value) error {
	v = dereference(v)

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("[%s] type is not supported. must be struct", v.Kind())
	}

	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if err := valueSetter(cfg, v.Field(i), typ.Field(i).Tag.Get(cfg.tag)); err != nil {
			return err
		}
	}

	return nil
}

func valueSetter(cfg *config, v reflect.Value, tag string) error {
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

	return fn(cfg, v, tag)
}

func dereference(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	return rv
}
