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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeInt64(t *testing.T) {
	o := &EncodeContext{
		tracer: NewDummyTracer(),
	}
	for _, mt := range []struct {
		Input  int64
		Output []byte
	}{
		{
			0,
			[]byte{0xE0},
		},
		{
			-8,
			[]byte{0xd8},
		},
		{
			15,
			[]byte{0xEF},
		},
	} {
		dst, err := EncodeToHessian4V2(o, nil, mt.Input)
		require.Nil(t, err)
		require.Equal(t, mt.Output, dst)
	}

	for _, n := range []int64{
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
		s := fmt.Sprintf("%d.bin", n)
		fn := filepath.Join("testdata", "long", s)
		data, err := ioutil.ReadFile(fn)
		require.Nil(t, err)
		dst, err := EncodeToHessian4V2(o, nil, n)
		require.Nil(t, err)
		require.Equal(t, data, dst, fn)
	}
}

func TestEncodeInt32(t *testing.T) {
	o := &EncodeContext{}
	for _, mt := range []struct {
		Input  int32
		Output []byte
	}{
		{
			0,
			[]byte{0x90},
		},
		{
			-16,
			[]byte{0x80},
		},
		{
			-2048,
			[]byte{0xc0, 0x00},
		},
		{
			-256,
			[]byte{0xc7, 0x00},
		},
		{
			2047,
			[]byte{0xcf, 0xff},
		},
	} {
		dst, err := EncodeToHessian4V2(o, nil, mt.Input)
		require.Nil(t, err)
		require.Equal(t, mt.Output, dst)
	}

	for _, n := range []int32{
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
		s := fmt.Sprintf("%d.bin", n)
		fn := filepath.Join("testdata", "number", s)
		data, err := ioutil.ReadFile(fn)
		require.Nil(t, err)
		dst, err := EncodeToHessian4V2(o, nil, n)
		require.Nil(t, err)
		require.Equal(t, data, dst, fn)
	}
}
