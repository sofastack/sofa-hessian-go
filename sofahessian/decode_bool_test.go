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

func TestDecodeBool(t *testing.T) {
	testDecodeBool(t, []byte{'T'}, true, false)
	testDecodeBool(t, []byte{'F'}, false, false)
	testDecodeBool(t, []byte{'C'}, false, true)
}

func testDecodeBool(t *testing.T, buf []byte, b bool, haserr bool) {
	rv, err := DecodeBoolHessian4V2(&DecodeContext{
		tracer: NewDummyTracer(),
	}, bufio.NewReader(bytes.NewReader(buf)))
	if haserr {
		require.NotNil(t, err)
		return
	}
	require.Nil(t, err)
	require.Equal(t, b, rv)

	testDecodeBoolV1(t, buf, b, haserr)
	testDecodeBool3V2(t, buf, b, haserr)
}

func testDecodeBoolV1(t *testing.T, buf []byte, b bool, haserr bool) {
	rv, err := DecodeBoolHessianV1(&DecodeContext{
		tracer: NewDummyTracer(),
	}, bufio.NewReader(bytes.NewReader(buf)))
	if haserr {
		require.NotNil(t, err)
		return
	}
	require.Nil(t, err)
	require.Equal(t, b, rv)
}

func testDecodeBool3V2(t *testing.T, buf []byte, b bool, haserr bool) {
	rv, err := DecodeBoolHessian3V2(&DecodeContext{
		tracer: NewDummyTracer(),
	}, bufio.NewReader(bytes.NewReader(buf)))
	if haserr {
		require.NotNil(t, err)
		return
	}
	require.Nil(t, err)
	require.Equal(t, b, rv)
}
