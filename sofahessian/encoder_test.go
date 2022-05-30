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
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type ResponseContext struct {
	Id   int64       `hessian:"id"`
	This interface{} `hessian:"this$0"`
}

func (c ResponseContext) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionResponse$ResponseContext"
}

type ConnectionResponse struct {
	Ctx        *ResponseContext `hessian:"ctx"`
	Host       string           `hessian:"host"`
	Result     int32            `hessian:"result"`
	ErrorMsg   string           `hessian:"errorMsg"`
	ErrorStack string           `hessian:"errorStack"`
	FromAppKey string           `hessian:"fromAppKey"`
	ToAppKey   string           `hessian:"toAppKey"`
}

func (c ConnectionResponse) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionResponse"
}

func TestEncoder(t *testing.T) {
	t.Run("should encode complex type", func(t *testing.T) {
		Bool := false
		Boolp := &Bool
		Boolpp := &Boolp
		Int32 := int32(13234)
		Int32p := &Int32
		Int32pp := &Int32p
		Int64 := int64(16434)
		Int64p := &Int64
		Int64pp := &Int64p
		Float64 := float64(16434.2)
		Float64p := &Float64
		Float64pp := &Float64p

		String := "1.2.3😎😎😎😎😎😎😎😎😎😎😎"
		Stringp := &String
		Stringpp := &Stringp

		ms := time.Now().UnixNano() / 1000 / 1000
		Time := time.Unix(0, ms*1000*1000)
		Timep := &Time
		Timepp := &Timep

		Binary := []byte("🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻")
		Binaryp := &Binary
		Binarypp := &Binaryp

		Foo := structfoo{
			A: 132323,
		}
		Foop := &Foo
		Foopp := &Foop

		Map := map[string]string{
			"a": "b",
			"4": "5",
		}
		Mapp := &Map
		Mappp := &Mapp

		UntypedSlice := []string{"a", "1"}

		Slice := []interface{}{
			"a",
			1,
			"c",
		}
		Slicep := &Slice
		Slicepp := &Slicep

		c := complex{
			Bool:         Bool,
			Boolp:        Boolp,
			Boolpp:       Boolpp,
			Int32:        Int32,
			Int32p:       Int32p,
			Int32pp:      Int32pp,
			Int64:        Int64,
			Int64p:       Int64p,
			Int64pp:      Int64pp,
			Float64:      Float64,
			Float64p:     Float64p,
			Float64pp:    Float64pp,
			String:       String,
			Stringp:      Stringp,
			Stringpp:     Stringpp,
			Time:         Time,
			Timep:        Timep,
			Timepp:       Timepp,
			Binary:       Binary,
			Binaryp:      Binaryp,
			Binarypp:     Binarypp,
			Foo:          Foo,
			Foop:         Foop,
			Foopp:        Foopp,
			Map:          Map,
			Mapp:         Mapp,
			Mappp:        Mappp,
			Slice:        Slice,
			Slicep:       Slicep,
			Slicepp:      Slicepp,
			UntypedSlice: UntypedSlice,
		}

		testEncoderEncodeDecodeByJSONEqual(
			t, c, c,
		)
	})
}

func testEncoderEncodeDecodeByJSONEqual(t *testing.T, x, y interface{}) {
	encoder := NewEncoder(NewEncodeContext().DisableObjectrefs())
	err := encoder.Encode(x)
	require.Nil(t, err)

	cr := NewClassRegistry()
	cr.RegisterJavaClass(complex{})
	z, err := DecodeHessian4V2(NewDecodeContext().
		SetClassRegistry(cr).
		SetTracer(NewDummyTracer()),
		bufio.NewReader(
			bytes.NewReader(
				encoder.Bytes())))

	require.Nil(t, err)
	g1, err := json.MarshalIndent(z, "", "    ")
	require.Nil(t, err)
	g2, err := json.MarshalIndent(y, "", "    ")
	require.Nil(t, err)
	require.Equal(t, string(g2), string(g1))
}
