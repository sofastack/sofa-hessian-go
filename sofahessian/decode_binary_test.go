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
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeBinary(t *testing.T) {
	t.Run("should read zero length binary data", func(t *testing.T) {
		testDecodeBinary(t, []byte{0x20}, nil)
	})

	t.Run("should read short binary", func(t *testing.T) {
		o := &DecodeContext{}
		r := bytes.NewReader([]byte{0x20, 0x23, 0x23, 0x02, 0x03, 0x20})
		br := bufio.NewReader(r)

		c, err := DecodeBinaryHessian4V2(o, br)
		require.Nil(t, err)
		require.Equal(t, []byte(nil), c)

		c, err = DecodeBinaryHessian4V2(o, br)
		require.Nil(t, err)
		require.Equal(t, []byte{0x23, 2, 3}, c)

		c, err = DecodeBinaryHessian4V2(o, br)
		require.Nil(t, err)
		require.Equal(t, []byte(nil), c)
	})

	t.Run("should read max length short data", func(t *testing.T) {
		testDecodeBinary(t, bytes.Repeat([]byte{0x2f}, 16), bytes.Repeat([]byte{0x2f}, 15))

		for _, i := range []int{15, 16} {
			testDecodeBinary(t,
				readFile(t, fmt.Sprintf("testdata/bytes/%d.bin", i)),
				bytes.Repeat([]byte{0x41}, i),
			)
		}
	})

	t.Run("should read goled files", func(t *testing.T) {
		for _, i := range []int{
			65535,
			32768, 32769, 32767, 42769, 82769,
		} {
			testDecodeBinary(t,
				readFile(t, fmt.Sprintf("testdata/bytes/%d.bin", i)),
				bytes.Repeat([]byte{0x41}, i),
			)
		}
	})
}

func testDecodeBinary(t *testing.T, b []byte, dst []byte) {
	o := &DecodeContext{
		tracer: NewDummyTracer(),
	}
	r := bytes.NewReader(b)
	br := bufio.NewReader(r)

	c, err := DecodeBinaryHessian4V2(o, br)
	require.Nil(t, err)
	require.Equal(t, len(dst), len(c))
	require.Equal(t, hex.EncodeToString(dst), hex.EncodeToString(c))
}

func testDecodeBinaryV1(t *testing.T, b []byte, dst []byte) {
	o := &DecodeContext{
		tracer: NewDummyTracer(),
	}
	r := bytes.NewReader(b)
	br := bufio.NewReader(r)

	c, err := DecodeBinaryHessianV1(o, br)
	require.Nil(t, err)
	require.Equal(t, len(dst), len(c))
	require.Equal(t, hex.EncodeToString(dst), hex.EncodeToString(c))
}

func testDecodeBinary3V2(t *testing.T, b []byte, dst []byte) {
	o := &DecodeContext{
		tracer: NewDummyTracer(),
	}
	r := bytes.NewReader(b)
	br := bufio.NewReader(r)

	c, err := DecodeBinaryHessian3V2(o, br)
	require.Nil(t, err)
	require.Equal(t, len(dst), len(c))
	require.Equal(t, hex.EncodeToString(dst), hex.EncodeToString(c))
}
