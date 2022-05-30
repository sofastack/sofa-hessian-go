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

// Package sofahessian implements the hessian 2.0 serialization protocol with golang
// in http://hessian.caucho.com/doc/hessian-serialization.html
//
// hessian v2 grammar
// top        ::= value
//
//            # 8-bit binary data split into 64k chunks
// binary     ::= x41 b1 b0 <binary-data> binary # non-final chunk
//            ::= 'B' b1 b0 <binary-data>        # final chunk
//            ::= [x20-x2f] <binary-data>        # binary data of
//                                                  #  length 0-15
//            ::= [x34-x37] <binary-data>        # binary data of
//                                                  #  length 0-1023
//
//            # boolean true/false
// boolean    ::= 'T'
//            ::= 'F'
//
//            # definition for an object (compact map)
// class-def  ::= 'C' string int string*
//
//            # time in UTC encoded as 64-bit long milliseconds since
//            #  epoch
// date       ::= x4a b7 b6 b5 b4 b3 b2 b1 b0
//            ::= x4b b3 b2 b1 b0       # minutes since epoch
//
//            # 64-bit IEEE double
// double     ::= 'D' b7 b6 b5 b4 b3 b2 b1 b0
//            ::= x5b                   # 0.0
//            ::= x5c                   # 1.0
//            ::= x5d b0                # byte cast to double
//                                      #  (-128.0 to 127.0)
//            ::= x5e b1 b0             # short cast to double
//            ::= x5f b3 b2 b1 b0       # 32-bit float cast to double
//
//            # 32-bit signed integer
// int        ::= 'I' b3 b2 b1 b0
//            ::= [x80-xbf]             # -x10 to x3f
//            ::= [xc0-xcf] b0          # -x800 to x7ff
//            ::= [xd0-xd7] b1 b0       # -x40000 to x3ffff
//
//            # list/vector
// list       ::= x55 type value* 'Z'   # variable-length list
// 	   ::= 'V' type int value*   # fixed-length list
//            ::= x57 value* 'Z'        # variable-length untyped list
//            ::= x58 int value*        # fixed-length untyped list
// 	   ::= [x70-77] type value*  # fixed-length typed list
// 	   ::= [x78-7f] value*       # fixed-length untyped list
//
//            # 64-bit signed long integer
// long       ::= 'L' b7 b6 b5 b4 b3 b2 b1 b0
//            ::= [xd8-xef]             # -x08 to x0f
//            ::= [xf0-xff] b0          # -x800 to x7ff
//            ::= [x38-x3f] b1 b0       # -x40000 to x3ffff
//            ::= x59 b3 b2 b1 b0       # 32-bit integer cast to long
//
//            # map/object
// map        ::= 'M' type (value value)* 'Z'  # key, value map pairs
// 	   ::= 'H' (value value)* 'Z'       # untyped key, value
//
//            # null value
// null       ::= 'N'
//
//            # Object instance
// object     ::= 'O' int value*
// 	   ::= [x60-x6f] value*
//
//            # value reference (e.g. circular trees and graphs)
// ref        ::= x51 int            # reference to nth map/list/object
//
//            # UTF-8 encoded character string split into 64k chunks
// string     ::= x52 b1 b0 <utf8-data> string  # non-final chunk
//            ::= 'S' b1 b0 <utf8-data>         # string of length
//                                              #  0-65535
//            ::= [x00-x1f] <utf8-data>         # string of length
//                                              #  0-31
//            ::= [x30-x34] <utf8-data>         # string of length
//                                              #  0-1023
//
//            # map/list types for OO languages
// type       ::= string                        # type name
//            ::= int                           # type reference
//
//            # main production
// value      ::= null
//            ::= binary
//            ::= boolean
//            ::= class-def value
//            ::= date
//            ::= double
//            ::= int
//            ::= list
//            ::= long
//            ::= map
//            ::= object
//            ::= ref
//            ::= string
package sofahessian

import "errors"

const (
	INTDIRECTMIN = -0x10    // -16
	INTDIRECTMAX = 0x2f     // 47
	INTZERO      = 0x90     // 144
	INTBYTEMIN   = -0x800   // -2048
	INTBYTEMAX   = 0x7ff    // 2047
	INTBYTEZERO  = 0xc8     // 200
	INTSHORTMIN  = -0x40000 // -262144
	INTSHORTMAX  = 0x3ffff  // 262143
	INTSHORTZERO = 0xd4     // 212

	LONGDIRECTMIN = -0x08
	LONGDIRECTMAX = 0x0f
	LONGZERO      = 0xe0
	LONGBYTEMIN   = -0x800
	LONGBYTEMAX   = 0x7ff
	LONGBYTEZERO  = 0xf8
	LONGSHORTMIN  = -0x40000
	LONGSHORTMAX  = 0x3ffff
	LONGSHORTZERO = 0x3c

	LISTDIRECT          = 0x70
	LISTDIRECTUNTYPED   = 0x78
	LISTDIRECTMAX       = 0x7
	LISTVARIABLE        = 0x55
	LISTFIXED           = 'V'
	LISTVARIABLEUNTYPED = 0x57
	LISTFIXEDUNTYPED    = 0x58
)

var (
	ErrEncodeTypeReferencesIsNil         = errors.New("hessian: type references cannot be nil")
	ErrEncodeTooManyNestedJSON           = errors.New("hessian: too many nested JSON")
	ErrEncodeNotSliceType                = errors.New("hessian: type is not slice")
	ErrEncodeNotStructType               = errors.New("hessian: type is not struct")
	ErrEncodeNotMapType                  = errors.New("hessian: type is not map")
	ErrEncodeSliceElemCannotBeInterfaced = errors.New("hessian: slice[i] cannot be interfaced")
	ErrEncodeSliceCannotBeInterfaced     = errors.New("hessian: slice cannot be interfaced")
	ErrEncodeStructCannotBeInterfaced    = errors.New("hessian: struct cannot be interfaced")
	ErrEncodeMapCannotBeInterfaced       = errors.New("hessian: map cannot be interfaced")
	ErrEncodePtrCannotBeInterfaced       = errors.New("hessian: pointer cannot be interfaced")
	ErrEncodeCannotInvalidValue          = errors.New("hessian: cannot encode invalid value")
	ErrDecodeBufferNotEnough             = errors.New("hessian: buffer not enough")
	ErrDecodeCannotDecodeInt32           = errors.New("hessian: malformed int32")
	ErrDecodeCannotDecodeInt64           = errors.New("hessian: malformed int64")
	ErrDecodeMalformedBool               = errors.New("hessian: malformed bool")
	ErrDecodeMalformedNil                = errors.New("hessian: malformed nil")
	ErrDecodeMalformedDouble             = errors.New("hessian: malformed double")
	ErrDecodeMalformedMap                = errors.New("hessian: malformed map")
	ErrDecodeMalformedUntypedMap         = errors.New("hessian: malformed untyped map")
	ErrDecodeMalformedTypedMap           = errors.New("hessian: malformed typed map")
	ErrDecodeMalformedDate               = errors.New("hessian: malformed date")
	ErrDecodeMalformedString             = errors.New("hessian: malformed string")
	ErrDecodeMalformedBinary             = errors.New("hessian: malformed binary")
	ErrDecodeMalformedObject             = errors.New("hessian: malformed object")
	ErrDecodeUnmatchedObject             = errors.New("hessian: unmatched object")
	ErrDecodeMalformedList               = errors.New("hessian: malformed list")
	ErrDecodeMalformedListEnd            = errors.New("hessian: malformed list end")
	ErrDecodeMalformedReference          = errors.New("hessian: malformed reference")
	ErrDecodeMalformedJSON               = errors.New("hessian: malformed json")
	ErrDecodeClassRefsIsNil              = errors.New("hessian: classrefs is nil")
	ErrDecodeClassRefsOverflow           = errors.New("hessian: classrefs overflow")
	ErrDecodeObjectRefsIsNil             = errors.New("hessian: objectrefs is nil")
	ErrDecodeTypeRefsIsNil               = errors.New("hessian: typerefs is nil")
	ErrDecodeTypeRefsOverflow            = errors.New("hessian: typerefs overflow")
	ErrDecodeObjectRefsOverflow          = errors.New("hessian: objectrefs overflow")
	ErrDecodeTypedMapKeyNotString        = errors.New("hessian: typedmap key is not string")
	ErrDecodeTypedMapValueNotAssign      = errors.New("hessian: typedmap value cannot set")
	ErrDecodeMaxDepthExceeded            = errors.New("hessian: encode depth exceeded")
	ErrEncodeMaxDepthExceeded            = errors.New("hessian: decode depth exceeded")
	ErrDecodeMaxListLengthExceeded       = errors.New("hessian: decode max list length exceeded")
	ErrDecodeMaxObjectFieldsExceeded     = errors.New("hessian: decode max object fields exceeded")
	ErrDecodeUnknownEncoding             = errors.New("hessian: decode unknown encoding")
	ErrDecodeObjectFieldCannotBeNull     = errors.New("hessian: decode object field cannot be null")
	ErrDecodeMapUnhashable               = errors.New("hessian: cannot set the value to the unhashable key of map")
)

type Version uint8

const (
	Hessian4xV2 Version = iota // default Hessian4x
	Hessian3xV2
	HessianV1
)
