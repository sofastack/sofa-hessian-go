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

func TestDecodeFloat64(t *testing.T) {
	testDecodeFloat64(t, []byte{0x5b}, 0.0)
	testDecodeFloat64(t, []byte{0x5c}, 1.0)

	testDecodeFloat64(t, []byte{0x5d, 0x00}, 0.0)
	testDecodeFloat64(t, []byte{0x5d, 0x01}, 1.0)
	testDecodeFloat64(t, []byte{0x5d, 0x7f}, 127.0)
	testDecodeFloat64(t, []byte{0x5d, 0x80}, -128.0)

	testDecodeFloat64(t, []byte{0x5e, 0x00, 0x00}, 0.0)
	testDecodeFloat64(t, []byte{0x5e, 0x00, 0x01}, 1.0)
	testDecodeFloat64(t, []byte{0x5e, 0x00, 0x80}, 128.0)
	testDecodeFloat64(t, []byte{0x5e, 0x00, 0x7f}, 127.0)
	testDecodeFloat64(t, []byte{0x5e, 0x80, 0x00}, -32768.0)
	testDecodeFloat64(t, []byte{0x5e, 0x7f, 0xff}, 32767.0)

	testDecodeFloat64(t, []byte{0x5f, 0x00, 0x00, 0x00, 0x00}, 0.0)
	testDecodeFloat64(t, []byte{0x5f, 0x00, 0x00, 0x2f, 0xda}, 12.25)

	testDecodeFloat64(t, []byte{0x44, 0x40, 0x24, 0, 0, 0, 0, 0, 0}, 10.0)

	t.Run("should read float64 from golden file", func(t *testing.T) {
		type x struct {
			fn string
			i  float64
		}
		for _, i := range []x{
			{
				"testdata/double/0.bin",
				0,
			},
			{
				"testdata/double/1.bin",
				1,
			},
			{
				"testdata/double/10.bin",
				10,
			},
			{
				"testdata/double/10.123.bin",
				10.123,
			},
			{
				"testdata/double/10.1.bin",
				10.1,
			},
			{
				"testdata/double/-128.bin",
				-128,
			},
			{
				"testdata/double/-127.9999.bin",
				-127.9999,
			},
			{
				"testdata/double/127.bin",
				127,
			},
			{
				"testdata/double/126.9989.bin",
				126.9989,
			},
			{
				"testdata/double/-32768.bin",
				-32768,
			},
			{
				"testdata/double/-32767.999.bin",
				-32767.999,
			},
			{
				"testdata/double/32767.bin",
				32767,
			},
			{
				"testdata/double/32766.99999.bin",
				32766.99999,
			},
			{
				"testdata/double/32768.bin",
				32768,
			},
			{
				"testdata/double/32768.bin",
				32768,
			},
			{
				"testdata/double/32767.99999.bin",
				32767.99999,
			},
			{
				"testdata/double/-0x800000.bin",
				-0x800000,
			},
			{
				"testdata/double/-2147483649.bin",
				-2147483649,
			},
			{
				"testdata/double/-2147483648.bin",
				-2147483648,
			},
			{
				"testdata/double/-2147483647.bin",
				-2147483647,
			},
			{
				"testdata/double/-2147483610.123.bin",
				-2147483610.123,
			},
			{
				"testdata/double/2147483648.bin",
				2147483648,
			},
			{
				"testdata/double/2147483647.bin",
				2147483647,
			},
			{
				"testdata/double/2147483646.bin",
				2147483646,
			},
			{
				"testdata/double/2147483646.456.bin",
				2147483646.456,
			},
		} {
			testDecodeFloat64(t, readFile(t, i.fn), i.i)
		}
	})
}

func testDecodeFloat64(t *testing.T, b []byte, x float64) {
	o := &DecodeContext{
		tracer: NewDummyTracer(),
	}

	y, err := DecodeFloat64Hessian4V2(o, bufio.NewReader(
		bytes.NewReader(b),
	))
	require.Nil(t, err)
	require.Equal(t, x, y)
}

func TestDecodeHessian1Float64(t *testing.T) {
	testDecodeHessian1Float64(t, []byte{0x44, 0xd4, 0xb2, 0x49, 0xad, 0x25, 0x94, 0xc3, 0x7d}, -1e100)
	testDecodeHessian1Float64(t, []byte{0x44, 0xbf, 0xf1, 0xf7, 0xce, 0xd9, 0x16, 0x87, 0x2b}, -1.123)
	testDecodeHessian1Float64(t, []byte{0x44, 0xbf, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, -1)
	testDecodeHessian1Float64(t, []byte{0x44, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0)
	testDecodeHessian1Float64(t, []byte{0x44, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 1)
	testDecodeHessian1Float64(t, []byte{0x44, 0x3f, 0xf1, 0xc6, 0xa7, 0xef, 0x9d, 0xb2, 0x2d}, 1.111)
	// testDecodeHessian1Float64(t, []byte{0x44, 0xbf, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 1e320)
}

func testDecodeHessian1Float64(t *testing.T, b []byte, x float64) {
	o := &DecodeContext{
		tracer: NewDummyTracer(),
	}

	y, err := DecodeFloat64HessianV1(o, bufio.NewReader(
		bytes.NewReader(b),
	))
	require.Nil(t, err)
	require.Equal(t, x, y)
}
