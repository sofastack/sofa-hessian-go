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
	"math"
)

// EncodeFloat64ToHessian4V2 encodes float64 to dst.
// A 64-bit IEEE floating pointer number.
//
// double ::= D b7 b6 b5 b4 b3 b2 b1 b0
//        ::= x5b
//        ::= x5c
//        ::= x5d b0
//        ::= x5e b1 b0
//        ::= x5f b3 b2 b1 b0
func EncodeFloat64ToHessian4V2(o *EncodeContext, dst []byte, d float64) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodefloat64")
		defer o.tracer.OnTraceStop("encodefloat64")
	}

	truncate := int64(d)
	if float64(truncate) == d {
		switch {
		case truncate == 0:
			dst = append(dst, 0x5b)
			return dst, nil
		case truncate == 1:
			dst = append(dst, 0x5c)
			return dst, nil
		case -0x80 <= truncate && truncate < 0x80:
			dst = append(dst, 0x5d, uint8(truncate))
			return dst, nil
		case -0x8000 <= truncate && truncate < 0x8000:
			dst = append(dst, 0x5e, 0, 0)
			binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(truncate))
			return dst, nil
		}
	}

	mills := (int64)(d * 1000)
	if mills >= math.MinInt32 && mills <= math.MaxInt32 && 0.001*float64(mills) == d {
		dst = append(dst, 0x5F, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(dst[len(dst)-4:], uint32(mills))
		return dst, nil
	}

	dst = append(dst, "D00000000"...)
	binary.BigEndian.PutUint64(dst[len(dst)-8:], math.Float64bits(d))
	return dst, nil
}
