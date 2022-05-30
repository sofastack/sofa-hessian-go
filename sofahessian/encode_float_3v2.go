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

// EncodeFloat64ToHessian3V2 encodes float64 to dst with hessian3 v2 protocol.
func EncodeFloat64ToHessian3V2(o *EncodeContext, dst []byte, d float64) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodefloat64")
		defer o.tracer.OnTraceStop("encodefloat64")
	}

	truncate := int64(d)
	if float64(truncate) == d {
		switch {
		case truncate == 0:
			dst = append(dst, 0x67)
			return dst, nil
		case truncate == 1:
			dst = append(dst, 0x68)
			return dst, nil
		case -0x80 <= truncate && truncate < 0x80:
			dst = append(dst, 0x69, uint8(truncate))
			return dst, nil
		case -0x8000 <= truncate && truncate < 0x8000:
			dst = append(dst, 0x6a, 0, 0)
			binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(truncate))
			return dst, nil
		}
	}

	f32 := float32(d)
	if float64(f32) == d {
		dst = append(dst, 0x6b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(dst[len(dst)-4:], math.Float32bits(float32(truncate)))
	} else {
		dst = append(dst, 0x44, 0, 0, 0, 0, 0, 0, 0, 0)
		binary.BigEndian.PutUint64(dst[len(dst)-8:], math.Float64bits(d))
	}

	return dst, nil
}
