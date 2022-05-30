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

const (
	defaultMaxDecodeListLength   = 8192
	defaultMaxDecodeObjectFields = 1024
)

// DecodeContext holds the context of decoding.
type DecodeContext struct {
	depth                int
	maxdepth             int
	maxlistlength        int
	maxobjectfields      int
	version              Version
	disallowMissingField bool
	classrefs            *DecodeClassRefs
	typerefs             *DecodeTypeRefs
	objectrefs           *DecodeObjectRefs
	classRegistry        *ClassRegistry
	tracer               Tracer
}

func NewDecodeContext() *DecodeContext {
	return &DecodeContext{
		classrefs:  NewDecodeClassRefs(),
		typerefs:   NewDecodeTypeRefs(),
		objectrefs: NewDecodeObjectRefs(),
	}
}

func (d *DecodeContext) DisallowMissingField() *DecodeContext {
	d.disallowMissingField = true
	return d
}

func (d *DecodeContext) SetMaxListLength(m int) *DecodeContext {
	d.maxlistlength = m
	return d
}

func (d *DecodeContext) GetMaxListLength() int {
	if d.maxlistlength <= 0 {
		return defaultMaxDecodeListLength
	}
	return d.maxlistlength
}

func (d *DecodeContext) GetMaxObjectFields() int {
	if d.maxobjectfields <= 0 {
		return defaultMaxDecodeObjectFields
	}
	return d.maxobjectfields
}

func (d *DecodeContext) addDepth() { d.depth++ }
func (d *DecodeContext) subDepth() { d.depth-- }

func (d *DecodeContext) SetMaxDepth(depth int) *DecodeContext {
	d.maxdepth = depth
	return d
}

func (d *DecodeContext) Reset() {
	d.version = 0
	d.classrefs.Reset()
	d.typerefs.Reset()
	d.objectrefs.Reset()
	d.tracer = nil
	d.classRegistry = nil
}

func (d *DecodeContext) GetVersion() Version { return d.version }

func (d *DecodeContext) SetVersion(ver Version) *DecodeContext {
	d.version = ver
	return d
}

func (d *DecodeContext) SetClassRegistry(cr *ClassRegistry) *DecodeContext {
	d.classRegistry = cr
	return d
}

func (d *DecodeContext) SetTracer(tracer Tracer) *DecodeContext {
	d.tracer = tracer
	return d
}

func (d *DecodeContext) SetClassrefs(refs *DecodeClassRefs) *DecodeContext {
	d.classrefs = refs
	return d
}

func (d *DecodeContext) SetTyperefs(refs *DecodeTypeRefs) *DecodeContext {
	d.typerefs = refs
	return d
}

func (d *DecodeContext) SetObjectrefs(refs *DecodeObjectRefs) *DecodeContext {
	d.objectrefs = refs
	return d
}

func (d DecodeContext) loadClassTypeSchema(name string) (*ClassTypeSchema, bool) {
	// Check option
	if d.classRegistry != nil {
		ci, ok := d.classRegistry.Load(name)
		if ok {
			return ci, ok
		}
	}

	// Check global
	ci, ok := Load(name)
	return ci, ok
}

func (d *DecodeContext) getClassrefs(i int) (ClassDefinition, error) {
	if d.classrefs == nil {
		return ClassDefinition{}, ErrDecodeClassRefsIsNil
	}

	cd, ok := d.classrefs.Get(i)
	if !ok {
		return ClassDefinition{}, ErrDecodeClassRefsOverflow
	}
	return cd, nil
}

func (d *DecodeContext) addClassrefs(cd ClassDefinition) error {
	if d.classrefs == nil {
		return ErrDecodeClassRefsIsNil
	}

	d.classrefs.Append(cd)
	return nil
}

func (d *DecodeContext) getObjectrefs(refid int) (interface{}, error) {
	if d.objectrefs == nil {
		return nil, ErrDecodeObjectRefsIsNil
	}

	obj, ok := d.objectrefs.Get(refid)
	if !ok {
		return nil, ErrDecodeObjectRefsOverflow
	}

	return obj, nil
}

func (d *DecodeContext) addObjectrefs(i interface{}) error {
	if d.objectrefs == nil {
		return ErrDecodeObjectRefsIsNil
	}
	d.objectrefs.Append(i)

	return nil
}

func (d *DecodeContext) addTyperefs(typ string) error {
	if d.typerefs == nil {
		return ErrDecodeTypeRefsIsNil
	}
	d.typerefs.Append(typ)

	return nil
}

func (d *DecodeContext) getTyperefs(i int) (string, error) {
	if d.typerefs == nil {
		return "", ErrDecodeTypeRefsIsNil
	}

	typ, ok := d.typerefs.Get(i)
	if !ok {
		return "", ErrDecodeTypeRefsOverflow
	}

	return typ, nil
}
