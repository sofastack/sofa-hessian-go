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
)

func BenchmarkEncodeAndDecodeJSON(b *testing.B) {
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

	c := complex{
		Bool:      Bool,
		Boolp:     Boolp,
		Boolpp:    Boolpp,
		Int32:     Int32,
		Int32p:    Int32p,
		Int32pp:   Int32pp,
		Int64:     Int64,
		Int64p:    Int64p,
		Int64pp:   Int64pp,
		Float64:   Float64,
		Float64p:  Float64p,
		Float64pp: Float64pp,
		String:    String,
		Stringp:   Stringp,
		Stringpp:  Stringpp,
		Time:      Time,
		Timep:     Timep,
		Timepp:    Timepp,
		Binary:    Binary,
		Binaryp:   Binaryp,
		Binarypp:  Binarypp,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dst, err := json.Marshal(c)
		if err != nil {
			b.Fatal(err)
		}

		b.SetBytes(int64(len(dst)))

		err = json.Unmarshal(dst, &c)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeAndDecodeHessian4V2(b *testing.B) {
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

	c := complex{
		Bool:      Bool,
		Boolp:     Boolp,
		Boolpp:    Boolpp,
		Int32:     Int32,
		Int32p:    Int32p,
		Int32pp:   Int32pp,
		Int64:     Int64,
		Int64p:    Int64p,
		Int64pp:   Int64pp,
		Float64:   Float64,
		Float64p:  Float64p,
		Float64pp: Float64pp,
		String:    String,
		Stringp:   Stringp,
		Stringpp:  Stringpp,
		Time:      Time,
		Timep:     Timep,
		Timepp:    Timepp,
		Binary:    Binary,
		Binaryp:   Binaryp,
		Binarypp:  Binarypp,
	}
	dst := make([]byte, 0, 1024)
	r := bytes.NewReader(nil)
	br := bufio.NewReader(r)

	ectx := NewEncodeContext()
	dctx := NewDecodeContext()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ectx.objectrefs.Reset()
		dst, err := EncodeToHessian4V2(ectx, dst[:0], c)
		if err != nil {
			b.Fatal(err)
		}

		b.SetBytes(int64(len(dst)))
		r.Reset(dst)
		br.Reset(r)

		dctx.objectrefs.Reset()
		_, err = DecodeHessian4V2(dctx, br)
		if err != nil {
			b.Fatal(err)
		}
	}
}
