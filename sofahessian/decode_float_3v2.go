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
	"math"
)

func DecodeFloat64Hessian3V2(o *DecodeContext, reader *bufio.Reader) (float64, error) {
	var i float64
	err := DecodeFloat64ToHessian3V2(o, reader, &i)
	return i, err
}

func DecodeFloat64ToHessian3V2(o *DecodeContext, reader *bufio.Reader, i *float64) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	switch c1 {
	case 0x44:
		u64, err := readUint64FromReader(reader)
		if err != nil {
			return err
		}
		*i = math.Float64frombits(u64)
		return nil

	case 0x67:
		*i = 0.0

	case 0x68:
		*i = 1.0

	case 0x69:
		u8, err := reader.ReadByte()
		if err != nil {
			return err
		}
		*i = float64(u8)

	case 0x6a:
		u16, err := readUint16FromReader(reader)
		if err != nil {
			return err
		}
		*i = float64(u16)

	case 0x6b:
		u32, err := readUint32FromReader(reader)
		if err != nil {
			return err
		}
		*i = float64(math.Float32frombits(u32))
		return nil

	default:
		return ErrDecodeMalformedDouble
	}

	return nil
}
