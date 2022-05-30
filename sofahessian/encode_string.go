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
	"encoding/binary"
	"unicode/utf8"
)

// EncodeStringToHessian4V2 encodes string to dst.
//
// # UTF-8 encoded character string split into 64k chunks
// string0000 ::= x52 b1 b0 <utf8-data> string  # non-final chunk
// 000000000000::= 'S' b1 b0 <utf8-data>00000000 # string of length
// 00000000000000000000000000000000000000000000  #  0-65535
// 000000000000::= [x00-x1f] <utf8-data>00000000 # string of length
// 00000000000000000000000000000000000000000000  #  0-31
// 000000000000::= [x30-x34] <utf8-data>00000000 # string of length
// 00000000000000000000000000000000000000000000  #  0-1023
// A 16-bit unicode character string encoded in UTF-8. Strings are encoded in chunks.
// x53 ('S') represents the final chunk and x52 ('R') represents any non-final chunk.
// Each chunk has a 16-bit unsigned integer length value.
//
// The length is the number of 16-bit characters, which may be different than the number of bytes.
//
// String chunks may not split surrogate pairs.
func EncodeStringToHessian4V2(o *EncodeContext, dst []byte, value string) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodestring")
		defer o.tracer.OnTraceStop("encodestring")
	}

	var (
		n   = utf8.RuneCountInString(value)
		err error
	)

	if n > 0x8000 {
		h := acquireUTF8Helper(value)
		n, err = h.Write(value)
		if err != nil {
			releaseUTF8Helper(h)
			return nil, err
		}
		os := h.GetOffsetSlice()
		cs := h.GetCountSlice()

		for n > 0x8000 {
			ulen := uint16(0x8000)
			offset := cs[0x8000]

			dst = append(dst, "R00"...)
			binary.BigEndian.PutUint16(dst[len(dst)-2:], ulen)
			dst = append(dst, value[:offset]...)
			value = value[offset:]
			n = n - os[offset]
		}

		releaseUTF8Helper(h)
	}

	if n < 32 {
		dst = append(dst, uint8(n))
	} else if n < 1024 {
		dst = append(dst, 48+uint8(n>>8), uint8(n&0xFF))
	} else {
		dst = append(dst, "S00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(n))
	}
	dst = append(dst, value...)

	return dst, nil
}

func EncodeStringToHessian3V2(o *EncodeContext, dst []byte, value string) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodestring")
		defer o.tracer.OnTraceStop("encodestring")
	}

	var (
		n   = utf8.RuneCountInString(value)
		err error
	)

	if n > 0x8000 {
		h := acquireUTF8Helper(value)
		n, err = h.Write(value)
		if err != nil {
			releaseUTF8Helper(h)
			return nil, err
		}
		os := h.GetOffsetSlice()
		cs := h.GetCountSlice()

		for n > 0x8000 {
			ulen := uint16(0x8000)
			offset := cs[0x8000]

			dst = append(dst, "s00"...)
			binary.BigEndian.PutUint16(dst[len(dst)-2:], ulen)
			dst = append(dst, value[:offset]...)
			value = value[offset:]
			n = n - os[offset]
		}

		releaseUTF8Helper(h)
	}

	if n < 32 {
		dst = append(dst, uint8(n))
	} else {
		dst = append(dst, "S00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(n))
	}
	dst = append(dst, value...)

	return dst, nil
}

func EncodeStringToHessianV1(o *EncodeContext, dst []byte, value string) ([]byte, error) {
	n := utf8.RuneCountInString(value)

	dst = append(dst, "S00"...)
	binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(n))
	dst = append(dst, value...)

	return dst, nil
}
