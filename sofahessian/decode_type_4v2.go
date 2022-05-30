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

import "bufio"

func DecodeTypeHessian4V2(o *DecodeContext, reader *bufio.Reader) (string, error) {
	var (
		b   []byte
		err error
	)

	b, err = DecodeTypeToHessian4V2(o, reader, b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func DecodeTypeToHessian4V2(o *DecodeContext, reader *bufio.Reader, dst []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodetype")
		defer o.tracer.OnTraceStop("decodetype")
	}

	codes, err := reader.Peek(1)
	if err != nil {
		return dst, err
	}

	switch codes[0] {
	case 0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f,
		0x30, 0x31, 0x32, 0x33,
		0x52, 0x53:
		var typ []byte
		if typ, err = DecodeStringToHessian4V2(o, reader, nil); err != nil {
			return nil, err
		}
		// Copy bytes to string
		if err = o.addTyperefs(string(typ)); err != nil {
			return dst, err
		}
		dst = append(dst, typ...)
		return dst, err

	default:
		var (
			refid int32
			typ   string
		)
		refid, err = DecodeInt32Hessian4V2(o, reader)
		if err != nil {
			return dst, err
		}

		typ, err = o.getTyperefs(int(refid))
		if err != nil {
			return dst, err
		}

		dst = append(dst, typ...)
		return dst, nil
	}
}
