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
	"reflect"
)

func EncodeObjectHessian3V2(o *EncodeContext, obj interface{}) ([]byte, error) {
	return EncodeObjectToHessian3V2(o, nil, obj)
}

func EncodeObjectToHessian3V2(o *EncodeContext, dst []byte, obj interface{}) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeobject")
		defer o.tracer.OnTraceStop("encodeobject")
	}

	if obj == nil {
		return EncodeNilToHessian3V2(o, dst)
	}

	var (
		t     = reflect.TypeOf(obj)
		err   error
		refid int
	)

	dst, refid, err = encodeObjectrefToHessian3V2(o, dst, obj)
	if err != nil {
		return dst, err
	}
	if refid >= 0 {
		return dst, nil
	}

	v := decAllocReflectValue(reflect.ValueOf(obj))
	t = decAllocReflectType(t)
	dst, err = encodeObjectDefinitionHessian3V2(o, dst, obj, t, v)
	if err != nil {
		return dst, err
	}

	// Write object field
	for i := 0; i < v.NumField(); i++ {
		tfield := t.Field(i)
		vfield := v.Field(i)
		if !vfield.CanInterface() {
			continue
		}
		if hname := tfield.Tag.Get("hessian"); hname != "" {
			if hname == "-" {
				continue
			}
		}

		dst, err = EncodeToHessian3V2(o, dst, vfield.Interface())
		if err != nil {
			return dst, err
		}
	}

	return dst, nil
}

func encodeObjectDefinitionHessian3V2(o *EncodeContext, dst []byte,
	obj interface{}, t reflect.Type, v reflect.Value) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebjectdefinition")
		defer o.tracer.OnTraceStop("encodeobjectdefinition")
	}

	classname := getInterfaceName(obj)
	dst, ref, err := encodeObjectBeginToHessian3V2(o, dst, classname)
	if err != nil {
		return dst, err
	}

	if ref == -1 {
		count := 0
		for i := 0; i < v.NumField(); i++ {
			vfield := v.Field(i)
			if !vfield.CanInterface() {
				continue
			}
			tfield := t.Field(i)
			if hname := tfield.Tag.Get("hessian"); hname != "" {
				if hname == "-" {
					continue
				}
			}
			count++
		}

		if dst, err = EncodeInt32ToHessian3V2(o, dst, int32(count)); err != nil {
			return dst, err
		}

		for i := 0; i < v.NumField(); i++ {
			tfield := t.Field(i)
			vfield := v.Field(i)
			if !vfield.CanInterface() {
				continue
			}
			name := tfield.Name
			if hname := tfield.Tag.Get("hessian"); hname != "" {
				if hname == "-" {
					continue
				}
				name = hname
			}
			if dst, err = EncodeStringToHessian3V2(o, dst, name); err != nil {
				return dst, err
			}
		}

		// Set the ref
		dst, _, err = encodeObjectBeginToHessian3V2(o, dst, classname)
		if err != nil {
			return dst, err
		}
	}

	return dst, nil
}

func encodeObjectBeginToHessian3V2(o *EncodeContext, dst []byte, typ string) ([]byte, int, error) {
	refid, referenced, err := o.addClassrefs(typ)
	if err != nil {
		return dst, -1, err
	}

	if referenced {
		if refid >= 0 {
			dst = append(dst, 0x6f)
			dst, err = EncodeInt32ToHessian3V2(o, dst, int32(refid))
			return dst, refid, err
		}
	}

	// class definition
	dst = append(dst, 0x4F)
	dst, err = EncodeInt32ToHessian3V2(o, dst, int32(len(typ)))
	if err != nil {
		return dst, -1, err
	}
	dst = append(dst, typ...)

	return dst, -1, err
}

func EncodeClassDefinitionToHessian3V2(o *EncodeContext, dst []byte, classname string,
	fields []string) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeobjectdefinition")
		defer o.tracer.OnTraceStop("encodeobjectdefinition")
	}

	dst, ref, err := encodeObjectBeginToHessian3V2(o, dst, classname)
	if err != nil {
		return dst, err
	}

	if ref == -1 {
		if dst, err = EncodeInt32ToHessian3V2(o, dst, int32(len(fields))); err != nil {
			return dst, err
		}

		for i := range fields {
			if dst, err = EncodeStringToHessian3V2(o, dst, fields[i]); err != nil {
				return dst, err
			}
		}

		// Set the ref
		dst, _, err = encodeObjectBeginToHessian3V2(o, dst, classname)
		if err != nil {
			return dst, err
		}
	}

	return dst, nil
}
