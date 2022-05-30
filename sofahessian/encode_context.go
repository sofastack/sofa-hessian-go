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
	"errors"
)

// EncodeContext holds the context of encoding.
type EncodeContext struct {
	maxdepth          int
	depth             int
	version           Version
	classrefs         *EncodeClassrefs
	typerefs          *EncodeTypeRefs
	objectrefs        *EncodeObjectrefs
	disableObjectrefs bool
	tracer            Tracer
	less              func(keyi, keyj, valuei, valuej interface{}) bool
}

// NewEncodeContext returns a new EncodeContext.
func NewEncodeContext() *EncodeContext {
	return &EncodeContext{
		classrefs:  NewEncodeClassrefs(),
		typerefs:   NewEncodeTyperefs(),
		objectrefs: NewEncodeObjectrefs(),
	}
}

func (e *EncodeContext) addDepth() { e.depth++ }
func (e *EncodeContext) subDepth() { e.depth-- }

func (e *EncodeContext) SetMaxDepth(depth int) *EncodeContext {
	e.maxdepth = depth
	return e
}

// Reset resets the context.
func (e *EncodeContext) Reset() {
	e.version = 0
	e.classrefs.Reset()
	e.typerefs.Reset()
	e.objectrefs.Reset()
	e.disableObjectrefs = false
	e.tracer = nil
	e.less = nil
}

func (e *EncodeContext) GetVersion() Version {
	return e.version
}

// SetVersion sets the version to encode.
func (e *EncodeContext) SetVersion(ver Version) *EncodeContext {
	e.version = ver
	return e
}

func (e *EncodeContext) DisableObjectrefs() *EncodeContext {
	e.disableObjectrefs = true
	return e
}

func (e *EncodeContext) SetLessFunc(fn func(keyi, keyj, valuei, valuej interface{}) bool) *EncodeContext {
	e.less = fn
	return e
}

func (e *EncodeContext) SetTracer(tracer Tracer) *EncodeContext {
	e.tracer = tracer
	return e
}

func (e *EncodeContext) SetClassrefs(ref *EncodeClassrefs) *EncodeContext {
	e.classrefs = ref
	return e
}

func (e *EncodeContext) SetObjectrefs(ref *EncodeObjectrefs) *EncodeContext {
	e.objectrefs = ref
	return e
}

func (e *EncodeContext) SetTyperefs(ref *EncodeTypeRefs) *EncodeContext {
	e.typerefs = ref
	return e
}

func (e *EncodeContext) getTyperefs(typ string) (int, bool, error) {
	if e.typerefs == nil {
		return -1, false, errors.New("hessian: encode type references is nil")
	}

	n, ok := e.typerefs.Get(typ)
	return n, ok, nil
}

func (e *EncodeContext) addTyperefs(typ string) error {
	if e.typerefs == nil {
		return errors.New("hessian: encode type references is nil")
	}

	e.typerefs.Set(typ)
	return nil
}

func (e *EncodeContext) getClassrefs(cls interface{}) (int, error) {
	if e.classrefs == nil {
		return -1, errors.New("hessian: encode class references is nil")
	}

	return e.classrefs.Get(cls), nil
}

func (e *EncodeContext) addClassrefs(cls interface{}) (ref int, referenced bool, err error) {
	if e.classrefs == nil {
		return -1, false, errors.New("hessian: encode class references is nil")
	}

	ref, referenced = e.classrefs.Add(cls)
	return ref, referenced, nil
}

func (e *EncodeContext) getObjectrefs(cls interface{}) (int, error) {
	if e.objectrefs == nil {
		return -1, errors.New("hessian: encode object references is nil")
	}

	return e.objectrefs.Get(cls), nil
}

func (e *EncodeContext) addObjectrefs(obj interface{}) (int, error) {
	if e.objectrefs == nil {
		return -1, errors.New("hessian: encode object references is nil")
	}

	return e.objectrefs.Add(obj), nil
}
