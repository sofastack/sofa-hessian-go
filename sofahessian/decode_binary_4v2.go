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

func DecodeBinaryHessian4V2(o *DecodeContext, reader *bufio.Reader) ([]byte, error) {
	p := []byte(nil)
	return DecodeBinaryToHessian4V2(o, reader, p)
}

func DecodeBinaryToHessian4V2(o *DecodeContext, reader *bufio.Reader, dst []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodebinary")
		defer o.tracer.OnTraceStop("decodebinary")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return dst, err
	}

	if c1 >= 0x20 && c1 <= 0x2f {
		length := int(c1) - 0x20
		dst = allocAtLeast(dst, length)
		err = readAtLeastBytesFromReader(reader, length, dst[len(dst)-length:])
		return dst, err
	}

	var (
		length uint16
		u16    uint16
		c2     byte
	)

	for c1 == 0x41 {
		u16, err = readUint16FromReader(reader)
		if err != nil {
			return dst, err
		}

		length = u16
		dst = allocAtLeast(dst, int(length))
		err = readAtLeastBytesFromReader(reader, int(length), dst[len(dst)-int(length):])
		if err != nil {
			return dst, err
		}

		c1, err = reader.ReadByte()
		if err != nil {
			return dst, err
		}
	}

	switch {
	case c1 == 0x42:
		u16, err = readUint16FromReader(reader)
		if err != nil {
			return dst, err
		}
		length = u16
		dst = allocAtLeast(dst, int(length))
		err = readAtLeastBytesFromReader(reader, int(length), dst[len(dst)-int(length):])
		if err != nil {
			return dst, err
		}
	case c1 >= 0x20 && c1 <= 0x2f:
		length = uint16(c1) - 0x20
		dst = allocAtLeast(dst, int(length))
		err = readAtLeastBytesFromReader(reader, int(length), dst[len(dst)-int(length):])
		if err != nil {
			return dst, err
		}
	case c1 >= 0x34 && c1 <= 0x37:
		c2, err = reader.ReadByte()
		if err != nil {
			return dst, err
		}
		length = uint16(c1-0x34)<<8 + uint16(c2)
		dst = allocAtLeast(dst, int(length))
		err = readAtLeastBytesFromReader(reader, int(length), dst[len(dst)-int(length):])
		if err != nil {
			return dst, err
		}
	default:
		return dst, ErrDecodeMalformedBinary
	}

	return dst, nil
}
