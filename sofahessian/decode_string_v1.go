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

func DecodeStringHessianV1(o *DecodeContext, reader *bufio.Reader) (string, error) {
	var (
		b   []byte
		err error
	)

	b, err = DecodeStringToHessianV1(o, reader, b)
	if err != nil {
		return "", err
	}
	return string(b), err
}

// DecodeStringToHessianV1 decodes dst to string.
func DecodeStringToHessianV1(o *DecodeContext, reader *bufio.Reader, s []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodestring")
		defer o.tracer.OnTraceStop("decodestring")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return s, err
	}

	for c1 == 's' {
		s, err = readLenAndUTF8StringFromReader(reader, s)
		if err != nil {
			return s, err
		}

		c1, err = reader.ReadByte()
		if err != nil {
			return s, err
		}
	}

	if c1 >= 0x00 && c1 <= 0x1F {
		s, err = readUTF8StringFromReader(reader, s, int(c1))
	} else if c1 == 0x53 {
		s, err = readLenAndUTF8StringFromReader(reader, s)
	} else if c1 == 0x74 { // hessian1
		s, err = readLenAndUTF8StringFromReader(reader, s)
	} else {
		return s, ErrDecodeMalformedString
	}

	return s, err
}
