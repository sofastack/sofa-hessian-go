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

func DecodeTypeHessianV1(o *DecodeContext, reader *bufio.Reader) (string, error) {
	var (
		b   []byte
		err error
	)

	b, err = DecodeTypeToHessianV1(o, reader, b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func DecodeTypeToHessianV1(o *DecodeContext, reader *bufio.Reader, dst []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodetype")
		defer o.tracer.OnTraceStop("decodetype")
	}

	codes, err := reader.Peek(1)
	if err != nil {
		return dst, err
	}

	var (
		typ string
		u16 uint16
		u32 uint32
	)

	switch codes[0] {
	case 0x74:
		// nolint
		reader.ReadByte()

		u16, err = readUint16FromReader(reader)
		if err != nil {
			return dst, err
		}

		length := int(u16)
		dst = allocAtLeast(dst, length)
		if err = readAtLeastBytesFromReader(reader, length, dst[len(dst)-length:]); err != nil {
			return dst, err
		}
		typ = string(dst[len(dst)-length:])
		if err = o.addTyperefs(typ); err != nil {
			return dst, err
		}

	case 0x54, 0x75:
		// nolint
		reader.ReadByte()

		u32, err = readUint32FromReader(reader)
		if err != nil {
			return dst, err
		}

		typ, err = o.getTyperefs(int(u32))
		if err != nil {
			return dst, err
		}
		dst = append(dst, typ...)

	default:
	}

	return dst, nil
}
