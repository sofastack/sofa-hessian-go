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
	"encoding/binary"
	"reflect"
	"sort"
)

// EncodeMapToHessian3V2 encodes map to dst.
func EncodeMapToHessian3V2(o *EncodeContext, dst []byte, obj interface{}) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodemap")
		defer o.tracer.OnTraceStop("encodemap")
	}

	if obj == nil {
		return EncodeNilToHessian3V2(o, dst)
	}

	// Allow *map to reduce recursive encodeto call
	t := reflect.TypeOf(obj)
	if kind := t.Kind(); kind != reflect.Map {
		if kind == reflect.Ptr {
			if t.Elem().Kind() != reflect.Map {
				return dst, ErrEncodeNotMapType
			}
		} else {
			return dst, ErrEncodeNotMapType
		}
	}

	v := reflect.ValueOf(obj)

	var (
		err   error
		refid int
	)

	if !o.disableObjectrefs {
		// Map cannot be hashed, use pointer instead.
		dst, refid, err = encodeObjectrefToHessian3V2(o, dst, v.Pointer())
		if err != nil {
			return dst, err
		}

		if refid >= 0 {
			return dst, nil
		}
	}

	classname := getInterfaceName(obj)
	dst, err = EncodeMapBeginToHessian3V2(o, dst, classname)
	if err != nil {
		return dst, err
	}

	// Unwrap the pointer if can
	v = reflect.Indirect(v)

	// Map in golang is unordered but other languages maybe or not maybe unordered.
	keys := v.MapKeys()
	if o.less == nil { // Fast path
		for i := range keys {
			key := keys[i]
			if key.CanInterface() { // Fast path
				if dst, err = EncodeToHessian3V2(o, dst, key.Interface()); err != nil {
					return dst, err
				}
			} else {
				if dst, err = EncodeValueToHessian3v2(o, dst, key); err != nil {
					return dst, err
				}
			}

			value := v.MapIndex(key)
			if value.CanInterface() { // Fast path
				if dst, err = EncodeToHessian3V2(o, dst, value.Interface()); err != nil {
					return dst, err
				}
			} else {
				if dst, err = EncodeValueToHessian3v2(o, dst, value); err != nil {
					return dst, err
				}
			}
		}
	} else {
		keys := keys
		sorted := make([]reflect.Value, 0, len(keys))
		for i := range keys {
			sorted = append(sorted, keys[i])
		}

		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].CanInterface() && sorted[j].CanInterface() {
				ii := sorted[i]
				keyi := ii.Interface()
				valuei := v.MapIndex(ii)
				ji := sorted[j]
				keyj := ji.Interface()
				valuej := v.MapIndex(ji)
				return o.less(keyi, keyj, valuei, valuej)
			}
			return false
		})
		for i := 0; i < len(sorted); i++ {
			key := sorted[i]
			if key.CanInterface() { // Fast path
				if dst, err = EncodeToHessian3V2(o, dst, key.Interface()); err != nil {
					return dst, err
				}
			} else {
				if dst, err = EncodeValueToHessian3v2(o, dst, key); err != nil {
					return dst, err
				}
			}

			value := v.MapIndex(key)
			if value.CanInterface() { // Fast path
				if dst, err = EncodeToHessian3V2(o, dst, value.Interface()); err != nil {
					return dst, err
				}
			} else {
				if dst, err = EncodeValueToHessian3v2(o, dst, value); err != nil {
					return dst, err
				}
			}
		}
	}

	return EncodeMapEndToHessian3V2(o, dst)
}

func EncodeMapBeginToHessian3V2(o *EncodeContext, dst []byte, typ string) ([]byte, error) {
	dst = append(dst, 'M')
	if typ != "" {
		dst = append(dst, 0x74, 0, 0)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(len(typ)))
		dst = append(dst, typ...)
	}

	return dst, nil
}

func EncodeMapEndToHessian3V2(o *EncodeContext, dst []byte) ([]byte, error) {
	dst = append(dst, 0x7a)
	return dst, nil
}
