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

import "encoding/binary"

func EncodeBinaryToHessian3V2(o *EncodeContext, dst []byte, b []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebinary")
		defer o.tracer.OnTraceStop("encodebinary")
	}

	n := len(b)
	if n < 16 {
		dst = append(dst, uint8(n)+0x20)
		dst = append(dst, b...)
		return dst, nil
	}

	for n > 0x8000 {
		dst = append(dst, "b00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], 0x8000)
		dst = append(dst, b[:0x8000]...)
		b = b[0x8000:]
		n = len(b)
	}

	if n < 16 {
		dst = append(dst, uint8(n+0x20))
	} else {
		dst = append(dst, "B00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(n))
	}
	dst = append(dst, b...)

	return dst, nil
}
