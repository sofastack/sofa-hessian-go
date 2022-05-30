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
	"reflect"
	"time"
)

// EncodeValueToHessianV1 encodes reflect.Value to dst.
func EncodeValueToHessianV1(o *EncodeContext, dst []byte, value reflect.Value) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodevalue")
		defer o.tracer.OnTraceStop("encodevalue")
	}

	if value.Kind() == reflect.Invalid {
		return dst, ErrEncodeCannotInvalidValue
	} else if value.Kind() == reflect.Ptr && value.IsNil() {
		return EncodeNilToHessianV1(o, dst)
	}

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		if value.CanInterface() {
			return EncodeListToHessianV1(o, dst, value.Interface())
		}
		return dst, ErrEncodeSliceCannotBeInterfaced

	case reflect.Map:
		if value.CanInterface() {
			return EncodeMapToHessianV1(o, dst, value.Interface())
		}
		return dst, ErrEncodeMapCannotBeInterfaced

	case reflect.Struct:
		if value.CanInterface() {
			return EncodeObjectToHessianV1(o, dst, value.Interface())
		}
		return dst, ErrEncodeStructCannotBeInterfaced

	case reflect.Ptr:
		// **T => *T
		indir := value.Elem()
		if indir.Kind() == reflect.Struct {
			return EncodeObjectToHessianV1(o, dst, value.Interface())
		}

		if value.CanInterface() {
			return EncodeToHessianV1(o, dst, indir.Interface())
		}

		return dst, ErrEncodePtrCannotBeInterfaced

	case reflect.Bool:
		return EncodeBoolToHessianV1(o, dst, value.Bool())
	case reflect.Int:
		return EncodeInt64ToHessianV1(o, dst, value.Int())
	case reflect.Int8:
		return EncodeInt32ToHessianV1(o, dst, int32(value.Int()))
	case reflect.Int16:
		return EncodeInt32ToHessianV1(o, dst, int32(value.Int()))
	case reflect.Int32:
		return EncodeInt32ToHessianV1(o, dst, int32(value.Int()))
	case reflect.Int64:
		return EncodeInt64ToHessianV1(o, dst, value.Int())
	case reflect.Uint:
		return EncodeInt64ToHessianV1(o, dst, int64(value.Uint()))
	case reflect.Uint8:
		return EncodeInt32ToHessianV1(o, dst, int32(value.Uint()))
	case reflect.Uint16:
		return EncodeInt32ToHessianV1(o, dst, int32(value.Uint()))
	case reflect.Uint32:
		return EncodeInt32ToHessianV1(o, dst, int32(value.Uint()))
	case reflect.Uint64:
		return EncodeInt64ToHessianV1(o, dst, int64(value.Uint()))
	case reflect.Float32:
		return EncodeFloat64ToHessianV1(o, dst, value.Float())
	case reflect.Float64:
		return EncodeFloat64ToHessianV1(o, dst, value.Float())
	case reflect.String:
		return EncodeStringToHessianV1(o, dst, value.String())
	case reflect.Interface:
		return EncodeToHessianV1(o, dst, value.Elem())
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
		return dst, errors.New("hessian: cannot encode type " + value.Kind().String())
	}
}

// EncodeHessianV1 encodes the interface to dst.
func EncodeHessianV1(o *EncodeContext, value interface{}) ([]byte, error) {
	return EncodeToHessianV1(o, nil, value)
}

// EncodeToHessianV1 encodes the interface to dst.
func EncodeToHessianV1(o *EncodeContext, dst []byte, value interface{}) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encode")
		defer o.tracer.OnTraceStop("encode")
	}

	if o.maxdepth > 0 {
		o.addDepth()
		if o.depth > o.maxdepth {
			return nil, ErrEncodeMaxDepthExceeded
		}
		defer o.subDepth()
	}

	switch v := value.(type) {
	// Fast path without reflection
	case HessianEncoder:
		return v.HessianEncode(o, dst)
	case *[]byte:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeBinaryToHessianV1(o, dst, *v)
	case []byte:
		return EncodeBinaryToHessianV1(o, dst, v)

	case *string:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeStringToHessianV1(o, dst, *v)
	case string:
		return EncodeStringToHessianV1(o, dst, v)

	case int:
		return EncodeInt64ToHessianV1(o, dst, int64(v))
	case *int:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeInt64ToHessianV1(o, dst, int64(*v))
	case uint:
		return EncodeInt64ToHessianV1(o, dst, int64(v))
	case *uint:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeInt64ToHessianV1(o, dst, int64(*v))
	case uint8:
		return EncodeInt32ToHessianV1(o, dst, int32(v))
	case int8:
		return EncodeInt32ToHessianV1(o, dst, int32(v))
	case uint16:
		return EncodeInt32ToHessianV1(o, dst, int32(v))
	case int16:
		return EncodeInt32ToHessianV1(o, dst, int32(v))
	case uint32:
		return EncodeInt32ToHessianV1(o, dst, int32(v))
	case int32:
		return EncodeInt32ToHessianV1(o, dst, v)
	case uint64:
		return EncodeInt64ToHessianV1(o, dst, int64(v))
	case int64:
		return EncodeInt64ToHessianV1(o, dst, v)

	case *uint8:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeInt32ToHessianV1(o, dst, int32(*v))
	case *int8:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeInt32ToHessianV1(o, dst, int32(*v))
	case *uint16:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeInt32ToHessianV1(o, dst, int32(*v))
	case *int16:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeInt32ToHessianV1(o, dst, int32(*v))
	case *uint32:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeInt64ToHessianV1(o, dst, int64(*v))
	case *int32:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeInt32ToHessianV1(o, dst, *v)
	case *uint64:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeInt64ToHessianV1(o, dst, int64(*v))
	case *int64:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeInt64ToHessianV1(o, dst, *v)

	case *float32:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeFloat64ToHessianV1(o, dst, float64(*v))
	case float32:
		return EncodeFloat64ToHessianV1(o, dst, float64(v))
	case *float64:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeFloat64ToHessianV1(o, dst, *v)
	case float64:
		return EncodeFloat64ToHessianV1(o, dst, v)

	case *bool:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeBoolToHessianV1(o, dst, *v)
	case bool:
		return EncodeBoolToHessianV1(o, dst, v)

	case nil:
		return EncodeNilToHessianV1(o, dst)

	case *time.Time:
		if v == nil {
			return EncodeNilToHessianV1(o, dst)
		}
		return EncodeDateToHessianV1(o, dst, *v)
	case time.Time:
		return EncodeDateToHessianV1(o, dst, v)

	default:
		return EncodeValueToHessianV1(o, dst, reflect.ValueOf(v))
	}
}
