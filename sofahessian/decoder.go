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
	"io"
	"time"
)

// Decoder reads the hessian type from reader.
type Decoder struct {
	ctx    *DecodeContext
	reader *bufio.Reader
}

func NewDecoder(ctx *DecodeContext, reader io.Reader) *Decoder {
	return &Decoder{
		ctx:    ctx,
		reader: bufio.NewReader(reader),
	}
}

func (d *Decoder) GetContext() *DecodeContext {
	return d.ctx
}

func (d *Decoder) GetReader() *bufio.Reader {
	return d.reader
}

func (d *Decoder) DecodeDate() (time.Time, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeDateHessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeDateHessianV1(d.ctx, d.reader)
	}
	return DecodeDateHessian3V2(d.ctx, d.reader)
}

func (d *Decoder) DecodeNil() error {
	if d.ctx.version == Hessian4xV2 {
		return DecodeNilHessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeNilHessianV1(d.ctx, d.reader)
	}
	return DecodeNilHessian3V2(d.ctx, d.reader)
}

func (d *Decoder) DecodeBinary() ([]byte, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeBinaryHessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeBinaryHessianV1(d.ctx, d.reader)
	}
	return DecodeBinaryHessian3V2(d.ctx, d.reader)
}

func (d *Decoder) DecodeString() (string, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeStringHessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeStringHessianV1(d.ctx, d.reader)
	}
	return DecodeStringHessian3V2(d.ctx, d.reader)
}

func (d *Decoder) DecodeToString(s []byte) ([]byte, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeStringToHessian4V2(d.ctx, d.reader, s)
	} else if d.ctx.version == HessianV1 {
		return DecodeStringToHessianV1(d.ctx, d.reader, s)
	}
	return DecodeStringToHessian3V2(d.ctx, d.reader, s)
}

func (d *Decoder) DecodeInt64() (int64, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeInt64Hessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeInt64HessianV1(d.ctx, d.reader)
	}
	return DecodeInt64Hessian3V2(d.ctx, d.reader)
}

func (d *Decoder) DecodeToInt64(i *int64) error {
	if d.ctx.version == Hessian4xV2 {
		return DecodeInt64ToHessian4V2(d.ctx, d.reader, i)
	} else if d.ctx.version == HessianV1 {
		return DecodeInt64ToHessianV1(d.ctx, d.reader, i)
	}
	return DecodeInt64ToHessian3V2(d.ctx, d.reader, i)
}

func (d *Decoder) DecodeInt32() (int32, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeInt32Hessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeInt32HessianV1(d.ctx, d.reader)
	}
	return DecodeInt32Hessian3V2(d.ctx, d.reader)
}

func (d *Decoder) DecodeToInt32(i *int32) error {
	if d.ctx.version == Hessian4xV2 {
		return DecodeInt32ToHessian4V2(d.ctx, d.reader, i)
	} else if d.ctx.version == HessianV1 {
		return DecodeInt32ToHessianV1(d.ctx, d.reader, i)
	}
	return DecodeInt32ToHessian3V2(d.ctx, d.reader, i)
}

func (d *Decoder) DecodeFloat64() (float64, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeFloat64Hessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeFloat64HessianV1(d.ctx, d.reader)
	}
	return DecodeFloat64Hessian3V2(d.ctx, d.reader)
}

func (d *Decoder) DecodeToFloat64(f *float64) error {
	if d.ctx.version == Hessian4xV2 {
		return DecodeFloat64ToHessian4V2(d.ctx, d.reader, f)
	} else if d.ctx.version == HessianV1 {
		return DecodeFloat64ToHessianV1(d.ctx, d.reader, f)
	}
	return DecodeFloat64ToHessian3V2(d.ctx, d.reader, f)
}

func (d *Decoder) DecodeList() (interface{}, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeListHessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeListHessianV1(d.ctx, d.reader)
	}
	return DecodeListHessian3V2(d.ctx, d.reader)
}

func (d *Decoder) DecodeMap() (interface{}, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeMapHessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeMapHessianV1(d.ctx, d.reader)
	}
	return DecodeMapHessian3V2(d.ctx, d.reader)
}

func (d *Decoder) DecodeObject() (interface{}, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeObjectHessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeObjectHessianV1(d.ctx, d.reader)
	}
	return DecodeObjectHessian3V2(d.ctx, d.reader)
}

func (d *Decoder) DecodeToObject(obj interface{}) error {
	if d.ctx.version == Hessian4xV2 {
		return DecodeObjectToHessian4V2(d.ctx, d.reader, obj)
	} else if d.ctx.version == HessianV1 {
		return DecodeObjectToHessianV1(d.ctx, d.reader, obj)
	}
	return DecodeObjectToHessian3V2(d.ctx, d.reader, obj)
}

func (d *Decoder) Decode() (interface{}, error) {
	if d.ctx.version == Hessian4xV2 {
		return DecodeHessian4V2(d.ctx, d.reader)
	} else if d.ctx.version == HessianV1 {
		return DecodeHessianV1(d.ctx, d.reader)
	}
	return DecodeHessian3V2(d.ctx, d.reader)
}

// Reset resets the decoder status.
func (d *Decoder) Reset(reader io.Reader) {
	d.reader.Reset(reader)
	d.ctx = NewDecodeContext()
}

// ResetWithContext resets the decoder status with context.
func (d *Decoder) ResetWithContext(ctx *DecodeContext, reader io.Reader) {
	d.reader.Reset(reader)
	d.ctx = ctx
}
