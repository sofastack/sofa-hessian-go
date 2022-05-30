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

func DecodeInt32Hessian4V2(o *DecodeContext, reader *bufio.Reader) (int32, error) {
	var i int32
	err := DecodeInt32ToHessian4V2(o, reader, &i)
	return i, err
}

func DecodeInt32ToHessian4V2(o *DecodeContext, reader *bufio.Reader, i *int32) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodeint32")
		defer o.tracer.OnTraceStop("decodeint32")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	if c1 >= 0x80 && c1 <= 0xbf {
		*i = int32(c1) - 0x90
		return nil
	}

	if c1 >= 0xc0 && c1 <= 0xcf {
		c2, err := reader.ReadByte()
		if err != nil {
			return err
		}
		*i = ((int32(c1) - 0xc8) << 8) + int32(c2)
		return nil
	}

	if c1 >= 0xd0 && c1 <= 0xd7 {
		u16, err := readUint16FromReader(reader)
		if err != nil {
			return err
		}

		*i = (int32(c1)-0xd4)<<16 + int32(u16)
		return nil
	}

	if c1 == 0x49 {
		u32, err := readUint32FromReader(reader)
		if err != nil {
			return err
		}
		*i = int32(u32)
		return nil
	}

	return ErrDecodeCannotDecodeInt32
}

func DecodeInt64Hessian4V2(o *DecodeContext, reader *bufio.Reader) (int64, error) {
	var i int64
	err := DecodeInt64ToHessian4V2(o, reader, &i)
	return i, err
}

func DecodeInt64ToHessian4V2(o *DecodeContext, reader *bufio.Reader, i *int64) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodeint64")
		defer o.tracer.OnTraceStop("decodeint64")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	if c1 >= 0xd8 && c1 <= 0xef {
		*i = int64(c1) - 0xe0
		return nil
	}

	if c1 >= 0xf0 && c1 <= 0xff {
		c2, err := reader.ReadByte()
		if err != nil {
			return err
		}

		*i = (int64(c1)-0xf8)<<8 + int64(c2)
		return nil
	}

	if c1 >= 0x38 && c1 <= 0x3f {
		u16, err := readUint16FromReader(reader)
		if err != nil {
			return err
		}

		*i = ((int64(c1) - 0x3c) << 16) + int64(u16)
		return nil
	}

	if c1 == 0x59 {
		i32, err := readInt32FromReader(reader)
		if err != nil {
			return err
		}
		*i = int64(i32)
		return nil
	}

	if c1 == 0x4c {
		u64, err := readUint64FromReader(reader)
		if err != nil {
			return err
		}
		*i = int64(u64)
		return nil
	}

	return ErrDecodeCannotDecodeInt64
}
