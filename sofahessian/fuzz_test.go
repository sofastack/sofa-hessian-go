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
	"encoding/json"
	"log"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFuzz0(t *testing.T) {
	testFuzz(t, h(t, "4800700070024830005a"), false)
}

func TestFuzz1(t *testing.T) {
	testFuzz(t, []byte("000000000000000000000000000000000000000000000000000000"), false)
}

func TestFuzz2(t *testing.T) {
	testFuzz(t, h(t, "03db3030"), false)
}

func TestFuzz3(t *testing.T) {
	testFuzz(t, h(t, "7d56013089"), true)
}

func TestFuzz4(t *testing.T) {
	testFuzz(t, h(t, "7d430089"), true)
}

func TestFuzz5(t *testing.T) {
	testFuzz(t, h(t, "760130b5db002000d9"), false)
}

func TestFuzz6(t *testing.T) {
	testFuzz(t, h(t, "4a3030303030303030"), false)
}

func TestFuzz7(t *testing.T) {
	testFuzz(t, []byte("HH@HH"), true)
}

func TestFuzz8(t *testing.T) {
	testFuzz(t, []byte("H~ Q\x91\xbd\xbf\xef\x9b\xed"), true)
}

func TestFuzz9(t *testing.T) {
	testFuzz(t, h(t, "44ffff303030303030"), false)
}

func TestFuzz10(t *testing.T) {
	testFuzz(t, h(t, "4874013084bd84f330da5a"), false)
}

func TestFuzz11(t *testing.T) {
	testFuzz(t, h(t, "4300910060484d01305a485a5a"), true)
}

func TestFuzz12(t *testing.T) {
	testFuzz(t, h(t, "4300910060484d01305a485a5a"), true)
}

func TestFuzz13(t *testing.T) {
	testFuzz(t, h(t, "4800700070024830005a"), true)
}

func TestFuzz14(t *testing.T) {
	testFuzz(t, h(t, "7d007001307899ef"), true)
}

func TestFuzz15(t *testing.T) {
	testFuzz(t, h(t, "4309303030303030303030930930303030303030303001300630303030303060544d01305abd"), true)
}

func TestFuzz16(t *testing.T) {
	testFuzz(t, h(t, "48007000701a3030bf3030ff303030303030303030bd30303030303030303030dd5a"), true)
}

func TestFuzz17(t *testing.T) {
	testFuzz(t, h(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709"), true)
}

func h(t *testing.T, s string) []byte {
	d, err := hex.DecodeString(s)
	require.Nil(t, err)
	return d
}

func UnixMilli(t time.Time) int64 {
	ns := t.UnixNano()
	if ns < 0 {
		return (ns - 999999) / 1000000
	}
	return ns / 1000000
}

func testFuzz(t *testing.T, data []byte, allowerr bool) {
	dctx := NewDecodeContext().SetTracer(NewDummyTracer()).SetMaxDepth(20)
	ectx := NewEncodeContext()
	br := bufio.NewReader(bytes.NewReader(data))
	v, err := DecodeHessian4V2(dctx, br)
	if fv, ok := v.(float64); ok && math.IsNaN(fv) {
		return
	}
	if allowerr {
		if err != nil {
			return
		}
	}
	require.Nil(t, err)

	dst, err := EncodeHessian4V2(ectx, v)
	require.Nil(t, err)
	o, err := DecodeHessian4V2(dctx, bufio.NewReader(bytes.NewReader(dst)))
	if err == ErrDecodeMapUnhashable {
		return
	}
	if allowerr {
		if err != nil {
			return
		}
	}
	require.Nil(t, err)

	if mv, ok := v.(map[interface{}]interface{}); ok {
		if ov, ok := o.(map[interface{}]interface{}); ok {
			require.Equal(t, len(mv), len(ov))
			jm, err := json.Marshal(mv)
			if err != nil {
				return
			}
			om, err := json.Marshal(ov)
			require.Nil(t, err)
			require.Equal(t, jm, om)
		}
		log.Fatal("expect map[interface{}]interface{}")
	}

	require.Equal(t, v, o)
}
