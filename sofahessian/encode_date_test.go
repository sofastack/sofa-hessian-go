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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEncodeDate(t *testing.T) {
	o := &EncodeContext{}
	for _, mt := range []struct {
		Input  string
		Output []byte
	}{
		{
			"09:51:31 May 8, 1998 UTC",
			[]byte{0x4a, 0x00, 0x00, 0x00, 0xd0, 0x4b, 0x92, 0x84, 0xb8},
		},
	} {
		layout := "15:04:05 Jan 2, 2006 MST"
		x, err := time.Parse(layout, mt.Input)
		require.Nil(t, err)
		require.Equal(t, mt.Input, x.Format(layout))
		dst, err := EncodeDateToHessian4V2(o, nil, x)
		require.Nil(t, err)
		require.Equal(t, mt.Output, dst, mt.Input)
	}
}
