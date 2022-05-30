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

func DecodeInt32HessianV1(o *DecodeContext, reader *bufio.Reader) (int32, error) {
	var i int32
	err := DecodeInt32ToHessianV1(o, reader, &i)
	return i, err
}

func DecodeInt64HessianV1(o *DecodeContext, reader *bufio.Reader) (int64, error) {
	var i int64
	err := DecodeInt64ToHessianV1(o, reader, &i)
	return i, err
}

func DecodeInt32ToHessianV1(o *DecodeContext, reader *bufio.Reader, i *int32) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if c1 == 'I' {
		ix, err := readInt32FromReader(reader)
		*i = ix
		return err
	}

	return ErrDecodeCannotDecodeInt32
}

func DecodeInt64ToHessianV1(o *DecodeContext, reader *bufio.Reader, i *int64) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if c1 == 'L' {
		ix, err := readUint64FromReader(reader)
		if err != nil {
			return err
		}
		*i = int64(ix)
		return err
	}

	return ErrDecodeCannotDecodeInt32
}
