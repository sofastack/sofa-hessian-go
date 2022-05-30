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
)

// EncodeInt64ToHessian4V2 encodes int64 to dst.
//
// long ::= L b7 b6 b5 b4 b3 b2 b1 b0
// 0000  ::= [xd8-xef]
// 0000  ::= [xf0-xff] b0
// 0000  ::= [x38-x3f] b1 b0
// 0000  ::= x4c b3 b2 b1 b0
// A 64-bit signed integer. An long is represented by the octet x4c ('L' ) followed
// by the 8-bytes of the integer in big-endian order.
func EncodeInt64ToHessian4V2(o *EncodeContext, dst []byte, n int64) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeint64")
		defer o.tracer.OnTraceStop("encodeint64")
	}

	switch {
	case LONGDIRECTMIN <= n && n <= LONGDIRECTMAX:
		dst = append(dst, uint8(n)+0xe0)

	case LONGBYTEMIN <= n && n <= LONGBYTEMAX:
		dst = append(dst, uint8(n>>8)+LONGBYTEZERO, uint8(n&0xFF))

	case LONGSHORTMIN <= n && n <= LONGSHORTMAX:
		dst = append(dst, uint8(n>>16)+LONGSHORTZERO, 0, 0)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(n&0xFFFF))

	case -0x80000000 <= n && n <= 0x7fffffff:
		dst = append(dst, "Y0000"...)
		binary.BigEndian.PutUint32(dst[len(dst)-4:], uint32(n))

	default:
		dst = append(dst, "L00000000"...)
		binary.BigEndian.PutUint64(dst[len(dst)-8:], uint64(n))
	}

	return dst, nil
}

// EncodeInt32ToHessian4V2 encodes int32 to dst.
//
// int ::= 'I' b3 b2 b1 b0
// 0000 ::= [x80-xbf]
// 0000 ::= [xc0-xcf] b0
// 0000 ::= [xd0-xd7] b1 b0
//
// A 32-bit signed integer. An integer is represented by the octet x49 ('I') followed
// by the 4 octets of the integer in big-endian order.
// value = (b3 << 24) + (b2 << 16) + (b1 << 8) + b0;
func EncodeInt32ToHessian4V2(o *EncodeContext, dst []byte, n int32) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeint32")
		defer o.tracer.OnTraceStop("encodeint32")
	}

	switch {
	case INTDIRECTMIN <= n && n <= INTDIRECTMAX:
		dst = append(dst, uint8(n+INTZERO))

	case INTBYTEMIN <= n && n <= INTBYTEMAX:
		dst = append(dst, uint8(n>>8)+INTBYTEZERO, uint8(n&0xFF))

	case INTSHORTMIN <= n && n <= INTSHORTMAX:
		dst = append(dst, uint8(n>>16)+INTSHORTZERO, 0, 0)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(n&0xFFFF))

	default:
		dst = append(dst, "I0000"...)
		binary.BigEndian.PutUint32(dst[len(dst)-4:], uint32(n))
	}
	return dst, nil
}

// EncodeInt64ToHessian3V2 encodes int64 to dst with hessian3 v2 protocol.
func EncodeInt64ToHessian3V2(o *EncodeContext, dst []byte, n int64) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeint64")
		defer o.tracer.OnTraceStop("encodeint64")
	}

	switch {
	case LONGDIRECTMIN <= n && n <= LONGDIRECTMAX:
		dst = append(dst, uint8(n)+0xe0)

	case LONGBYTEMIN <= n && n <= LONGBYTEMAX:
		dst = append(dst, uint8(n>>8)+LONGBYTEZERO, uint8(n&0xFF))

	case LONGSHORTMIN <= n && n <= LONGSHORTMAX:
		dst = append(dst, uint8(n>>16)+LONGSHORTZERO, 0, 0)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(n&0xFFFF))

	case -0x80000000 <= n && n <= 0x7fffffff:
		dst = append(dst, "w0000"...)
		binary.BigEndian.PutUint32(dst[len(dst)-4:], uint32(n))

	default:
		dst = append(dst, "L00000000"...)
		binary.BigEndian.PutUint64(dst[len(dst)-8:], uint64(n))
	}

	return dst, nil
}

// EncodeInt64ToHessianV1 encodes int64 to dst with hessian3 v2 protocol.
func EncodeInt64ToHessianV1(o *EncodeContext, dst []byte, n int64) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeint64")
		defer o.tracer.OnTraceStop("encodeint64")
	}
	dst = append(dst, "L00000000"...)
	binary.BigEndian.PutUint64(dst[len(dst)-8:], uint64(n))

	return dst, nil
}

// EncodeInt32ToHessian3V2 encodes int32 to dst with hessian3 v2 protocol.
func EncodeInt32ToHessian3V2(o *EncodeContext, dst []byte, n int32) ([]byte, error) {
	return EncodeInt32ToHessian4V2(o, dst, n)
}

// EncodeInt32ToHessianV1 encodes int32 to dst with hessian3 v2 protocol.
func EncodeInt32ToHessianV1(o *EncodeContext, dst []byte, n int32) ([]byte, error) {
	dst = append(dst, "I0000"...)
	binary.BigEndian.PutUint32(dst[len(dst)-4:], uint32(n))
	return dst, nil
}

// EncodeInt32RefToHessianV1 encodes int32 to dst with hessian3 v2 protocol.
func EncodeInt32RefToHessianV1(o *EncodeContext, dst []byte, n int32) ([]byte, error) {
	dst = append(dst, "0000"...)
	binary.BigEndian.PutUint32(dst[len(dst)-4:], uint32(n))
	return dst, nil
}
