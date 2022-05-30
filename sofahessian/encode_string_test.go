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
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeString(t *testing.T) {
	o := &EncodeContext{
		tracer: NewDummyTracer(),
	}
	for _, mt := range []struct {
		Input  string
		Output []byte
	}{
		{
			"",
			[]byte{0x00},
		},
		{
			"hello",
			[]byte{0x05, 'h', 'e', 'l', 'l', 'o'},
		},
		{
			"\u00c3",
			[]byte{0x01, 0xc3, 0x83},
		},
	} {
		dst, err := EncodeToHessian4V2(o, nil, mt.Input)
		require.Nil(t, err)
		require.Equal(t, mt.Output, dst)
	}

	for _, n := range []int{
		32767,
		32768,
		32769,
		65534,
		65535,
		65536,
		65537,
	} {
		s := fmt.Sprintf("large_string_%d.bin", n)
		fn := filepath.Join("testdata", "string", s)
		data, err := ioutil.ReadFile(fn)
		require.Nil(t, err)
		expect := strings.Repeat("A", n)
		dst, err := EncodeToHessian4V2(o, nil, expect)
		require.Nil(t, err)
		require.Equal(t, data, dst, fn)
	}

	for _, n := range []int{
		32767,
		32768,
		32769,
		65534,
		65535,
		65536,
		65537,
	} {
		s := fmt.Sprintf("utf8_%d.bin", n)
		fn := filepath.Join("testdata", "string", s)
		data, err := ioutil.ReadFile(fn)
		require.Nil(t, err)
		expect := string(bytes.Repeat([]byte("锋"), n))
		dst, err := EncodeToHessian4V2(o, nil, expect)
		require.Nil(t, err)
		require.Equal(t, data, dst, fn)
	}
}
