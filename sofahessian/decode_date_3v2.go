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
	"time"
)

func DecodeDateHessian3V2(o *DecodeContext, reader *bufio.Reader) (time.Time, error) {
	var t time.Time
	err := DecodeDateToHessian3V2(o, reader, &t)
	return t, err
}

func DecodeDateToHessian3V2(o *DecodeContext, reader *bufio.Reader, t *time.Time) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	if c1 != 'd' {
		return ErrDecodeMalformedDate
	}

	u64, err := readUint64FromReader(reader)
	if err != nil {
		return err
	}

	if u64/1000/3600/24/365 >= 2262-1970 { // try to save the time when year after 2262
		var mt time.Time
		*t = mt.Round(time.Millisecond * time.Duration(u64))
	} else {
		*t = time.Unix(int64(u64/1000), int64(u64)%1000*10e5)
	}
	return nil
}
