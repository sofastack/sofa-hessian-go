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

import "reflect"

// EncodeListToHessian4V2 encodes list to dst.
// list ::= x55 type value* 'Z'   # variable-length list
//      ::= 'V' type int value*   # fixed-length list
//      ::= x57 value* 'Z'        # variable-length untyped list
//      ::= x58 int value*        # fixed-length untyped list
//      ::= [x70-77] type value*  # fixed-length typed list
//      ::= [x78-7f] value*       # fixed-length untyped list
func EncodeListToHessian4V2(o *EncodeContext, dst []byte, obj interface{}) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodelist")
		defer o.tracer.OnTraceStop("encodelist")
	}

	if obj == nil {
		return EncodeNilToHessian4V2(o, dst)
	}

	value := reflect.ValueOf(obj)
	classname := getInterfaceName(obj)
	return encodeListToHessian4V2(o, dst, value, classname)
}

func encodeListToHessian4V2(o *EncodeContext, dst []byte, slice reflect.Value, typ string) ([]byte, error) {
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
			dst, refid, err = encodeObjectrefToHessian4V2(o, dst, slice.Pointer())
		} else { // Array
			dst, refid, err = encodeObjectrefToHessian4V2(o, dst, slice.Interface())
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

	dst, end, err = EncodeListBeginToHessian4V2(o, dst, length, typ)
	if err != nil {
		return dst, err
	}

	for i := 0; i < length; i++ {
		if si := slice.Index(i); si.CanInterface() {
			if dst, err = EncodeToHessian4V2(o, dst, si.Interface()); err != nil {
				return dst, err
			}
		} else {
			return dst, ErrEncodeSliceElemCannotBeInterfaced
		}
	}

	if end {
		dst = append(dst, 'z')
	}

	return dst, nil
}

func EncodeListWithLengthToHessian4V2(o *EncodeContext, dst []byte, length int,
	fn func(i int, o *EncodeContext, dst []byte) ([]byte, error)) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodelist")
		defer o.tracer.OnTraceStop("encodelist")
	}

	var (
		end bool
		err error
	)

	dst, end, err = EncodeListBeginToHessian4V2(o, dst, length, "")
	if err != nil {
		return dst, err
	}

	for i := 0; i < length; i++ {
		dst, err = fn(i, o, dst)
		if err != nil {
			return dst, err
		}
	}

	return EncodeListEndToHessian4V2(o, dst, end)
}

func EncodeListEndToHessian4V2(o *EncodeContext, dst []byte, end bool) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodelistend")
		defer o.tracer.OnTraceStop("encodelistend")
	}
	if end {
		dst = append(dst, 'z')
	}
	return dst, nil
}

// EncodeListBeginToHessian4V2 encodes slice to dst.
//
// list ::= x55 type value* 'Z'   # variable-length list
//      ::= 'V' type int value*   # fixed-length list
//      ::= x57 value* 'Z'        # variable-length untyped list
//      ::= x58 int value*        # fixed-length untyped list
//      ::= [x70-77] type value*  # fixed-length typed list
//      ::= [x78-7f] value*       # fixed-length untyped list
func EncodeListBeginToHessian4V2(o *EncodeContext, dst []byte, length int, typ string) ([]byte, bool, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodelistbegin")
		defer o.tracer.OnTraceStop("encodelistbegin")
	}

	var err error
	if length < 0 {
		if typ == "" {
			dst = append(dst, 'V')
			dst, err = encodeTyperefToHessian4V2(o, dst, typ)
			return dst, true, err
		}
		dst = append(dst, 'W')
		return dst, true, err

	} else if length <= LISTDIRECTMAX {
		if typ != "" {
			dst = append(dst, uint8(LISTDIRECT+length))
			dst, err = encodeTyperefToHessian4V2(o, dst, typ)
			return dst, false, err
		}

		dst = append(dst, uint8(LISTDIRECTUNTYPED+length))
		return dst, false, err
	}

	if typ != "" {
		dst = append(dst, LISTFIXED)
		dst, err = encodeTyperefToHessian4V2(o, dst, typ)
		if err != nil {
			return dst, false, err
		}

	} else {
		dst = append(dst, LISTFIXEDUNTYPED)
	}

	dst, err = EncodeInt32ToHessian4V2(o, dst, int32(length))
	return dst, false, err
}
