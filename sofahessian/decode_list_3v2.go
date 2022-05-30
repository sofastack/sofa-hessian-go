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
	"bufio"
	"fmt"
	"reflect"
)

func decodeBoundedListHessian3V2(o *DecodeContext, reader *bufio.Reader,
	list []interface{}, length int32) ([]interface{}, error) {
	for i := 0; i < int(length); i++ {
		obj, err := DecodeHessian3V2(o, reader)
		if err != nil {
			return list, err
		}
		list = append(list, obj)
	}
	return list, nil
}

func DecodeListHessian3V2(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	var i interface{}
	err := DecodeListToHessian3V2(o, reader, &i)
	return i, err
}

func DecodeListToHessian3V2(o *DecodeContext, reader *bufio.Reader, obj *interface{}) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodelist")
		defer o.tracer.OnTraceStop("decodelist")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	var (
		typ    string
		end    bool
		length int32
		c2     byte
		i32    int32
		u32    uint32
		c3     byte
	)

	if c1 == 0x56 {
		typ, err = DecodeTypeHessian3V2(o, reader)
		if err != nil {
			return err
		}
		end = true

		c2, err = reader.ReadByte()
		if err != nil {
			return err
		}

		switch c2 {
		case 'n':
			c3, err = reader.ReadByte()
			if err != nil {
				return err
			}
			length = int32(c3)

		case 'l':
			u32, err = readUint32FromReader(reader)
			if err != nil {
				return err
			}
			length = int32(u32)

		default:
			return ErrDecodeMalformedList
		}

	} else {
		i32, err = DecodeInt32Hessian3V2(o, reader)
		if err != nil {
			return err
		}

		typ, err = o.getTyperefs(int(i32))
		if err != nil {
			return err
		}
		end = false
		length, err = DecodeInt32Hessian3V2(o, reader)
		if err != nil {
			return err
		}
	}

	if length < 0 || int(length) >= o.GetMaxListLength() {
		return ErrDecodeMaxListLengthExceeded
	}

	if typ != "" {
		ci, ok := o.loadClassTypeSchema(typ)
		if ok { // concrete type
			if ci.base.Kind() != reflect.Slice && ci.base.Kind() != reflect.Array {
				return fmt.Errorf("hessian: expect slice/array type but got %s", ci.base.Kind().String())
			}

			value := reflect.MakeSlice(ci.base, int(length), int(length))
			*obj = value.Interface()
			if err = o.addObjectrefs(*obj); err != nil {
				return err
			}

			var list []interface{}

			list, err = decodeBoundedListHessian3V2(o, reader, nil, length)
			if err != nil {
				return err
			}

			if len(list) != int(length) {
				return fmt.Errorf("hessian: expect [%d]T but got [%d]T", length, len(list))
			}

			for i := range list {
				if err = safeSetValueByReflect(value.Index(i), list[i]); err != nil {
					return err
				}
			}

		} else {
			if length > 0 {
				list := make([]interface{}, 0, length)
				jl := &JavaList{class: typ, value: list}
				if err = o.addObjectrefs(jl); err != nil {
					return err
				}
				jl.value, err = decodeBoundedListHessian3V2(o, reader, jl.value, length)
				if err != nil {
					return err
				}
				*obj = jl

			} else {
				*obj = &JavaList{class: typ, value: []interface{}{}}
			}
		}

	} else {
		list := make([]interface{}, length)
		if err := o.addObjectrefs(list); err != nil {
			return err
		}
		for i := 0; i < int(length); i++ {
			obj, err := DecodeHessian3V2(o, reader)
			if err != nil {
				return err
			}
			list[i] = obj
		}
		*obj = list
	}

	if end { // read list end byte
		c1, err := reader.ReadByte()
		if err != nil {
			return err
		}
		if c1 != 'z' {
			return ErrDecodeMalformedListEnd
		}
	}

	return nil
}
