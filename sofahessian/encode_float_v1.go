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

// EncodeFloat64ToHessianV1 encodes float64 to dst with hessian3 v2 protocol.
func EncodeFloat64ToHessianV1(o *EncodeContext, dst []byte, d float64) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodefloat64")
		defer o.tracer.OnTraceStop("encodefloat64")
	}

	dst = append(dst, "D00000000"...)
	binary.BigEndian.PutUint64(dst[len(dst)-8:], math.Float64bits(d))

	return dst, nil
}
