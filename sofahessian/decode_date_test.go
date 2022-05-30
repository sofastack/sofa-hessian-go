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
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeDate(t *testing.T) {
	o := &DecodeContext{
		tracer: NewDummyTracer(),
	}
	t.Run("it should read 09:51:31 May 8, 1998 UTC", func(t *testing.T) {
		br := bufio.NewReader(bytes.NewReader(
			readFile(t, "testdata/date/894621091000.bin"),
		))
		dt, err := DecodeDateHessian4V2(o, br)
		require.Nil(t, err)
		require.Equal(t, "1998-05-08 09:51:31 +0000 UTC", dt.UTC().String())
	})
}

func TestHessian1DecodeDate(t *testing.T) {
	o := &DecodeContext{
		tracer: NewDummyTracer(),
	}
	t.Run("it should read 2010-10-10 00:00:00 +0000 UTC", func(t *testing.T) {
		br := []byte{0x64, 0x00, 0x00, 0x01, 0x2b, 0x93, 0x6f, 0xd0, 0x00}
		dt, err := DecodeDateHessianV1(o, bufio.NewReader(
			bytes.NewReader(br),
		))
		require.Nil(t, err)
		require.Equal(t, "2010-10-10 00:00:00 +0000 UTC", dt.UTC().String())
	})
}
