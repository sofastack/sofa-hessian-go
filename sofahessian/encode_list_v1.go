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
)

func EncodeListToHessianV1(o *EncodeContext, dst []byte, obj interface{}) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodelist")
		defer o.tracer.OnTraceStop("encodelist")
	}

	if obj == nil {
		return EncodeNilToHessianV1(o, dst)
	}

	value := reflect.ValueOf(obj)
	classname := getInterfaceName(obj)
	return encodeListToHessianV1(o, dst, value, classname)
}

func EncodeListBeginToHessianV1(o *EncodeContext, dst []byte, length int, typ string) ([]byte, bool, error) {
	refid, _, err := o.getTyperefs(typ)
	if err != nil {
		return dst, false, err
	}

	if refid >= 0 {
		dst = append(dst, 'R')
		dst, err = EncodeInt32RefToHessianV1(o, dst, int32(refid))
		if err != nil {
			return dst, false, err
		}
		return dst, false, err
	}

	dst = append(dst, 'V')
	dst, err = encodeTyperefToHessian3V2(o, dst, typ)
	if err != nil {
		return dst, false, err
	}

	dst = append(dst, "l0000"...)
	binary.BigEndian.PutUint32(dst[len(dst)-4:], uint32(length))

	return dst, true, nil
}

func EncodeListEndToHessianV1(o *EncodeContext, dst []byte, end bool) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodelistend")
		defer o.tracer.OnTraceStop("encodelistend")
	}

	if end {
		dst = append(dst, 'z')
	}
	return dst, nil
}

func encodeListToHessianV1(o *EncodeContext, dst []byte, slice reflect.Value, typ string) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodelistbegin")
		defer o.tracer.OnTraceStop("encodelistbegin")
	}

	// Unwrap the pointer if we can
	slice = reflect.Indirect(slice)

	if slice.Kind() != reflect.Slice &&
		slice.Kind() != reflect.Array {
		return dst, ErrEncodeNotSliceType
	}

	var (
		err   error
		refid int
	)

	if !o.disableObjectrefs {
		if slice.Kind() == reflect.Slice {
			// []interface{} cannot be hashed so we use address instead.
			dst, refid, err = encodeObjectrefToHessianV1(o, dst, slice.Pointer())
		} else { // Array
			dst, refid, err = encodeObjectrefToHessianV1(o, dst, slice.Interface())
		}

		if err != nil {
			return dst, err
		}

		if refid >= 0 {
			return dst, nil
		}
	}

	var (
		end    bool
		length = slice.Len()
	)

	dst, end, err = EncodeListBeginToHessianV1(o, dst, length, typ)
	if err != nil {
		return dst, err
	}

	for i := 0; i < length; i++ {
		if slice.Index(i).CanInterface() {
			if dst, err = EncodeToHessianV1(o, dst, slice.Index(i).Interface()); err != nil {
				return dst, err
			}
		} else {
			return dst, ErrEncodeSliceElemCannotBeInterfaced
		}
	}

	return EncodeListEndToHessianV1(o, dst, end)
}
