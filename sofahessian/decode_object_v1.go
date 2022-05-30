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
	"errors"
	"reflect"
)

func DecodeObjectToHessianV1(o *DecodeContext, reader *bufio.Reader, obj interface{}) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodeotobject")
		defer o.tracer.OnTraceStop("decodetoobject")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	var refid int32

	if c1 == 0x43 {
		err = decodeObjectDefinitionHessianV1(o, reader)
		if err != nil {
			return err
		}
		return DecodeObjectToHessianV1(o, reader, obj)

	} else if c1 == 0x4f {
		refid, err = DecodeInt32HessianV1(o, reader)
		if err != nil {
			return err
		}

	} else if c1 >= 0x60 && c1 <= 0x6f {
		refid = int32(c1) - 0x60
	} else {
		return ErrDecodeMalformedObject
	}

	cd, err := o.getClassrefs(int(refid))
	if err != nil {
		return err
	}

	name := getInterfaceName(obj)
	if name != cd.class {
		return ErrDecodeUnmatchedObject
	}

	structvalue := decAllocReflectValue(reflect.ValueOf(obj))
	if err := o.addObjectrefs(obj); err != nil {
		return err
	}
	rt := decAllocReflectType(reflect.TypeOf(obj))

	for i := range cd.fields {
		field := cd.fields[i]
		fi, ok := lookupReflectField(rt, field)
		if !ok {
			if o.disallowMissingField {
				return errors.New("hessian: malformed class field (not found) " + field)
			}
			_, err := DecodeHessianV1(o, reader)
			if err != nil {
				return err
			}

		} else {

			key := structvalue.Field(fi)
			if !key.CanSet() {
				return errors.New("hessian: malformed class field (unassignable) " + field)
			}

			value, err := DecodeHessianV1(o, reader)
			if err != nil {
				return err
			}

			err = safeSetValueByReflect(key, value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DecodeObjectHessianV1(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodeobject")
		defer o.tracer.OnTraceStop("decodeobject")
	}

	c, err := reader.Peek(2)
	if err != nil {
		return nil, err
	}
	if c[1] == 't' {
		return DecodeObjectHessianV1WithType(o, reader)
	}
	return DecodeMapHessianV1(o, reader)
}

func DecodeObjectHessianV1WithType(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	// read t
	_, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	typ, err := DecodeStringHessianV1(o, reader)
	if err != nil {
		return nil, err
	}
	cd := ClassDefinition{
		class:  typ,
		fields: make([]string, 0, 0),
	}

	if err = o.addClassrefs(cd); err != nil {
		return nil, err
	}

	var (
		fieldKey   string
		fieldValue interface{}
		b          []byte
	)

	ci, ok := o.loadClassTypeSchema(cd.class)
	if !ok { // Generic java object
		jo := &JavaObject{
			class:  cd.class,
			names:  make([]string, 0, len(cd.fields)),
			values: make([]interface{}, 0, len(cd.fields)),
		}

		if err = o.addObjectrefs(jo); err != nil {
			return nil, err
		}

		b, err = reader.Peek(1)
		if err != nil {
			return jo, err
		}
		for b[0] != 'z' && err == nil {
			fieldKey, err = DecodeStringHessianV1(o, reader)
			if err != nil {
				return jo, err
			}
			if len(fieldKey) == 0 {
				return nil, ErrDecodeObjectFieldCannotBeNull
			}
			fieldValue, err = DecodeHessianV1(o, reader)
			if err != nil {
				return jo, err
			}
			jo.names = append(jo.names, fieldKey)
			jo.values = append(jo.values, fieldValue)
			b, err = reader.Peek(1)
		}

		// read z
		_, err = reader.ReadByte()
		if err != nil {
			return jo, err
		}

		return jo, nil
	}

	// Concrete type
	value := reflect.New(ci.base)
	structvalue := value.Elem()

	if err = o.addObjectrefs(value.Interface()); err != nil {
		return nil, err
	}

	b, err = reader.Peek(1)
	if err != nil {
		return nil, err
	}

	for b[0] != 'z' {
		fieldKey, err = DecodeStringHessianV1(o, reader)
		if err != nil {
			return nil, err
		}

		fi, ok := lookupReflectField(ci.base, fieldKey)
		if !ok {
			if o.disallowMissingField {
				return nil, errors.New("hessian: malformed class field (not found) " + fieldKey)
			}
			_, err = DecodeHessianV1(o, reader)
			if err != nil {
				return nil, err
			}

		} else {
			key := structvalue.Field(fi)
			if !key.CanSet() {
				return nil, errors.New("hessian: malformed class field (unassignable) " + fieldKey)
			}

			fieldValue, err = DecodeHessianV1(o, reader)
			if err != nil {
				return nil, err
			}

			err = safeSetValueByReflect(key, fieldValue)
			if err != nil {
				return nil, err
			}
		}

		b, err = reader.Peek(1)
		if err != nil {
			return nil, err
		}
	}

	// read z
	_, err = reader.ReadByte()
	if err != nil {
		return value.Interface(), err
	}

	return value.Interface(), nil
}

func decodeObjectDefinitionHessianV1(o *DecodeContext, reader *bufio.Reader) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodeobjectdefinition")
		defer o.tracer.OnTraceStop("decodeobjectdefinition")
	}

	u32, err := DecodeInt32HessianV1(o, reader)
	if err != nil {
		return err
	}

	// TODO(detailyang): config it to avoid DDOS
	typ := string(make([]byte, u32))
	err = readAtLeastBytesFromReader(reader, int(u32), s2b(typ))
	if err != nil {
		return err
	}

	// TODO(detailyang): cache class definition
	fieldslen, err := DecodeInt32HessianV1(o, reader)
	if err != nil {
		return err
	}

	if fieldslen < 0 || int(fieldslen) > o.GetMaxObjectFields() {
		return ErrDecodeMaxObjectFieldsExceeded
	}

	cd := ClassDefinition{
		class:  typ,
		fields: make([]string, 0, fieldslen),
	}

	for i := 0; i < int(fieldslen); i++ {
		field, err := DecodeStringHessianV1(o, reader)
		if err != nil {
			return err
		}
		cd.fields = append(cd.fields, field)
	}

	if err := o.addClassrefs(cd); err != nil {
		return err
	}

	return nil
}
