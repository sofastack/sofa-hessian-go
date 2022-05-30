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
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeAndDecodeFloat64V1(t *testing.T) {
	for _, f := range []float64{
		1.0,
		2.0,
		3.423,
		4.234,
		8.234,
		1004,
	} {
		dst, err := EncodeFloat64ToHessianV1(NewEncodeContext(), nil, f)
		require.Nil(t, err)
		y, err := DecodeFloat64HessianV1(NewDecodeContext(), bufio.NewReader(bytes.NewReader(dst)))
		require.Nil(t, err)
		require.Equal(t, f, y)
	}

	for _, f := range []float32{
		1.0,
		2.0,
		3.423,
		4.234,
		8.234,
		1004,
	} {
		dst, err := EncodeFloat64ToHessianV1(NewEncodeContext(), nil, float64(f))
		require.Nil(t, err)
		y, err := DecodeFloat64HessianV1(NewDecodeContext(), bufio.NewReader(bytes.NewReader(dst)))
		require.Nil(t, err)
		require.Equal(t, float64(f), y)
	}
}

func TestEncodeFloat644V2(t *testing.T) {
	o := &EncodeContext{
		tracer: NewDummyTracer(),
	}
	for _, mt := range []struct {
		Input  float64
		Output []byte
	}{
		{
			0,
			[]byte{0x5B},
		},
		{
			1.0,
			[]byte{0x5c},
		},
	} {
		dst, err := EncodeToHessian4V2(o, nil, mt.Input)
		require.Nil(t, err)
		require.Equal(t, mt.Output, dst)
	}

	for _, n := range []string{
		"2147483646.456",
		"0",
		"1",
		"10.1",
		"10.123",
		"10",
		"126.9989",
		"127",
		"2147483646",
		"2147483647",
		"2147483648",
		"32766.99999",
		"32767.99999",
		"32767",
		"32768",
		"-127.9999",
		"-128",
		"-2147483610.123",
		"-2147483647.0",
		"-2147483647",
		"-2147483648",
		"-2147483649",
		"-32767.999",
		"-32768",
	} {
		s := fmt.Sprintf("%s.bin", n)
		fn := filepath.Join("testdata", "double", s)
		data, err := ioutil.ReadFile(fn)
		require.Nil(t, err)
		f, err := strconv.ParseFloat(n, 64)
		require.Nil(t, err)
		dst, err := EncodeToHessian4V2(o, nil, f)
		require.Nil(t, err)
		require.Equal(t, data, dst, fmt.Sprintf("%f-%s", f, fn))
	}
}
