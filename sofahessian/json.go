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
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/valyala/fastjson"
)

func init() {
	RegisterJSONToHessianEncoder("int", JSONToHessianEncoderFunc(EncodeJSONInt32))
	RegisterJSONToHessianEncoder("java.lang.Integer", JSONToHessianEncoderFunc(EncodeJSONInt32))
	RegisterJSONToHessianEncoder("long", JSONToHessianEncoderFunc(EncodeJSONInt64))
	RegisterJSONToHessianEncoder("java.lang.Long", JSONToHessianEncoderFunc(EncodeJSONInt64))
	RegisterJSONToHessianEncoder("double", JSONToHessianEncoderFunc(EncodeJSONFloat64))
	RegisterJSONToHessianEncoder("java.lang.Double", JSONToHessianEncoderFunc(EncodeJSONFloat64))
	RegisterJSONToHessianEncoder("date", JSONToHessianEncoderFunc(EncodeJSONDate))
	RegisterJSONToHessianEncoder("java.util.Date", JSONToHessianEncoderFunc(EncodeJSONDate))
	RegisterJSONToHessianEncoder("bytes", JSONToHessianEncoderFunc(EncodeJSONBytes))
}

type ToJSONEncoder interface {
	EncodeToJSON(ctx *JSONEncodeContext, dst []byte, obj interface{}) ([]byte, error)
}

type ToJSONEncoderFunc func(ctx *JSONEncodeContext, dst []byte, obj interface{}) ([]byte, error)

func (t ToJSONEncoderFunc) EncodeToJSON(ctx *JSONEncodeContext, dst []byte, obj interface{}) ([]byte, error) {
	return t(ctx, dst, obj)
}

// JSONToHessianEncoder represents the custom encoder to encode json to hessian.
type JSONToHessianEncoder interface {
	EncodeJSONToHessian(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value, classname string) error
}

type JSONToHessianEncoderFunc func(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value, classname string) error

func (j JSONToHessianEncoderFunc) EncodeJSONToHessian(ctx *JSONEncodeContext,
	encoder *Encoder, v *fastjson.Value, classname string) error {
	return j(ctx, encoder, v, classname)
}

var (
	globalJSONToHessianEncoderMap sync.Map // map[string]JSONToHessianEncoder
	globalToJSONEncoderMap        sync.Map // map[string]ToJSONEncoder
	globalJSONParserPool          fastjson.ParserPool
)

// RegisterJSONToHessianEncoder registers the JSONToHessianEncoder to global.
func RegisterJSONToHessianEncoder(classname string, e JSONToHessianEncoder) {
	globalJSONToHessianEncoderMap.Store(classname, e)
}

// RegisterToJSONEncoder registers the ToJSONEncoder to global.
func RegisterToJSONEncoder(classname string, e ToJSONEncoder) {
	globalToJSONEncoderMap.Store(classname, e)
}

// LoadToJSONEncoder loads the ToJSONEncoder from global.
func LoadToJSONEncoder(classname string) (ToJSONEncoder, bool) {
	i, ok := globalToJSONEncoderMap.Load(classname)
	if !ok {
		return nil, false
	}
	return i.(ToJSONEncoder), true
}

// LoadJSONToHessianEncoder loads the JSONToHessianEncoder from global.
func LoadJSONToHessianEncoder(classname string) (JSONToHessianEncoder, bool) {
	i, ok := globalJSONToHessianEncoderMap.Load(classname)
	if !ok {
		return nil, false
	}
	return i.(JSONToHessianEncoder), true
}

// JSONEncodeContext holds the context of json encoding.
type JSONEncodeContext struct {
	ptrlevel uint
	// ptrtracker tracks the pointer to avoid stack overflow.
	ptrtracker   map[interface{}]struct{}
	maxPtrCycles uint
	depth        int
	maxDepth     int
	tracer       Tracer
	less         func(keyi, keyj, valuei, valuej interface{}) bool
}

// NewJSONEncodeContext returns a new JSONEncodeContext.
func NewJSONEncodeContext() *JSONEncodeContext {
	return &JSONEncodeContext{
		ptrtracker: make(map[interface{}]struct{}),
	}
}

// Reset resets the context.
func (e *JSONEncodeContext) Reset() {
	e.ptrlevel = 0
	for k := range e.ptrtracker {
		delete(e.ptrtracker, k)
	}
	e.maxPtrCycles = 0
	e.depth = 0
	e.maxDepth = 0
	e.tracer = nil
	e.less = nil
}

func (e *JSONEncodeContext) LoadToJSONEncoder(className string) (ToJSONEncoder, bool) {
	return LoadToJSONEncoder(className)
}

func (e *JSONEncodeContext) SetMaxPtrCycles(m uint) *JSONEncodeContext {
	e.maxPtrCycles = m
	return e
}

func (e *JSONEncodeContext) SetLessFunc(fn func(keyi, keyj, valuei, valuej interface{}) bool) *JSONEncodeContext {
	e.less = fn
	return e
}

// SetTracer sets the tracer for JSONEncodeContext.
func (e *JSONEncodeContext) SetTracer(tracer Tracer) *JSONEncodeContext {
	e.tracer = tracer
	return e
}

// SetMaxDepth sets the maximum depth for recursive.
func (e *JSONEncodeContext) SetMaxDepth(depth int) *JSONEncodeContext {
	e.maxDepth = depth
	return e
}

// EncodeRawJSONBytesToHessian transforms the json bytes to hessian encoder.
func EncodeRawJSONBytesToHessian(ctx *JSONEncodeContext, encoder *Encoder, json []byte) error {
	p := globalJSONParserPool.Get()
	v, err := p.ParseBytes(json)
	if err != nil {
		globalJSONParserPool.Put(p)
		return ErrDecodeMalformedJSON
	}

	err = EncodeJSONToHessian(ctx, encoder, v)

	globalJSONParserPool.Put(p)

	return err
}

// EncodeRawJSONStringToHessian transforms the json string to hessian encoder.
func EncodeRawJSONStringToHessian(ctx *JSONEncodeContext, encoder *Encoder, json string) error {
	p := globalJSONParserPool.Get()
	v, err := p.Parse(json)
	if err != nil {
		globalJSONParserPool.Put(p)
		return ErrDecodeMalformedJSON
	}

	err = EncodeJSONToHessian(ctx, encoder, v)

	globalJSONParserPool.Put(p)

	return err
}

// EncodeJSONToHessian transforms fastjson to hessian encoder.
func EncodeJSONToHessian(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value) (err error) {
	ctx.depth++

	if ctx.maxDepth != 0 && ctx.depth > ctx.maxDepth {
		return ErrEncodeTooManyNestedJSON
	}

	switch v.Type() {
	case fastjson.TypeObject:
		err = EncodeJSONObjectToHessian(ctx, encoder, v)
		break

	case fastjson.TypeArray:
		err = EncodeJSONArrayToHessian(ctx, encoder, v, "")
		break

	case fastjson.TypeNull:
		err = EncodeJSONNULLToHessian(ctx, encoder, v)
		break

	case fastjson.TypeString:
		err = EncodeJSONStringToHessian(ctx, encoder, v)
		break

	case fastjson.TypeNumber:
		err = EncodeJSONNumberToHessian(ctx, encoder, v)
		break

	case fastjson.TypeTrue:
		err = EncodeJSONBoolToHessian(ctx, encoder, true)
		break

	case fastjson.TypeFalse:
		err = EncodeJSONBoolToHessian(ctx, encoder, false)
		break
	}

	ctx.depth--

	return err
}

// EncodeJSONBytes encodes JSON bytes(base64 encoding) to hessian encoder.
func EncodeJSONBytes(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value, classname string) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json base64")
		defer ctx.tracer.OnTraceStop("json base64")
	}

	i := v.GetStringBytes()
	n := base64.StdEncoding.DecodedLen(len(i))
	dst := make([]byte, n)
	rn, err := base64.StdEncoding.Decode(dst, i)
	if err != nil {
		return err
	}

	return encoder.EncodeBinary(dst[:rn])
}

// EncodeJSONDate encodes JSON boolean to hessian encoder.
func EncodeJSONDate(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value, classname string) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json java.lang.date")
		defer ctx.tracer.OnTraceStop("json java.lang.date")
	}

	realvalue := v.Get("$")
	if realvalue == nil {
		return fmt.Errorf("expect get the date from the key \"$\" but got nil")
	}

	i, err := realvalue.Int()
	if err != nil {
		return fmt.Errorf("expect get the date from the key \"$\": %s", err)
	}

	return encoder.EncodeDate(time.Unix(int64(i/1000), int64(i)%1000*10e5))
}

// EncodeJSONBoolToHessian encodes JSON boolean to hessian encoder.
func EncodeJSONBoolToHessian(ctx *JSONEncodeContext, encoder *Encoder, b bool) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json boolean")
		defer ctx.tracer.OnTraceStop("json boolean")
	}
	return encoder.EncodeBool(b)
}

// EncodeJSONInt32 encodes JSON int32 to hessian encoder.
func EncodeJSONInt32(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value, classname string) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json int32")
		defer ctx.tracer.OnTraceStop("json int32")
	}

	i, err := v.Int()
	if err != nil {
		return errors.New("expect get the int from the key \"$\"")
	}

	return encoder.EncodeInt32(int32(i))
}

// EncodeJSONFloat64 encodes JSON float64 to hessian encoder.
func EncodeJSONFloat64(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value, classname string) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json java.lang.double")
		defer ctx.tracer.OnTraceStop("json java.lang.double")
	}

	i, err := v.Float64()
	if err != nil {
		return errors.New("expect get the double from the key \"$\"")
	}
	return encoder.EncodeFloat64(i)
}

// EncodeJSONInt64 encodes JSON int64 to hessian encoder.
func EncodeJSONInt64(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value, classname string) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json int64")
		defer ctx.tracer.OnTraceStop("json int64")
	}

	i, err := v.Int()
	if err != nil {
		return errors.New("expect get the int from the key \"$\"")
	}

	return encoder.EncodeInt64(int64(i))
}

// EncodeJSONNumberToHessian encodes JSON number to hessian encoder.
func EncodeJSONNumberToHessian(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json number")
		defer ctx.tracer.OnTraceStop("json number")
	}

	f64 := v.GetFloat64()
	if f64 == float64(uint64(f64)) { // int
		return encoder.EncodeInt64(int64(f64))
	}

	return encoder.EncodeFloat64(f64)
}

// EncodeJSONStringToHessian encodes JSON string to hessian encoder.
func EncodeJSONStringToHessian(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json number")
		defer ctx.tracer.OnTraceStop("json number")
	}

	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json string")
		defer ctx.tracer.OnTraceStop("json string")
	}
	return encoder.EncodeString(b2s(v.GetStringBytes()))
}

// EncodeJSONNULLToHessian encodes JSON NULL to hessian encoder.
func EncodeJSONNULLToHessian(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json null")
		defer ctx.tracer.OnTraceStop("json null")
	}
	return encoder.EncodeNil()
}

// EncodeJSONArrayToHessian encodes JSON Array to hessian encoder.
func EncodeJSONArrayToHessian(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value, typ string) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json array")
		defer ctx.tracer.OnTraceStop("json array")
	}

	values, err := v.Array()
	if err != nil {
		return err
	}

	var end bool
	end, err = encoder.EncodeListBegin(len(values), typ)
	if err != nil {
		return err
	}

	for i := range values {
		if err = EncodeJSONToHessian(ctx, encoder, values[i]); err != nil {
			break
		}
	}

	if err = encoder.EncodeListEnd(end); err != nil {
		return err
	}

	return err
}

// EncodeJSONObjectToHessian encodes JSON Object to hessian encoder.
func EncodeJSONObjectToHessian(ctx *JSONEncodeContext, encoder *Encoder, v *fastjson.Value) error {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json object")
		defer ctx.tracer.OnTraceStop("json object")
	}

	object, err := v.Object()
	if err != nil {
		return err
	}

	// check the class
	cls := object.Get("$class")
	if cls == nil { // map
		if err = encoder.EncodeMapBegin(); err != nil {
			return err
		}

		object.Visit(func(key []byte, v *fastjson.Value) {
			if err != nil {
				return
			}

			// key
			if err = encoder.EncodeString(string(key)); err != nil {
				return
			}

			// value
			if err = EncodeJSONToHessian(ctx, encoder, v); err != nil {
				return
			}
		})
		if err != nil {
			return err
		}

		if err = encoder.EncodeMapEnd(); err != nil {
			return err
		}

		return nil
	}

	// class
	if cls.Type() != fastjson.TypeString {
		return errors.New("expect get the string from the \"$class\"")
	}

	classname := string(cls.GetStringBytes())
	if classname == "" {
		return errors.New("$class cannot be empty")
	}

	realvalue := object.Get("$")
	if realvalue == nil {
		return errors.New("expect get the real value from the key \"$\"")
	}

	je, ok := LoadJSONToHessianEncoder(classname)
	if ok {
		return je.EncodeJSONToHessian(ctx, encoder, v, classname)
	}

	if realvalue.Type() == fastjson.TypeArray { // array
		return EncodeJSONArrayToHessian(ctx, encoder, realvalue, classname)
	}

	var realobject *fastjson.Object
	realobject, err = realvalue.Object()
	if err != nil {
		return err
	}

	fields := make([]string, 0, 8)
	realobject.Visit(func(key []byte, v *fastjson.Value) {
		fields = append(fields, string(key))
	})

	// encode class definition
	if err = encoder.EncodeClassDefinition(classname, fields); err != nil {
		return err
	}

	// encode class field
	for i := range fields {
		field := fields[i]
		fieldvalue := realobject.Get(field)
		if err = EncodeJSONToHessian(ctx, encoder, fieldvalue); err != nil {
			return err
		}
	}

	return err
}

// EncodeToJSON encodes the interface to dst.
func EncodeToJSON(ctx *JSONEncodeContext, dst []byte, obj interface{}) ([]byte, error) {
	ctx.depth++

	if ctx.maxDepth != 0 && ctx.depth > ctx.maxDepth {
		return dst, ErrEncodeTooManyNestedJSON
	}

	if obj == nil {
		return EncodeNilToJSON(ctx, dst)
	}

	var err error
	switch t := obj.(type) { // fast path without reflection
	case *[]byte:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		return EncodeBinaryToJSON(ctx, dst, *t)
	case []byte:
		return EncodeBinaryToJSON(ctx, dst, t)
	case *uint8:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInt32ToJSON(ctx, dst, int32(*t))
	case uint8:
		dst, err = EncodeInt32ToJSON(ctx, dst, int32(t))
	case *int8:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInt32ToJSON(ctx, dst, int32(*t))
	case int8:
		dst, err = EncodeInt32ToJSON(ctx, dst, int32(t))
	case *uint16:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInt32ToJSON(ctx, dst, int32(*t))
	case uint16:
		dst, err = EncodeInt32ToJSON(ctx, dst, int32(t))
	case *int16:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
	case int16:
		dst, err = EncodeInt32ToJSON(ctx, dst, int32(t))
	case *uint32:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInt32ToJSON(ctx, dst, int32(*t))
	case uint32:
		dst, err = EncodeInt32ToJSON(ctx, dst, int32(t))
	case *int32:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInt32ToJSON(ctx, dst, *t)
	case int32:
		dst, err = EncodeInt32ToJSON(ctx, dst, t)
	case *uint64:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInt64ToJSON(ctx, dst, int64(*t))
	case uint64:
		dst, err = EncodeInt64ToJSON(ctx, dst, int64(t))
	case *int64:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInt64ToJSON(ctx, dst, *t)
	case int64:
		dst, err = EncodeInt64ToJSON(ctx, dst, t)
	case *uint:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInt64ToJSON(ctx, dst, int64(*t))
	case uint:
		dst, err = EncodeInt64ToJSON(ctx, dst, int64(t))
	case *int:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInt64ToJSON(ctx, dst, int64(*t))
	case int:
		dst, err = EncodeInt64ToJSON(ctx, dst, int64(t))
	case *float32:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeFloat32ToJSON(ctx, dst, *t)
	case float32:
		dst, err = EncodeFloat32ToJSON(ctx, dst, t)
	case *float64:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeFloat64ToJSON(ctx, dst, *t)
	case float64:
		dst, err = EncodeFloat64ToJSON(ctx, dst, t)
	case *bool:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeBoolToJSON(ctx, dst, *t)
	case bool:
		dst, err = EncodeBoolToJSON(ctx, dst, t)
	case *string:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeStringToJSON(ctx, dst, *t)
	case string:
		dst, err = EncodeStringToJSON(ctx, dst, t)
	case *map[string]interface{}:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeStringInterfaceMapToJSON(ctx, dst, *t)
	case map[string]interface{}:
		dst, err = EncodeStringInterfaceMapToJSON(ctx, dst, t)
	case *map[interface{}]interface{}:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInterfaceInterfaceMapToJSON(ctx, dst, *t)
	case map[interface{}]interface{}:
		dst, err = EncodeInterfaceInterfaceMapToJSON(ctx, dst, t)
	case *[]interface{}:
		if t == nil {
			dst, err = EncodeNilToJSON(ctx, dst)
			break
		}
		dst, err = EncodeInterfaceSliceToJSON(ctx, dst, *t)
	case []interface{}:
		dst, err = EncodeInterfaceSliceToJSON(ctx, dst, t)

	case *JavaObject:
		if ctx.maxPtrCycles != 0 {
			if ctx.ptrlevel++; ctx.ptrlevel >= ctx.maxPtrCycles {
				ptr := t
				if _, ok := ctx.ptrtracker[ptr]; ok {
					// cannot encode process CyclicReference and use {
					//    "$class": "CyclicReference",
					//    "$": "type name"
					// } replace
					return EncodeStringInterfaceMapToJSON(ctx, dst, map[string]interface{}{
						"$class": "CyclicReference",
						"$":      "JavaObject",
					})
				}
				ctx.ptrtracker[ptr] = struct{}{}
				defer delete(ctx.ptrtracker, ptr)
			}
			// maybe slow down
			defer func() {
				ctx.ptrlevel--
			}()
		}

		if je, ok := ctx.LoadToJSONEncoder(t.class); ok {
			dst, err = je.EncodeToJSON(ctx, dst, obj)
		} else {
			m := acquireMStringInterface()
			for i := 0; i < t.Len(); i++ {
				name := t.names[i]
				value := t.values[i]
				(*m)[name] = value
			}
			container := map[string]interface{}{
				"$class": t.class,
				"$":      *m,
			}
			dst, err = EncodeStringInterfaceMapToJSON(ctx, dst, container)
			releaseMStringInterface(m)
		}

	case *JavaMap:
		if ctx.maxPtrCycles != 0 {
			if ctx.ptrlevel++; ctx.ptrlevel >= ctx.maxPtrCycles {
				ptr := t
				if _, ok := ctx.ptrtracker[ptr]; ok {
					// cannot encode process CyclicReference and use {
					//    "$class": "CyclicReference",
					//    "$": "type name"
					// } replace
					return EncodeStringInterfaceMapToJSON(ctx, dst, map[string]interface{}{
						"$class": "CyclicReference",
						"$":      "JavaObject",
					})
				}
				ctx.ptrtracker[ptr] = struct{}{}
				defer delete(ctx.ptrtracker, ptr)
			}
			// maybe slow down
			defer func() {
				ctx.ptrlevel--
			}()
		}

		if je, ok := ctx.LoadToJSONEncoder(t.class); ok {
			dst, err = je.EncodeToJSON(ctx, dst, obj)
		} else {
			container := map[string]interface{}{
				"$class": t.class,
				"$":      t.m,
			}
			dst, err = EncodeStringInterfaceMapToJSON(ctx, dst, container)
		}

	case *JavaList:
		if ctx.maxPtrCycles != 0 {
			if ctx.ptrlevel++; ctx.ptrlevel >= ctx.maxPtrCycles {
				ptr := t
				if _, ok := ctx.ptrtracker[ptr]; ok {
					// cannot encode process CyclicReference and use {
					//    "$class": "CyclicReference",
					//    "$": "type name"
					// } replace
					return EncodeStringInterfaceMapToJSON(ctx, dst, map[string]interface{}{
						"$class": "CyclicReference",
						"$":      "JavaObject",
					})
				}
				ctx.ptrtracker[ptr] = struct{}{}
				defer delete(ctx.ptrtracker, ptr)
			}
			// maybe slow down
			defer func() {
				ctx.ptrlevel--
			}()
		}

		if je, ok := ctx.LoadToJSONEncoder(t.class); ok {
			dst, err = je.EncodeToJSON(ctx, dst, obj)
		} else {
			container := map[string]interface{}{
				"$class": t.class,
				"$":      t.value,
			}
			dst, err = EncodeStringInterfaceMapToJSON(ctx, dst, container)
		}

	default:
		dst, err = EncodeValueToJSON(ctx, dst, reflect.ValueOf(obj))
	}

	ctx.depth--

	return dst, err
}

// EncodeValueToJSON encodes reflect.Value to dst.
func EncodeValueToJSON(ctx *JSONEncodeContext, dst []byte, value reflect.Value) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodevalue")
		defer ctx.tracer.OnTraceStop("json encodevalue")
	}

	if value.Kind() == reflect.Invalid {
		return dst, ErrEncodeCannotInvalidValue
	} else if value.Kind() == reflect.Ptr && value.IsNil() {
		return EncodeNilToJSON(ctx, dst)
	}

	if value.Kind() == reflect.Ptr && ctx.maxPtrCycles != 0 {
		if ctx.ptrlevel++; ctx.ptrlevel >= ctx.maxPtrCycles {
			if value.CanInterface() {
				ptr := value.Interface()
				if _, ok := ctx.ptrtracker[ptr]; ok {
					// cannot encode process CyclicReference and use {
					//    "$class": "CyclicReference",
					//    "$": "type name"
					// } replace
					return EncodeStringInterfaceMapToJSON(ctx, dst, map[string]interface{}{
						"$class": "CyclicReference",
						"$":      value.Type().String(),
					})
				}
				ctx.ptrtracker[ptr] = struct{}{}
				defer delete(ctx.ptrtracker, ptr)
			}
		}
		// maybe slow down
		defer func() {
			ctx.ptrlevel--
		}()
	}

	// Unwrap pointer if needs
	value = reflect.Indirect(value)

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		if value.CanInterface() {
			return EncodeInterfaceSliceValueToJSON(ctx, dst, value)
		}
		return dst, ErrEncodeSliceCannotBeInterfaced

	case reflect.Map:
		if value.CanInterface() {
			return EncodeStringInterfaceMapValueToJSON(ctx, dst, value)
		}
		return dst, ErrEncodeMapCannotBeInterfaced

	case reflect.Struct:
		if value.CanInterface() {
			return EncodeStructValueToJSON(ctx, dst, reflect.TypeOf(value.Interface()), value)
		}
		return dst, ErrEncodeStructCannotBeInterfaced

	case reflect.Ptr:
		for {
			value = reflect.Indirect(value)
			if value.Kind() != reflect.Ptr {
				break
			}
		}

		if value.CanInterface() {
			return EncodeToJSON(ctx, dst, value.Interface())
		}

		return dst, ErrEncodePtrCannotBeInterfaced

	case reflect.Bool:
		return EncodeBoolToJSON(ctx, dst, value.Bool())
	case reflect.Int:
		return EncodeInt64ToJSON(ctx, dst, value.Int())
	case reflect.Int8:
		return EncodeInt32ToJSON(ctx, dst, int32(value.Int()))
	case reflect.Int16:
		return EncodeInt32ToJSON(ctx, dst, int32(value.Int()))
	case reflect.Int32:
		return EncodeInt32ToJSON(ctx, dst, int32(value.Int()))
	case reflect.Int64:
		return EncodeInt64ToJSON(ctx, dst, value.Int())
	case reflect.Uint:
		return EncodeInt64ToJSON(ctx, dst, int64(value.Uint()))
	case reflect.Uint8:
		return EncodeInt32ToJSON(ctx, dst, int32(value.Uint()))
	case reflect.Uint16:
		return EncodeInt32ToJSON(ctx, dst, int32(value.Uint()))
	case reflect.Uint32:
		return EncodeInt32ToJSON(ctx, dst, int32(value.Uint()))
	case reflect.Uint64:
		return EncodeInt64ToJSON(ctx, dst, int64(value.Uint()))
	case reflect.Float32:
		return EncodeFloat64ToJSON(ctx, dst, value.Float())
	case reflect.Float64:
		return EncodeFloat64ToJSON(ctx, dst, value.Float())
	case reflect.String:
		return EncodeStringToJSON(ctx, dst, value.String())
	case reflect.Interface:
		return EncodeToJSON(ctx, dst, value.Elem())
	case reflect.Uintptr:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.UnsafePointer:
		fallthrough
	default:
		return dst, errors.New("hessian: cannot encode type " + value.Kind().String() + " to json")
	}
}

// EncodeInterfaceSliceValueToJSON encodes reflect.Value to dst.
func EncodeInterfaceSliceValueToJSON(ctx *JSONEncodeContext, dst []byte, value reflect.Value) ([]byte, error) {
	classname := getInterfaceName(value.Interface())
	if je, ok := ctx.LoadToJSONEncoder(classname); ok {
		return je.EncodeToJSON(ctx, dst, value.Interface())
	}

	slice := make([]interface{}, value.Len())
	for i := 0; i < value.Len(); i++ {
		if vi := value.Index(i); vi.CanInterface() {
			slice[i] = vi.Interface()
		} else {
			return dst, ErrEncodeSliceElemCannotBeInterfaced
		}
	}

	if classname == "" {
		return EncodeInterfaceSliceToJSON(ctx, dst, slice)
	}

	return EncodeStringInterfaceMapToJSON(ctx, dst, map[string]interface{}{
		"$class": classname,
		"$":      slice,
	})
}

// EncodeInterfaceSliceToJSON encodes reflect.Value to dst.
func EncodeInterfaceSliceToJSON(ctx *JSONEncodeContext, dst []byte, m []interface{}) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encode[]interface{}")
		defer ctx.tracer.OnTraceStop("json encode[]interface{}")
	}

	var (
		err error
		n   = len(m)
	)

	dst = append(dst, '[')
	for i := range m {
		dst, err = EncodeToJSON(ctx, dst, m[i])
		if err != nil {
			return dst, err
		}

		if n = n - 1; n != 0 {
			dst = append(dst, ',')
		}
	}
	dst = append(dst, ']')
	return dst, nil
}

// EncodeInterfaceInterfaceMapToJSON encodes reflect.Value to dst.
func EncodeInterfaceInterfaceMapToJSON(ctx *JSONEncodeContext, dst []byte,
	m map[interface{}]interface{}) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodemap[interface{}]interface{}")
		defer ctx.tracer.OnTraceStop("json encodemap[interface{}]interface{}")
	}

	var (
		err error
		n   = len(m)
	)
	dst = append(dst, '{')

	if ctx.less == nil {
		for k := range m {
			key, ok := k.(string)
			if !ok {
				return dst, errors.New("hessian: expect encode map[string]interface{}")
			}
			if dst, err = EncodeStringToJSON(ctx, dst, key); err != nil {
				return dst, err
			}

			dst = append(dst, ':')
			dst, err = EncodeToJSON(ctx, dst, m[k])
			if err != nil {
				return dst, err
			}

			if n = n - 1; n != 0 {
				dst = append(dst, ',')
			}
		}
	} else {
		sorted := make([]string, 0, n)
		for k := range m {
			key, ok := k.(string)
			if !ok {
				return dst, errors.New("hessian: expect encode map[string]interface{}")
			}
			sorted = append(sorted, key)
		}

		sort.Slice(sorted, func(i, j int) bool {
			ii := sorted[i]
			keyi := ii
			valuei := m[ii]
			jj := sorted[j]
			keyj := jj
			valuej := m[jj]
			return ctx.less(keyi, keyj, valuei, valuej)
		})

		n = len(sorted)
		for i := range sorted {
			k := sorted[i]
			if dst, err = EncodeStringToJSON(ctx, dst, k); err != nil {
				return dst, err
			}

			dst = append(dst, ':')
			dst, err = EncodeToJSON(ctx, dst, m[k])
			if err != nil {
				return dst, err
			}

			if n = n - 1; n != 0 {
				dst = append(dst, ',')
			}
		}
	}

	dst = append(dst, '}')
	return dst, nil
}

// EncodeStructValueToJSON encodes reflect.Value to dst.
func EncodeStructValueToJSON(ctx *JSONEncodeContext, dst []byte, typ reflect.Type,
	value reflect.Value) ([]byte, error) {
	classname := getInterfaceName(value.Interface())
	if je, ok := ctx.LoadToJSONEncoder(classname); ok {
		return je.EncodeToJSON(ctx, dst, value.Interface())
	}

	m := acquireMStringInterface()

	for i := 0; i < value.NumField(); i++ {
		tfield := typ.Field(i)
		vfield := value.Field(i)
		if !vfield.CanInterface() {
			continue
		}
		if vfield.CanInterface() {
			name := tfield.Name
			if hname := tfield.Tag.Get("hessian"); hname != "" {
				if hname == "-" {
					continue
				}
				name = hname
			}
			(*m)[name] = vfield.Interface()
		}
	}

	var err error

	if classname == "" {
		dst, err = EncodeStringInterfaceMapToJSON(ctx, dst, *m)
	} else {
		dst, err = EncodeStringInterfaceMapToJSON(ctx, dst, map[string]interface{}{
			"$class": classname,
			"$":      m,
		})
	}

	releaseMStringInterface(m)

	return dst, err
}

// EncodeStringInterfaceMapValueToJSON encodes reflect.Value to dst.
func EncodeStringInterfaceMapValueToJSON(ctx *JSONEncodeContext, dst []byte, value reflect.Value) ([]byte, error) {
	classname := getInterfaceName(value.Interface())
	if je, ok := ctx.LoadToJSONEncoder(classname); ok {
		return je.EncodeToJSON(ctx, dst, value.Interface())
	}

	keys := value.MapKeys()
	m := make(map[interface{}]interface{}, len(keys))
	for i := range keys {
		key := keys[i]
		value = value.MapIndex(key)
		if key.CanInterface() && value.CanInterface() { // Fast path
			m[key.Interface()] = value.Interface()
		} else {
			return dst, fmt.Errorf("hessian: cannot encode type m[%s]%s", key.Kind(), value.Kind())
		}
	}

	if classname == "" {
		return EncodeInterfaceInterfaceMapToJSON(ctx, dst, m)
	}

	return EncodeStringInterfaceMapToJSON(ctx, dst, map[string]interface{}{
		"$class": classname,
		"$":      m,
	})
}

// EncodeStringInterfaceMapToJSON encodes map[string]interface{} to dst.
func EncodeStringInterfaceMapToJSON(ctx *JSONEncodeContext, dst []byte, m map[string]interface{}) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodemap[string]interface{}")
		defer ctx.tracer.OnTraceStop("json encodemap[string]interface{}")
	}

	var (
		err error
		n   = len(m)
	)
	dst = append(dst, '{')

	if ctx.less != nil { // sort keys
		sorted := make([]string, 0, n)
		for i := range m {
			sorted = append(sorted, i)
		}

		sort.Slice(sorted, func(i, j int) bool {
			ii := sorted[i]
			keyi := ii
			valuei := m[ii]
			jj := sorted[j]
			keyj := jj
			valuej := m[jj]
			return ctx.less(keyi, keyj, valuei, valuej)
		})

		for i := range sorted {
			k := sorted[i]
			if dst, err = EncodeStringToJSON(ctx, dst, k); err != nil {
				return dst, err
			}

			dst = append(dst, ':')
			dst, err = EncodeToJSON(ctx, dst, m[k])
			if err != nil {
				return dst, err
			}

			if n = n - 1; n != 0 {
				dst = append(dst, ',')
			}
		}

	} else {
		for k := range m {
			if dst, err = EncodeStringToJSON(ctx, dst, k); err != nil {
				return dst, err
			}

			dst = append(dst, ':')
			dst, err = EncodeToJSON(ctx, dst, m[k])
			if err != nil {
				return dst, err
			}

			if n = n - 1; n != 0 {
				dst = append(dst, ',')
			}
		}
	}

	dst = append(dst, '}')
	return dst, nil
}

// EncodeStringToJSON encodes string to dst.
func EncodeStringToJSON(ctx *JSONEncodeContext, dst []byte, s string) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodestring")
		defer ctx.tracer.OnTraceStop("json encodestring")
	}
	dst = escapeString(dst, s)
	return dst, nil
}

// EncodeNilToJSON encodes nil to dst.
func EncodeNilToJSON(ctx *JSONEncodeContext, dst []byte) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodenil")
		defer ctx.tracer.OnTraceStop("json encodenil")
	}
	dst = append(dst, "null"...)
	return dst, nil
}

// EncodeBinaryToJSON encodes []byte to dst.
func EncodeBinaryToJSON(ctx *JSONEncodeContext, dst []byte, b []byte) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodebinary")
		defer ctx.tracer.OnTraceStop("json encodebinary")
	}
	n := base64.StdEncoding.EncodedLen(len(b))
	b64 := make([]byte, n)
	base64.StdEncoding.Encode(b64, b)
	dst = append(dst, '"')
	dst = append(dst, b64...)
	dst = append(dst, '"')
	return dst, nil
}

// EncodeInt32ToJSON encodes in32 to dst.
func EncodeInt32ToJSON(ctx *JSONEncodeContext, dst []byte, i int32) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodeint32")
		defer ctx.tracer.OnTraceStop("json encodeint32")
	}
	dst = strconv.AppendInt(dst, int64(i), 10)
	return dst, nil
}

// EncodeBoolToJSON encodes bool to dst.
func EncodeBoolToJSON(ctx *JSONEncodeContext, dst []byte, b bool) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodebool")
		defer ctx.tracer.OnTraceStop("json encodebool")
	}
	if b == true {
		return append(dst, "true"...), nil
	}
	return append(dst, "false"...), nil
}

// EncodeInt64ToJSON encodes in64 to dst.
func EncodeInt64ToJSON(ctx *JSONEncodeContext, dst []byte, i int64) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodeint64")
		defer ctx.tracer.OnTraceStop("json encodeint64")
	}
	dst = strconv.AppendInt(dst, i, 10)
	return dst, nil
}

// EncodeFloat64ToJSON encodes float64 to dst.
func EncodeFloat64ToJSON(ctx *JSONEncodeContext, dst []byte, f float64) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodefloat64")
		defer ctx.tracer.OnTraceStop("json encodefloat64")
	}
	// Convert as if by ES6 number to string conversion.
	// This matches most other JSON generators.
	// See golang.org/issue/6384 and golang.org/issue/14135.
	// Like fmt %g, but the exponent cutoffs are different
	// and exponents themselves are not padded to two digits.
	abs := math.Abs(f)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			fmt = 'e'
		}
	}

	dst = strconv.AppendFloat(dst, f, fmt, -1, 64)
	if fmt == 'e' {
		// clean up e-09 to e-9
		n := len(dst)
		if n >= 4 && dst[n-4] == 'e' && dst[n-3] == '-' && dst[n-2] == '0' {
			dst[n-2] = dst[n-1]
			dst = dst[:n-1]
		}
	}
	return dst, nil
}

// EncodeFloat32ToJSON encodes float32 to dst.
func EncodeFloat32ToJSON(ctx *JSONEncodeContext, dst []byte, f float32) ([]byte, error) {
	if ctx.tracer != nil {
		ctx.tracer.OnTraceStart("json encodefloat32")
		defer ctx.tracer.OnTraceStop("json encodefloat32")
	}
	// Convert as if by ES6 number to string conversion.
	// This matches most other JSON generators.
	// See golang.org/issue/6384 and golang.org/issue/14135.
	// Like fmt %g, but the exponent cutoffs are different
	// and exponents themselves are not padded to two digits.
	abs := math.Abs(float64(f))
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if float32(abs) < 1e-6 || float32(abs) >= 1e21 {
			fmt = 'e'
		}
	}

	dst = strconv.AppendFloat(dst, float64(f), fmt, -1, 32)
	if fmt == 'e' {
		// clean up e-09 to e-9
		n := len(dst)
		if n >= 4 && dst[n-4] == 'e' && dst[n-3] == '-' && dst[n-2] == '0' {
			dst[n-2] = dst[n-1]
			dst = dst[:n-1]
		}
	}
	return dst, nil
}
