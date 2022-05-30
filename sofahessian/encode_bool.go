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

// EncodeBoolToHessian4V2 encodes bool to dst.
// The octet 'F' represents false and the octet T represents true.
// boolean ::= T
//         ::= F
func EncodeBoolToHessian4V2(o *EncodeContext, dst []byte, b bool) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebool")
		defer o.tracer.OnTraceStop("encodebool")
	}

	if b {
		dst = append(dst, 'T')
	} else {
		dst = append(dst, 'F')
	}
	return dst, nil
}

func EncodeBoolToHessian3V2(o *EncodeContext, dst []byte, b bool) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebool")
		defer o.tracer.OnTraceStop("encodebool")
	}

	if b {
		dst = append(dst, 'T')
	} else {
		dst = append(dst, 'F')
	}
	return dst, nil
}

func EncodeBoolToHessianV1(o *EncodeContext, dst []byte, b bool) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebool")
		defer o.tracer.OnTraceStop("encodebool")
	}

	if b {
		dst = append(dst, 'T')
	} else {
		dst = append(dst, 'F')
	}
	return dst, nil
}
