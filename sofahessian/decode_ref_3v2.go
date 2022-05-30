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

import "bufio"

func DecodeRefHessian3V2(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decoderef")
		defer o.tracer.OnTraceStop("decoderef")
	}

	var i interface{}
	err := DecodeRefToHessian3V2(o, reader, &i)
	return i, err
}

func DecodeRefToHessian3V2(o *DecodeContext, reader *bufio.Reader, obj *interface{}) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	var refid uint32
	if c1 == 0x4a {
		c1, err = reader.ReadByte()
		if err != nil {
			return err
		}
		refid = uint32(c1)

	} else if c1 == 0x4b {
		var u16 uint16
		u16, err = readUint16FromReader(reader)
		if err != nil {
			return err
		}
		refid = uint32(u16)

	} else if c1 == 0x52 {
		refid, err = readUint32FromReader(reader)
		if err != nil {
			return err
		}

	} else {
		return ErrDecodeMalformedReference
	}

	i, err := o.getObjectrefs(int(refid))
	if err != nil {
		return err
	}
	*obj = i

	return nil
}
