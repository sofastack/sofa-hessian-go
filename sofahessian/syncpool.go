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
	"io"
	"sync"
)

var (
	bbrPool sync.Pool

	mStringInterfacePool = sync.Pool{New: func() interface{} {
		m := make(map[string]interface{}, 8)
		return &m
	}}

	hessianEncodeCtxPool = sync.Pool{
		New: func() interface{} {
			return NewEncodeContext()
		},
	}
	hessianDecodeCtxPool = sync.Pool{
		New: func() interface{} {
			return NewDecodeContext()
		},
	}
	hessianEncoderPool = sync.Pool{
		New: func() interface{} {
			return NewEncoder(NewEncodeContext())
		},
	}
	hessianDecoderPool = sync.Pool{
		New: func() interface{} {
			return NewDecoder(NewDecodeContext(), nil)
		},
	}
	jsonctxPool = sync.Pool{
		New: func() interface{} {
			return NewJSONEncodeContext()
		},
	}
)

func acquireMStringInterface() *map[string]interface{} {
	return mStringInterfacePool.Get().(*map[string]interface{})
}

func releaseMStringInterface(m *map[string]interface{}) {
	for i := range *m {
		delete(*m, i)
	}
	mStringInterfacePool.Put(m)
}

var utf8HelperPool = sync.Pool{
	New: func() interface{} {
		return &utf8Helper{}
	},
}

func acquireUTF8Helper(s string) *utf8Helper {
	return utf8HelperPool.Get().(*utf8Helper)
}

func releaseUTF8Helper(h *utf8Helper) {
	h.Reset()
	utf8HelperPool.Put(h)
}

var u8Pool = sync.Pool{
	New: func() interface{} {
		var p [8]byte
		return &p
	},
}

func acquireU8() *[8]byte {
	return u8Pool.Get().(*[8]byte)
}

func releaseU8(p *[8]byte) {
	u8Pool.Put(p)
}

// AcquireHessianDecodeContext acquires decode context from sync.pool.
func AcquireHessianDecodeContext() *DecodeContext {
	return hessianDecodeCtxPool.Get().(*DecodeContext)
}

// ReleaseHessianDecodeContext releases decode context to sync.pool.
func ReleaseHessianDecodeContext(hd *DecodeContext) {
	hd.Reset()
	hessianDecodeCtxPool.Put(hd)
}

// AcquireHessianEncodeContext acquires encode context from sync.pool.
func AcquireHessianEncodeContext() *EncodeContext {
	return hessianEncodeCtxPool.Get().(*EncodeContext)
}

// ReleaseHessianEncodeContext releases encode context to sync.pool.
func ReleaseHessianEncodeContext(hc *EncodeContext) {
	hc.Reset()
	hessianEncodeCtxPool.Put(hc)
}

// AcquireHessianDecoder acquires decoder from sync.pool.
func AcquireHessianDecoder(ctx *DecodeContext, reader io.Reader) *Decoder {
	hd, ok := hessianDecoderPool.Get().(*Decoder)
	if !ok {
		panic("failed to type casting")
	}

	hd.ResetWithContext(ctx, reader)
	return hd
}

// ReleaseHessianDecoder releases decoder to sync.pool.
func ReleaseHessianDecoder(hd *Decoder) {
	hessianDecoderPool.Put(hd)
}

// AcquireHessianEncoder acquires encoder from sync.pool.
func AcquireHessianEncoder(ctx *EncodeContext) *Encoder {
	he, ok := hessianEncoderPool.Get().(*Encoder)
	if !ok {
		panic("failed to type casting")
	}

	he.ResetWithContext(ctx)
	return he
}

// ReleaseHessianEncoder releases decoder to sync.pool.
func ReleaseHessianEncoder(he *Encoder) {
	hessianEncoderPool.Put(he)
}

// AcquireJSONContext acquires json context from sync.pool.
func AcquireJSONContext() *JSONEncodeContext {
	return jsonctxPool.Get().(*JSONEncodeContext)
}

// ReleaseJSONContext releases json context to sync.pool.
func ReleaseJSONContext(jctx *JSONEncodeContext) {
	jctx.Reset()
	jsonctxPool.Put(jctx)
}

func AcquireBytesBufioReader(b []byte) *BytesBufioReader {
	i := bbrPool.Get()
	if i == nil {
		bbr := &BytesBufioReader{
			r: bytes.NewReader(b),
		}
		bbr.br = bufio.NewReaderSize(bbr.r, 8192)
		return bbr
	}

	bbr, ok := i.(*BytesBufioReader)
	if !ok {
		panic("failed to type casting")
	}
	bbr.Reset(b)

	return bbr
}

func ReleaseBytesBufioReader(bbr *BytesBufioReader) {
	bbrPool.Put(bbr)
}
