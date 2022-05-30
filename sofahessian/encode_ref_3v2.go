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

import "encoding/binary"

// EncodeRefHessian3V2 encodes refid to dst.
func EncodeRefHessian3V2(o *EncodeContext, dst []byte, refid uint32) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encoderef")
		defer o.tracer.OnTraceStop("encoderef")
	}

	if refid < 0x100 {
		dst = append(dst, 'J')
		dst = append(dst, uint8(refid))

	} else if refid < 0x10000 {
		dst = append(dst, "K00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(refid))

	} else {
		dst = append(dst, "R0000"...)
		binary.BigEndian.PutUint32(dst[len(dst)-4:], refid)
	}

	return dst, nil
}

func encodeTyperefToHessian3V2(o *EncodeContext, dst []byte, typ string) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodetyperef")
		defer o.tracer.OnTraceStop("encodetyperef")
	}

	if typ == "" {
		return dst, nil
	}

	refid, ok, err := o.getTyperefs(typ)
	if err != nil {
		return dst, err
	}

	if !ok {
		if err = o.addTyperefs(typ); err != nil {
			return dst, err
		}
		dst = append(dst, "t00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(len(typ)))
		dst = append(dst, typ...)
		return dst, nil
	}

	return EncodeInt32ToHessian3V2(o, dst, int32(refid))
}

func encodeObjectrefToHessian3V2(o *EncodeContext, dst []byte, obj interface{}) ([]byte, int, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeoobjectref")
		defer o.tracer.OnTraceStop("encodeoobjectref")
	}

	if o.disableObjectrefs {
		return dst, -1, nil
	}

	ref, err := o.getObjectrefs(obj)
	if err != nil {
		return dst, -1, err
	}

	if ref >= 0 {
		dst, err = EncodeRefHessian3V2(o, dst, uint32(ref))
		return dst, ref, err
	}
	_, err = o.addObjectrefs(obj)

	return dst, -1, err
}
