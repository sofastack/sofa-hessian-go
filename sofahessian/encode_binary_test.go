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
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeAndDecodeBinaryV1(t *testing.T) {
	for _, s := range []string{
		"1234",
		"abcd",
		"11111111111111",
		strings.Repeat("😯", 1000),
		"aaaaaa",
		"bbbbbbbb",
		"ccccccccccccccccccccccccccccccccc",
		"asdfasdfasldfj98723ljsdfo879723",
	} {
		dst, err := EncodeBinaryToHessianV1(NewEncodeContext(), nil, []byte(s))
		require.Nil(t, err)
		z, err := DecodeBinaryHessianV1(NewDecodeContext(), bufio.NewReader(bytes.NewReader(dst)))
		require.Nil(t, err)
		require.Equal(t, s, string(z))
	}
}

func TestEncodeBinary(t *testing.T) {
	o := &EncodeContext{
		tracer: NewDummyTracer(),
	}
	for _, mt := range []struct {
		Input  []byte
		Output []byte
	}{
		{
			nil,
			[]byte{0x20},
		},
		{
			[]byte{0x01, 0x02, 0x03},
			[]byte{0x23, 0x01, 0x02, 0x03},
		},
	} {
		dst, err := EncodeToHessian4V2(o, nil, mt.Input)
		require.Nil(t, err)
		require.Equal(t, mt.Output, dst)
	}

	for _, n := range []int{
		15,
		16,
		32767,
		32768,
		32769,
		42769,
		65535,
		82769,
	} {
		s := fmt.Sprintf("%d.bin", n)
		fn := filepath.Join("testdata", "bytes", s)
		data, err := ioutil.ReadFile(fn)
		require.Nil(t, err)
		expect := bytes.Repeat([]byte{0x41}, n)
		dst, err := EncodeToHessian4V2(o, nil, expect)
		require.Nil(t, err)
		require.Equal(t, data, dst, fn)
	}
}
