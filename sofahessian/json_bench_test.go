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
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fastjson"
)

func BenchmarkEncodeToJSON(b *testing.B) {
	b.Run("benchmark int8", func(b *testing.B) {
		benchmarkEncodeToJSON(b, int8(1))
	})

	b.Run("benchmark int16", func(b *testing.B) {
		benchmarkEncodeToJSON(b, int16(1))
	})

	b.Run("benchmark int32", func(b *testing.B) {
		benchmarkEncodeToJSON(b, int32(1))
	})

	b.Run("benchmark int64", func(b *testing.B) {
		benchmarkEncodeToJSON(b, int64(1))
	})

	b.Run("benchmark float64", func(b *testing.B) {
		benchmarkEncodeToJSON(b, float64(1.2345689))
	})

	b.Run("benchmark string", func(b *testing.B) {
		benchmarkEncodeToJSON(b, "abcdefghijk")
	})

	b.Run("benchmark true", func(b *testing.B) {
		benchmarkEncodeToJSON(b, true)
	})

	b.Run("benchmark false", func(b *testing.B) {
		benchmarkEncodeToJSON(b, false)
	})

	b.Run("benchmark null", func(b *testing.B) {
		benchmarkEncodeToJSON(b, nil)
	})

	b.Run("benchmark binary", func(b *testing.B) {
		benchmarkEncodeToJSON(b, "abcdefghijk")
	})

	b.Run("benchmark array", func(b *testing.B) {
		benchmarkEncodeToJSON(b, []interface{}{"abcdefghijk", 1, 1.2, 3.4, true})
	})

	b.Run("benchmark map", func(b *testing.B) {
		benchmarkEncodeToJSON(b, map[string]interface{}{
			"a": map[string]interface{}{
				"b": 1,
				"d": 1.2,
			},
		})
	})

	b.Run("benchmark demo car", func(b *testing.B) {
		benchmarkEncodeToJSON(b, &DemoCar{
			A:       "a",
			C:       "c",
			B:       "b",
			Model:   "model 1",
			Color:   "aquamarine",
			Mileage: 65536,
		})
	})
}

func benchmarkEncodeToJSON(b *testing.B, obj interface{}) {
	var (
		dst []byte
		err error
	)
	jctx := NewJSONEncodeContext().
		SetMaxPtrCycles(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dst, err = EncodeToJSON(jctx, dst[:0], obj)
		if err != nil {
			b.Fatal(err)
		}
		b.SetBytes(int64(len(dst)))
	}
}

func BenchmarkEncodeJSON(b *testing.B) {
	ectx := NewEncodeContext().SetTracer(NewDummyTracer())
	encoder := NewEncoder(ectx)
	jctx := NewJSONEncodeContext().SetTracer(NewDummyTracer())
	input := `{"a":1,"b":[1,2,3.0, true, false, null]}`

	var p fastjson.Parser
	v, err := p.Parse(input)
	require.Nil(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := encoder.EncodeFastJSON(jctx, v)
		if err != nil {
			b.Fatal(err)
		}
		b.SetBytes(int64(len(input)))
		encoder.b = encoder.b[:0]
	}
}
