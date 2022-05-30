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
)

func DecodeHessianV1(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	codes, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}

	if o.maxdepth > 0 {
		o.addDepth()
		if o.depth > o.maxdepth {
			return nil, ErrDecodeMaxDepthExceeded
		}
		defer o.subDepth()
	}

	switch codes[0] {
	case 66:
		return DecodeBinaryHessianV1(o, reader)
	case 68:
		return DecodeFloat64HessianV1(o, reader)
	case 70:
		return DecodeBoolHessianV1(o, reader)
	case 73:
		return DecodeInt32HessianV1(o, reader)
	case 76:
		return DecodeInt64HessianV1(o, reader)
	case 77:
		return DecodeObjectHessianV1(o, reader)
	case 78:
		return nil, DecodeNilHessianV1(o, reader)
	case 82:
		return DecodeRefHessianV1(o, reader)
	case 83:
		return DecodeStringHessianV1(o, reader)
	case 84:
		return DecodeBoolHessianV1(o, reader)
	case 86:
		return DecodeListHessianV1(o, reader)
	case 98:
		return DecodeBinaryHessianV1(o, reader)
	case 100:
		return DecodeDateHessianV1(o, reader)
	case 108:
		return readUint32FromReader(reader)
	case 115:
		return DecodeStringHessianV1(o, reader)
	case 116:
		return DecodeTypeHessianV1(o, reader)
	}

	return nil, ErrDecodeUnknownEncoding
}
