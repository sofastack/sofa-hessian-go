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
	"testing"

	"github.com/stretchr/testify/require"
)

// nolint
func TestDecode(t *testing.T) {
	t.Run("should read decode complex data", func(t *testing.T) {
		d := readFile(t, "testdata/complex/hessian3.txt")
		c, err := hex.DecodeString(string(d))
		require.Nil(t, err)
		decoder := NewDecoder(NewDecodeContext().SetVersion(Hessian3xV2), bytes.NewReader(c))
		_, err = decoder.Decode()
		require.Nil(t, err)

		_, err = decoder.Decode()
		require.Nil(t, err)

		s, err := decoder.DecodeString()
		require.Nil(t, err)
		require.Equal(t, 72420, len(s))
	})

	// nolint
	t.Run("should decode tr", func(t *testing.T) {
		x, err := hex.DecodeString(`4d74002a636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e526571756573745300036374784d740039636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e526571756573742452657175657374436f6e7465787453000269644c000000000041fb9653000674686973243052000000007a7a`)
		require.Nil(t, err)
		decoder := NewDecoder(NewDecodeContext().SetVersion(HessianV1), bytes.NewReader(x))
		obj, err := decoder.Decode()
		require.Nil(t, err)
		_ = obj

		x = []byte{79, 187, 99, 111, 109, 46, 116, 97, 111, 98, 97, 111, 46, 114, 101, 109, 111, 116, 105, 110, 103, 46, 105, 109, 112, 108, 46, 67, 111, 110, 110, 101, 99, 116, 105, 111, 110, 82, 101, 115, 112, 111, 110, 115, 101, 149, 4, 104, 111, 115, 116, 6, 114, 101, 115, 117, 108, 116, 8, 101, 114, 114, 111, 114, 77, 115, 103, 10, 101, 114, 114, 111, 114, 83, 116, 97, 99, 107, 3, 99, 116, 120, 111, 144, 78, 144, 78, 78, 79, 200, 59, 99, 111, 109, 46, 116, 97, 111, 98, 97, 111, 46, 114, 101, 109, 111, 116, 105, 110, 103, 46, 105, 109, 112, 108, 46, 67, 111, 110, 110, 101, 99, 116, 105, 111, 110, 82, 101, 115, 112, 111, 110, 115, 101, 36, 82, 101, 115, 112, 111, 110, 115, 101, 67, 111, 110, 116, 101, 120, 116, 146, 2, 105, 100, 6, 116, 104, 105, 115, 36, 48, 111, 145, 225, 74, 0}
		decoder = NewDecoder(NewDecodeContext().SetVersion(Hessian3xV2), bytes.NewReader(x))
		obj, err = decoder.Decode()
		require.Nil(t, err)
		_ = obj

		x, err = hex.DecodeString(`4fba636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e5265717565737491036374786f904fc839636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e526571756573742452657175657374436f6e7465787492026964067468697324306f91fb564a00`)
		require.Nil(t, err)
		decoder = NewDecoder(NewDecodeContext().SetVersion(Hessian3xV2), bytes.NewReader(x))
		obj, err = decoder.Decode()
		require.Nil(t, err)
		_ = obj

		x, err = hex.DecodeString(`4fba636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e5265717565737491036374786f904fc839636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e526571756573742452657175657374436f6e7465787492026964067468697324306f91fbe54a00`)
		require.Nil(t, err)
		decoder = NewDecoder(NewDecodeContext().SetVersion(Hessian3xV2), bytes.NewReader(x))
		obj, err = decoder.Decode()
		require.Nil(t, err)
		_ = obj
	})

	t.Run("should read java exception", func(t *testing.T) {
		yo := &JavaObject{
			class: "java.lang.StackTraceElement",
			names: []string{"declaringClass", "methodName", "fileName", "lineNumber"},
			values: []interface{}{
				"hessian.Main",
				"main",
				"Main.java",
				int32(1283),
			},
		}

		yoo := &JavaObject{
			class: "java.io.IOException",
			names: []string{"detailMessage", "cause", "stackTrace"},
			values: []interface{}{
				"this is a java IOException instance",
				nil,
				&JavaList{
					class: "[java.lang.StackTraceElement",
					value: []interface{}{
						yo,
					},
				},
			},
		}
		yoo.values[1] = yoo

		jo := &JavaObject{
			class: "java.lang.reflect.UndeclaredThrowableException",
			names: []string{"undeclaredThrowable", "detailMessage", "cause", "stackTrace"},
			values: []interface{}{
				yoo,
				nil,
				nil,
				&JavaList{
					class: "[java.lang.StackTraceElement",
					value: []interface{}{
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{"declaringClass", "methodName", "fileName", "lineNumber"},
							values: []interface{}{
								"hessian.Main",
								"main",
								"Main.java",
								int32(1283),
							},
						},
					},
				},
			},
		}

		testDecode(t,
			readFile(t, "testdata/exception/UndeclaredThrowableException.bin"),
			jo,
		)
	})

	t.Run("should decode circle", func(t *testing.T) {
		joo := &JavaObject{
			class: "hessian.ConnectionRequest$RequestContext",
			names: []string{"id", "this$0"},
			values: []interface{}{
				int32(101),
				nil,
			},
		}
		jo := &JavaObject{
			class: "hessian.ConnectionRequest",
			names: []string{"ctx"},
			values: []interface{}{
				joo,
			},
		}

		joo.values[1] = jo

		testDecode(t,
			readFile(t, "testdata/object/ConnectionRequest.bin"),
			jo,
		)
		// testEncodeToJSON(t, jo, "abcd")
	})

	t.Run("should decode circle for concrete type", func(t *testing.T) {
		req := &HessianConnectionRequest{}
		ctx := HessianConnectionRequestContext{
			Id:   101,
			This: &HessianConnectionRequest{nil},
		}
		ctxp := &ctx
		ctxpp := &ctxp
		ctxppp := &ctxpp
		ctxpppp := &ctxppp
		req.Ctx = &ctxpppp
		ctx.This = req

		testDecodeWithConcerteType(t,
			readFile(t, "testdata/object/ConnectionRequest.bin"),
			req,
		)
	})
}

func testDecode(t *testing.T, b []byte, x interface{}) {
	o := NewDecodeContext()
	y, err := DecodeHessian4V2(o, bufio.NewReader(bytes.NewReader(b)))
	require.Nil(t, err)
	require.Equal(t, x, y)
}

func testDecodeWithConcerteType(t *testing.T, b []byte, x interface{}) {
	hr := HessianConnectionRequest{}
	cr := NewClassRegistry()
	cr.RegisterJavaClass(hr)
	o := NewDecodeContext().SetClassRegistry(cr)
	y, err := DecodeHessian4V2(o, bufio.NewReader(bytes.NewReader(b)))
	require.Nil(t, err)
	require.Equal(t, x, y)
}
