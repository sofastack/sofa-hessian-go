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
	"reflect"
)

// JavaClassNameGetter represents the java class name interface.
type JavaClassNameGetter interface {
	GetJavaClassName() string
}

type HessianEncoder interface {
	HessianEncode(ctx *EncodeContext, b []byte) ([]byte, error)
}

type HessianDecoder interface {
	HessianDecode(ctx *DecodeContext, br *bufio.Reader) error
}

var (
	JavaClassNameGetterInterfaceType = reflect.TypeOf((*JavaClassNameGetter)(nil)).Elem()
	HessianEncoderInterfaceType      = reflect.TypeOf((*HessianEncoder)(nil)).Elem()
	HessianDecoderInterfaceType      = reflect.TypeOf((*HessianDecoder)(nil)).Elem()
)

type JavaList struct {
	class string
	value []interface{}
}

func (jl *JavaList) HessianEncode(ctx *EncodeContext, b []byte) ([]byte, error) {
	if ctx.version == Hessian4xV2 {
		return jl.HessianEncode4V2(ctx, b)
	} else if ctx.version == Hessian3xV2 {
		return jl.HessianEncode3V2(ctx, b)
	} else {
		return jl.HessianEncodeV1(ctx, b)
	}
}

func (jl *JavaList) HessianEncodeV1(o *EncodeContext, dst []byte) ([]byte, error) {
	if jl == nil {
		return EncodeNilToHessian4V2(o, dst)
	}

	value := reflect.ValueOf(jl.value)
	classname := jl.class
	return encodeListToHessian4V2(o, dst, value, classname)
}

func (jl *JavaList) HessianEncode4V2(o *EncodeContext, dst []byte) ([]byte, error) {
	if jl == nil {
		return EncodeNilToHessian4V2(o, dst)
	}

	classname := jl.class
	value := reflect.ValueOf(jl.value)
	return encodeListToHessian4V2(o, dst, value, classname)
}

func (jl *JavaList) HessianEncode3V2(o *EncodeContext, dst []byte) ([]byte, error) {
	if jl == nil {
		return EncodeNilToHessian3V2(o, dst)
	}

	value := reflect.ValueOf(jl.value)
	classname := jl.class
	return encodeListToHessian3V2(o, dst, value, classname)
}

func (jl *JavaList) GetJavaClassName() string {
	return jl.class
}

func (jl *JavaList) GetValue() []interface{} {
	return jl.value
}

type JavaMap struct {
	class string
	m     map[interface{}]interface{}
}

func NewJavaMap(class string, m map[interface{}]interface{}) *JavaMap {
	return &JavaMap{class: class, m: m}
}

func (j *JavaMap) GetJavaClassName() string { return j.class }

func (j *JavaMap) GetValue() map[interface{}]interface{} { return j.m }

func (j *JavaMap) HessianEncode(ctx *EncodeContext, b []byte) ([]byte, error) {
	if ctx.version == Hessian4xV2 {
		return j.HessianEncode4V2(ctx, b)
	} else if ctx.version == Hessian3xV2 {
		return j.HessianEncode3V2(ctx, b)
	} else {
		return j.HessianEncodeV1(ctx, b)
	}
}

func (j *JavaMap) HessianEncodeV1(ctx *EncodeContext, b []byte) ([]byte, error) {
	var err error
	b, err = EncodeMapBeginToHessianV1(ctx, b, j.class)
	if err != nil {
		return b, err
	}

	for k := range j.m {
		b, err = EncodeToHessianV1(ctx, b, k)
		if err != nil {
			return b, err
		}

		b, err = EncodeToHessianV1(ctx, b, j.m[k])
		if err != nil {
			return b, err
		}
	}

	if err != nil {
		return b, err
	}

	return EncodeMapEndToHessianV1(ctx, b)
}

func (j *JavaMap) HessianEncode3V2(ctx *EncodeContext, b []byte) ([]byte, error) {
	var err error
	b, err = EncodeMapBeginToHessian3V2(ctx, b, j.class)
	if err != nil {
		return b, err
	}

	for k := range j.m {
		b, err = EncodeToHessian3V2(ctx, b, k)
		if err != nil {
			return b, err
		}

		b, err = EncodeToHessian3V2(ctx, b, j.m[k])
		if err != nil {
			return b, err
		}
	}

	if err != nil {
		return b, err
	}

	return EncodeMapEndToHessian3V2(ctx, b)
}

func (j *JavaMap) HessianEncode4V2(ctx *EncodeContext, b []byte) ([]byte, error) {
	var err error
	b, err = EncodeMapBeginToHessian4V2(ctx, b, j.class)
	if err != nil {
		return b, err
	}

	for k := range j.m {
		b, err = EncodeToHessian4V2(ctx, b, k)
		if err != nil {
			return b, err
		}

		b, err = EncodeToHessian4V2(ctx, b, j.m[k])
		if err != nil {
			return b, err
		}
	}

	if err != nil {
		return b, err
	}

	return EncodeMapEndToHessian4V2(ctx, b)
}

type JavaObject struct {
	class  string
	names  []string
	values []interface{}
}

func (jo *JavaObject) HessianEncode(ctx *EncodeContext, b []byte) ([]byte, error) {
	if ctx.version == Hessian4xV2 {
		return jo.HessianEncode4V2(ctx, b)
	} else if ctx.version == Hessian3xV2 {
		return jo.HessianEncode3V2(ctx, b)
	} else if ctx.version == HessianV1 {
		// todo not working yet
	}

	return b, nil
}

func (jo *JavaObject) HessianEncode4V2(o *EncodeContext, dst []byte) ([]byte, error) {
	if jo == nil {
		return EncodeNilToHessian4V2(o, dst)
	}

	var (
		refid int
		err   error
	)
	dst, refid, err = encodeObjectrefToHessian4V2(o, dst, jo)
	if err != nil {
		return dst, err
	}
	if refid >= 0 {
		return dst, nil
	}

	v := reflect.ValueOf(*jo)
	t := reflect.TypeOf(*jo)
	dst, err = jo.EncodeObjectDefinitionHessian4V2(o, dst, jo, t, v)
	if err != nil {
		return dst, err
	}

	// Write object field
	for i := range jo.values {
		dst, err = EncodeToHessian4V2(o, dst, jo.values[i])
		if err != nil {
			return dst, err
		}
	}

	return dst, nil
}

func (jo *JavaObject) HessianEncode3V2(o *EncodeContext, dst []byte) ([]byte, error) {
	if jo == nil {
		return EncodeNilToHessian3V2(o, dst)
	}

	var (
		refid int
		err   error
	)
	dst, refid, err = encodeObjectrefToHessian3V2(o, dst, jo)
	if err != nil {
		return dst, err
	}
	if refid >= 0 {
		return dst, nil
	}

	dst, err = jo.EncodeObjectDefinitionHessian3V2(o, dst)
	if err != nil {
		return dst, err
	}

	// Write object field
	for i := range jo.values {
		dst, err = EncodeToHessian3V2(o, dst, jo.values[i])
		if err != nil {
			return dst, err
		}
	}

	return dst, nil
}

func (jo *JavaObject) EncodeObjectDefinitionHessian3V2(o *EncodeContext, dst []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebjectdefinition")
		defer o.tracer.OnTraceStop("encodeobjectdefinition")
	}

	classname := jo.class
	dst, ref, err := encodeObjectBeginToHessian3V2(o, dst, classname)
	if err != nil {
		return dst, err
	}

	if ref == -1 {
		count := len(jo.names)

		if dst, err = EncodeInt32ToHessian3V2(o, dst, int32(count)); err != nil {
			return dst, err
		}

		for i := range jo.names {
			if dst, err = EncodeStringToHessian3V2(o, dst, jo.names[i]); err != nil {
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

func (jo *JavaObject) EncodeObjectDefinitionHessian4V2(o *EncodeContext, dst []byte,
	obj interface{}, t reflect.Type, v reflect.Value) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebjectdefinition")
		defer o.tracer.OnTraceStop("encodeobjectdefinition")
	}

	classname := jo.class
	dst, ref, err := encodeObjectBeginToHessian4V2(o, dst, classname)
	if err != nil {
		return dst, err
	}

	if ref == -1 {
		count := len(jo.names)

		if dst, err = EncodeInt32ToHessian4V2(o, dst, int32(count)); err != nil {
			return dst, err
		}

		for i := range jo.names {
			if dst, err = EncodeStringToHessian4V2(o, dst, jo.names[i]); err != nil {
				return dst, err
			}
		}

		// Set the ref
		dst, _, err = encodeObjectBeginToHessian4V2(o, dst, classname)
		if err != nil {
			return dst, err
		}
	}

	return dst, nil
}

func (jo *JavaObject) Len() int {
	return len(jo.names)
}

func (jo *JavaObject) GetKey(i int) string {
	return jo.names[i]
}

func (jo *JavaObject) GetValue(i int) interface{} {
	return jo.values[i]
}

func (jo *JavaObject) GetJavaClassName() string {
	return jo.class
}
