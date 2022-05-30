// nolint
// Copyright 20xx The Alipay Authors.
//
// @authors[0]: bingwu.ybw(bingwu.ybw@antfin.com|detailyang@gmail.com)
// @authors[1]: robotx(robotx@antfin.com)
//
// *Legal Disclaimer*
// Within this source code, the comments in Chinese shall be the original, governing version. Any comment in other languages are for reference only. In the event of any conflict between the Chinese language version comments and other language version comments, the Chinese language version shall prevail.
// *法律免责声明*
// 关于代码注释部分，中文注释为官方版本，其它语言注释仅做参考。中文注释可能与其它语言注释存在不一致，当中文注释与其它语言注释存在不一致时，请以中文注释为准。
//
//

package sofahessian

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func safeIsNil(v reflect.Value) (rv bool) {
	defer func() {
		if r := recover(); r != nil {
			rv = false
		}
	}()

	rv = v.IsNil()
	return rv
}

func safeEqual(a, b interface{}) (rv bool) {
	defer func() {
		if r := recover(); r != nil {
			rv = false
		}
	}()

	if reflect.TypeOf(a).Kind() == reflect.Struct ||
		reflect.TypeOf(b).Kind() == reflect.Struct {
		return false
	}

	if a == b {
		rv = true
	} else {
		rv = false
	}
	return rv
}

func safeSetMap(m *map[interface{}]interface{}, key, value interface{}) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	// The key maybe not be hashed so we recover it.
	(*m)[key] = value
	ok = true
	return true
}

func safeSetValueByReflect(key reflect.Value, value interface{}) (err error) {
	if value == nil {
		key.Set(reflect.Zero(key.Type()))
		return nil
	}

	rv := reflect.ValueOf(value)
	if err = safeSetKeyValue(key, rv); err == nil {
		return nil
	}
	err = nil

	if key.Kind() == rv.Kind() {
		if key.Type().String() == rv.Type().String() {
			return safeSetKeyValue(key, rv)
		}
	}

	key = decAllocReflectValue(key)
	rv = decAllocReflectValue(rv)
	rt := reflect.TypeOf(rv)

	switch key.Kind() {
	case reflect.Bool:
		key.SetBool(rv.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		key.SetInt(rv.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		key.SetUint(rv.Uint())

	case reflect.Float32, reflect.Float64:
		key.SetFloat(rv.Float())

	case reflect.Complex64, reflect.Complex128:
		key.SetComplex(rv.Complex())

	case reflect.String:
		key.SetString(rv.String())

	case reflect.Ptr:
		return safeSetKeyValue(key, rv)

	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if key.Type() == rv.Type() { // Same slice
			return safeSetKeyValue(key, rv)
		}

		switch slice := value.(type) {
		case *JavaList:
			// discard the class name
			key.Set(reflect.MakeSlice(key.Type(), len(slice.value), len(slice.value)))
			for i := range slice.value {
				if err = safeSetSliceValue(key, reflect.ValueOf(slice.value[i]), i); err != nil {
					return err
				}
			}
			return nil
		case []interface{}:
			// discard the class name
			key.Set(reflect.MakeSlice(key.Type(), len(slice), len(slice)))
			for i := range slice {
				if err = safeSetSliceValue(key, reflect.ValueOf(slice[i]), i); err != nil {
					return err
				}
			}
			return nil
		case string:
			if rv.Len() == 0 {
				return
			}
			// copy string to []byte
			return safeSetKeyValue(key, reflect.ValueOf([]byte(slice)))

		default:
			return fmt.Errorf("hessian: cannot set the slice from %s to %s",
				key.Type().String(), rv.Type().String())
		}
	case reflect.Interface:
		return safeSetKeyValue(key, rv)

	case reflect.Map:
		if key.Type() == rv.Type() { // Same struct
			return safeSetKeyValue(key, rv)
		}

		// Try last
		switch m := value.(type) {
		case *JavaMap: // typed map
			// discard the class name
			key.Set(reflect.MakeMapWithSize(key.Type(), len(m.m)))
			for k, v := range m.m {
				if err = safeSetMapKeyValueByReflect(key, reflect.ValueOf(k), reflect.ValueOf(v)); err != nil {
					return err
				}
			}
			return nil

		case map[interface{}]interface{}: // untyped map
			key.Set(reflect.MakeMapWithSize(key.Type(), len(m)))
			for k, v := range m {
				if err = safeSetMapKeyValueByReflect(key, reflect.ValueOf(k), reflect.ValueOf(v)); err != nil {
					return err
				}
			}
			return nil
		default:
			return fmt.Errorf("hessian: cannot set the struct field from %s to %s",
				key.Type().String(), rv.Type().String())
		}

	case reflect.Struct:
		if key.Type() == rv.Type() { // Same struct
			return safeSetKeyValue(key, rv)
		}

		jo, ok := value.(*JavaObject)
		if !ok {
			return fmt.Errorf("hessian: cannot set the struct field from %s to %s",
				key.Type().String(), rv.Type().String())
		}

		for i := range jo.names {
			name := jo.names[i]
			value := jo.values[i]
			fkey := key.FieldByName(name)
			if !fkey.IsValid() {
				continue
			}
			if err = safeSetValueByReflect(fkey, value); err != nil {
				return err
			}
		}

	case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Uintptr:
		fallthrough
	default:
		err = errors.New("hessian: invalid kind " + rt.Kind().String())
	}

	return err
}

func safeSetSliceValue(slice, value reflect.Value, i int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("hessian: cannot set the map value from %s to %s", slice.Type(), value.Type())
		}
	}()
	slice.Index(i).Set(value)

	return err
}

func safeSetMapKeyValueByReflect(m, key, value reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("hessian: cannot set the map value from %s to %s", key.Type(), value.Type())
		}
	}()
	m.SetMapIndex(key, value)

	return err
}

func safeSetKeyValue(key, value reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("hessian: cannot set the value from %s to %s", key.Type(), value.Type())
		}
	}()
	key.Set(value)

	return err
}

func decAllocReflectType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

func decAllocIndirectReflectValue(v reflect.Value) (prev, now reflect.Value) {
	prev = v
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if v.CanSet() {
				v.Set(reflect.New(v.Type().Elem()))
			}
		}
		prev = v
		v = v.Elem()
	}
	return prev, v
}

func decAllocReflectValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if v.CanSet() {
				v.Set(reflect.New(v.Type().Elem()))
			}
		}
		v = v.Elem()
	}
	return v
}

func lookupReflectField(rt reflect.Type, name string) (int, bool) {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		val, ok := field.Tag.Lookup("hessian")
		if ok && val == name {
			return i, true
		}

		fieldname := rt.Field(i).Name
		if strings.EqualFold(fieldname, name) {
			return i, true
		}
	}
	return 0, false
}
