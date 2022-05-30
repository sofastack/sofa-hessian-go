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

// EncodeNilToHessian4V2 encodes nil to dst.
//
// null ::= N
func EncodeNilToHessian4V2(o *EncodeContext, dst []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodenil")
		defer o.tracer.OnTraceStop("encodenil")
	}
	dst = append(dst, 'N')
	return dst, nil
}

// EncodeNilToHessian3V2 encoddes nil to dst with hessian3 v2 protocol.
func EncodeNilToHessian3V2(e *EncodeContext, dst []byte) ([]byte, error) {
	return EncodeNilToHessian4V2(e, dst)
}

// EncodeNilToHessianV1 encoddes nil to dst with hessian3 v2 protocol.
func EncodeNilToHessianV1(e *EncodeContext, dst []byte) ([]byte, error) {
	return EncodeNilToHessian4V2(e, dst)
}
