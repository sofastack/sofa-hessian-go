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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeInt32(t *testing.T) {
	t.Run("should read int32 integers", func(t *testing.T) {
		testDecodeInt32(t, 0, []byte{0x90})
		testDecodeInt32(t, -16, []byte{0x80})
		testDecodeInt32(t, 47, []byte{0xbf})
		testDecodeInt32(t, 0, []byte{0xc8, 0x00})
		testDecodeInt32(t, -256, []byte{0xc7, 0x00})
		testDecodeInt32(t, -2048, []byte{0xc0, 0x00})
		testDecodeInt32(t, 2047, []byte{0xcf, 0xff})
		testDecodeInt32(t, 0, []byte{0xd4, 0x00, 0x00})
		testDecodeInt32(t, -262144, []byte{0xd0, 0x00, 0x00})
		testDecodeInt32(t, 262143, []byte{0xd7, 0xff, 0xff})
		testDecodeInt32(t, 0, []byte{'I', 0x00, 0x00, 0x00, 0x00})
		testDecodeInt32(t, 300, []byte{'I', 0x00, 0x00, 0x01, 0x2c})
	})

	t.Run("should read int32 integers from golden file", func(t *testing.T) {
		for _, i := range []int{
			0,
			1,
			10,
			16,
			2047,
			255,
			256,
			262143,
			262144,
			46,
			47,
			-16,
			-2048,
			-256,
			-262144,
			-262145,
		} {
			b := readFile(t, fmt.Sprintf("testdata/number/%d.bin", i))
			testDecodeInt32(t, i, b)
		}
	})

	t.Run("should read and write int32 integers", func(t *testing.T) {
		testDecodeHessian1Int32(t, -10000, []byte{0x49, 0xff, 0xff, 0xd8, 0xf0})
		testDecodeHessian1Int32(t, -1, []byte{0x49, 0xff, 0xff, 0xff, 0xff})
		testDecodeHessian1Int32(t, 0, []byte{0x49, 0x00, 0x00, 0x00, 0x00})
		testDecodeHessian1Int32(t, 100000, []byte{0x49, 0x00, 0x01, 0x86, 0xa0})
		testDecodeHessian1Int32(t, 2147483647, []byte{0x49, 0x7f, 0xff, 0xff, 0xff})
	})
}

func TestDecodeInt64(t *testing.T) {
	t.Run("should read compact long", func(t *testing.T) {
		testDecodeInt64(t, 0, []byte{0xe0})
		testDecodeInt64(t, -8, []byte{0xd8})
		testDecodeInt64(t, 15, []byte{0xef})
		testDecodeInt64(t, 0, []byte{0xf8, 0x00})
		testDecodeInt64(t, -2048, []byte{0xf0, 0x00})
		testDecodeInt64(t, -256, []byte{0xf7, 0x00})
		testDecodeInt64(t, 2047, []byte{0xff, 0xff})
		testDecodeInt64(t, 0, []byte{0x3c, 0x00, 0x00})
		testDecodeInt64(t, -262144, []byte{0x38, 0x00, 0x00})
		testDecodeInt64(t, 262143, []byte{0x3f, 0xff, 0xff})
		testDecodeInt64(t, 0, []byte{0x59, 0x00, 0x00, 0x00, 0x00})
		testDecodeInt64(t, 300, []byte{0x59, 0x00, 0x00, 0x01, 0x2c})
		testDecodeInt64(t, 2147483647, []byte{0x59, 0x7f, 0xff, 0xff, 0xff})
		testDecodeInt64(t, -2147483648, []byte{0x59, 0x80, 0x00, 0x00, 0x00})
		testDecodeInt64(t, 0, []byte{0x4c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		testDecodeInt64(t, 300, []byte{0x4c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x2c})
		testDecodeInt64(t, 2147483647, []byte{0x4c, 0x00, 0x00, 0x00, 0x00, 0x7f, 0xff, 0xff, 0xff})
	})

	t.Run("should read normal long", func(t *testing.T) {
		testDecodeInt64(t, 2147483648, []byte{0x4c, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00})
	})

	t.Run("should read int64 integers from golden file", func(t *testing.T) {
		for _, i := range []int{
			-7,
			-8,
			-9,
			-2048,
			-2049,
			-262144,
			-2147483647,
			-2147483648,
			0, 14, 15, 16, 255, 2047, 2048, 262143, 2147483646, 2147483647, 2147483648,
		} {
			b := readFile(t, fmt.Sprintf("testdata/long/%d.bin", i))
			testDecodeInt64(t, i, b)
		}
	})

	t.Run("should read compact long", func(t *testing.T) {
		testDecodeHessian1Int64(t, -9223372036854775808, []byte{0x4c, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		testDecodeInt64(t, -10000, []byte{0x4c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xd8, 0xf0})
		testDecodeInt64(t, -1, []byte{0x4c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
		testDecodeInt64(t, 0, []byte{0x4c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		testDecodeInt64(t, 10000, []byte{0x4c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x27, 0x10})
		testDecodeInt64(t, 9007199254740991, []byte{0x4c, 0x00, 0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
		testDecodeInt64(t, 9007199254740992, []byte{0x4c, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		testDecodeInt64(t, 9007199254740993, []byte{0x4c, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})
		testDecodeInt64(t, 9223372036854775807, []byte{0x4c, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	})
}

func testDecodeInt32(t *testing.T, i int, dst []byte) {
	o := &DecodeContext{}
	i32, err := DecodeInt32Hessian4V2(o, bufio.NewReader(bytes.NewReader(dst)))
	require.Nil(t, err)
	require.Equal(t, int32(i), i32)
}

func testDecodeInt64(t *testing.T, i int, dst []byte) {
	o := &DecodeContext{}
	i64, err := DecodeInt64Hessian4V2(o, bufio.NewReader(bytes.NewReader(dst)))
	require.Nil(t, err)
	require.Equal(t, int64(i), i64)
}

func testDecodeHessian1Int32(t *testing.T, i int, dst []byte) {
	o := &DecodeContext{}
	i32, err := DecodeInt32HessianV1(o, bufio.NewReader(bytes.NewReader(dst)))
	require.Nil(t, err)
	require.Equal(t, int32(i), i32)
}

func testDecodeHessian1Int64(t *testing.T, i int, dst []byte) {
	o := &DecodeContext{}
	i64, err := DecodeInt64HessianV1(o, bufio.NewReader(bytes.NewReader(dst)))
	require.Nil(t, err)
	require.Equal(t, int64(i), i64)
}
