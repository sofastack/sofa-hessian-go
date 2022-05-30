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
	"time"

	"github.com/valyala/fastjson"
)

// Encoder writes hessian type to underlying buffer.
type Encoder struct {
	b   []byte
	ctx *EncodeContext
}

// NewEncoder returns the encoder.
func NewEncoder(o *EncodeContext) *Encoder {
	return &Encoder{
		ctx: o,
	}
}

// GetContext reads the context of encoding.
func (e *Encoder) GetContext() *EncodeContext {
	return e.ctx
}

// EncodeStreamingJSONBytes encodes the json stream to encoder.
func (e *Encoder) EncodeStreamingJSONBytes(ctx *JSONEncodeContext, b []byte) error {
	var sc fastjson.Scanner
	sc.InitBytes(b)

	// streaming parse
	for sc.Next() {
		if err := e.EncodeFastJSON(ctx, sc.Value()); err != nil {
			return err
		}
	}

	return sc.Error()
}

// EncodeFastJSON encodes fastjson to dst.
func (e *Encoder) EncodeFastJSON(ctx *JSONEncodeContext, j *fastjson.Value) error {
	return EncodeJSONToHessian(ctx, e, j)
}

// EncodeJSONBytes encodes JSON bytes to encoder.
func (e *Encoder) EncodeJSONBytes(ctx *JSONEncodeContext, b []byte) error {
	return EncodeRawJSONBytesToHessian(ctx, e, b)
}

// EncodeJSONString encodes JSON string to encoder.
func (e *Encoder) EncodeJSONString(ctx *JSONEncodeContext, b string) error {
	return EncodeRawJSONStringToHessian(ctx, e, b)
}

// EncodeBinary encodes []byte to encoder.
func (e *Encoder) EncodeBinary(b []byte) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeBinaryToHessian4V2(e.ctx, e.b, b)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeBinaryToHessianV1(e.ctx, e.b, b)
	} else {
		e.b, err = EncodeBinaryToHessian3V2(e.ctx, e.b, b)
	}
	return err
}

// EncodeValue encodes reflect.Value to encoder.
func (e *Encoder) EncodeValue(value reflect.Value) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeValueToHessian4V2(e.ctx, e.b, value)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeValueToHessianV1(e.ctx, e.b, value)
	} else {
		e.b, err = EncodeValueToHessian3v2(e.ctx, e.b, value)
	}
	return err
}

// EncodeBool encodes bool to encoder.
func (e *Encoder) EncodeBool(b bool) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeBoolToHessian4V2(e.ctx, e.b, b)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeBoolToHessianV1(e.ctx, e.b, b)
	} else {
		e.b, err = EncodeBoolToHessian3V2(e.ctx, e.b, b)
	}
	return err
}

// EncodeDate encodes time.Time to encoder.
func (e *Encoder) EncodeDate(t time.Time) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeDateToHessian4V2(e.ctx, e.b, t)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeDateToHessianV1(e.ctx, e.b, t)
	} else {
		e.b, err = EncodeDateToHessian3V2(e.ctx, e.b, t)
	}
	return err
}

// EncodeFloat64 encodes float64 to encoder.
func (e *Encoder) EncodeFloat64(f64 float64) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeFloat64ToHessian4V2(e.ctx, e.b, f64)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeFloat64ToHessianV1(e.ctx, e.b, f64)
	} else {
		e.b, err = EncodeFloat64ToHessian3V2(e.ctx, e.b, f64)
	}
	return err
}

// EncodeInt64 encodes int64 to encoder.
func (e *Encoder) EncodeInt64(i64 int64) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeInt64ToHessian4V2(e.ctx, e.b, i64)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeInt64ToHessianV1(e.ctx, e.b, i64)
	} else {
		e.b, err = EncodeInt64ToHessian3V2(e.ctx, e.b, i64)
	}
	return err
}

// EncodeInt32 encodes int32 to encoder.
func (e *Encoder) EncodeInt32(i32 int32) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeInt32ToHessian4V2(e.ctx, e.b, i32)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeInt32ToHessianV1(e.ctx, e.b, i32)
	} else {
		e.b, err = EncodeInt32ToHessian3V2(e.ctx, e.b, i32)
	}
	return err
}

// EncodeList encodes []T to encoder.
func (e *Encoder) EncodeList(slice interface{}) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeListToHessian4V2(e.ctx, e.b, slice)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeListToHessianV1(e.ctx, e.b, slice)
	} else {
		e.b, err = EncodeListToHessian3V2(e.ctx, e.b, slice)
	}
	return err
}

// EncodeListBegin encodes list prefix to encoder.
func (e *Encoder) EncodeListBegin(length int, typ string) (end bool, err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, end, err = EncodeListBeginToHessian4V2(e.ctx, e.b, length, typ)
	} else if e.ctx.version == HessianV1 {
		e.b, end, err = EncodeListBeginToHessianV1(e.ctx, e.b, length, typ)
	} else {
		e.b, end, err = EncodeListBeginToHessian3V2(e.ctx, e.b, length, typ)
	}
	return end, err
}

// EncodeListEnd encodes list end to encoder.
func (e *Encoder) EncodeListEnd(end bool) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeListEndToHessian4V2(e.ctx, e.b, end)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeListEndToHessianV1(e.ctx, e.b, end)
	} else {
		e.b, err = EncodeListEndToHessian3V2(e.ctx, e.b, end)
	}
	return err
}

// EncodeMapBegin encodes map prefix to encoder.
func (e *Encoder) EncodeMapBegin() (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeMapBeginToHessian4V2(e.ctx, e.b, "")
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeMapBeginToHessianV1(e.ctx, e.b, "")
	} else {
		e.b, err = EncodeMapBeginToHessian3V2(e.ctx, e.b, "")
	}
	return err
}

// EncodeMapEnd encodes map end to encoder.
func (e *Encoder) EncodeMapEnd() (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeMapEndToHessian4V2(e.ctx, e.b)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeMapEndToHessianV1(e.ctx, e.b)
	} else {
		e.b, err = EncodeMapEndToHessian3V2(e.ctx, e.b)
	}
	return err
}

// EncodeMap encodes map[T]T to encoder.
func (e *Encoder) EncodeMap(m interface{}) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeMapToHessian4V2(e.ctx, e.b, m)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeMapToHessianV1(e.ctx, e.b, m)
	} else {
		e.b, err = EncodeMapToHessian3V2(e.ctx, e.b, m)
	}
	return err
}

// EncodeNil encodes nil to encoder.
func (e *Encoder) EncodeNil() (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeNilToHessian4V2(e.ctx, e.b)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeNilToHessianV1(e.ctx, e.b)
	} else {
		e.b, err = EncodeNilToHessian3V2(e.ctx, e.b)
	}
	return err
}

// EncodeObject encodes object to encoder.
func (e *Encoder) EncodeObject(obj interface{}) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeObjectToHessian4V2(e.ctx, e.b, obj)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeObjectToHessianV1(e.ctx, e.b, obj)
	} else {
		e.b, err = EncodeObjectToHessian3V2(e.ctx, e.b, obj)
	}
	return err
}

// EncodeString encodes string to encoder.
func (e *Encoder) EncodeString(s string) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeStringToHessian4V2(e.ctx, e.b, s)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeStringToHessianV1(e.ctx, e.b, s)
	} else {
		e.b, err = EncodeStringToHessian3V2(e.ctx, e.b, s)
	}
	return err
}

func (e *Encoder) EncodeClassDefinition(typ string, fields []string) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeClassDefinitionToHessian4V2(e.ctx, e.b, typ, fields)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeClassDefinitionToHessianV1(e.ctx, e.b, typ, fields)
	} else {
		e.b, err = EncodeClassDefinitionToHessian3V2(e.ctx, e.b, typ, fields)
	}
	return err
}

// Encode encodes the interface to encoder.
//
// Encode cannot represent cyclic data structures. Passing cyclic data structures will result in panic.
func (e *Encoder) Encode(i interface{}) (err error) {
	if e.ctx.version == Hessian4xV2 {
		e.b, err = EncodeToHessian4V2(e.ctx, e.b, i)
	} else if e.ctx.version == HessianV1 {
		e.b, err = EncodeToHessianV1(e.ctx, e.b, i)
	} else {
		e.b, err = EncodeToHessian3V2(e.ctx, e.b, i)
	}
	return err
}

// EncodeBytes encodes bytes to encoder.
func (e *Encoder) EncodeBytes(b []byte) {
	e.b = append(e.b, b...)
}

// Bytes returns the encoded bytes.
func (e *Encoder) Bytes() []byte {
	return e.b
}

// Reset resets the encoder status.
func (e *Encoder) Reset() {
	e.b = e.b[:0]
	e.ctx = NewEncodeContext()
}

// ResetWithContext resets the encoder status with context.
func (e *Encoder) ResetWithContext(o *EncodeContext) {
	e.b = e.b[:0]
	e.ctx = o
}
