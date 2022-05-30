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
	"fmt"
	"io"
)

func DecodeBoolHessian4V2(o *DecodeContext, reader io.Reader) (bool, error) {
	var b bool
	err := DecodeBoolToHessian4V2(o, reader, &b)
	return b, err
}

func DecodeBoolToHessian4V2(o *DecodeContext, reader io.Reader, b *bool) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodebool")
		defer o.tracer.OnTraceStop("decodebool")
	}

	var c [1]byte
	n, err := reader.Read(c[:])
	if err != nil {
		return err
	}
	if n < 1 {
		return fmt.Errorf("expect read one byte but got zero")
	}

	switch c[0] {
	case 'T':
		*b = true
	case 'F':
		*b = false
	default:
		return ErrDecodeMalformedBool
	}

	return nil
}
