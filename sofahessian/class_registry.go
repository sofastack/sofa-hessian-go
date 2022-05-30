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
	"reflect"
	"sync"
)

type ClassTypeSchema struct {
	raw   reflect.Type // the origin type
	base  reflect.Type // the base type after all indirections ( ***T => T)
	ebase reflect.Type // the type of element for [T]
	indir int          // number of indirections to reach the base type
}

var globalClassRegistry ClassRegistry

func Load(name string) (*ClassTypeSchema, bool) {
	return globalClassRegistry.Load(name)
}

func Register(name string, value interface{}) (bool, error) {
	return globalClassRegistry.Register(name, value)
}

func RegisterJavaClass(value JavaClassNameGetter) (bool, error) {
	return globalClassRegistry.Register(value.GetJavaClassName(), value)
}

type ClassRegistry struct {
	types sync.Map // map[string]*typeInfo
}

func NewClassRegistry() *ClassRegistry { return &ClassRegistry{} }

func (tr *ClassRegistry) Load(name string) (*ClassTypeSchema, bool) {
	i, ok := tr.types.Load(name)
	if ok {
		return i.(*ClassTypeSchema), true
	}
	return nil, false
}

func (tr *ClassRegistry) RegisterJavaClass(value JavaClassNameGetter) (bool, error) {
	return tr.Register(value.GetJavaClassName(), value)
}

func (tr *ClassRegistry) Register(name string, value interface{}) (bool, error) {
	_, ok := tr.types.Load(name)
	if ok {
		return false, nil
	}

	rt := reflect.TypeOf(value)
	rv := reflect.ValueOf(value)
	rv = decAllocReflectValue(rv)
	ut := new(ClassTypeSchema)
	ut.base = rt
	ut.raw = rt
	slowpoke := ut.base // walks half as fast as ut.base
	for {
		pt := ut.base
		if pt.Kind() != reflect.Ptr {
			break
		}
		ut.base = pt.Elem()
		if ut.base == slowpoke { // ut.base lapped slowpoke
			// recursive pointer type.
			return false, errors.New("hessian: can't represent recursive pointer type " + ut.base.String())
		}
		if ut.indir%2 == 0 {
			slowpoke = slowpoke.Elem()
		}
		ut.indir++
	}

	tr.types.Store(name, ut)

	// recursive register the class field type
	if kind := ut.base.Kind(); kind == reflect.Struct {
		for i := 0; i < ut.base.NumField(); i++ {
			field := ut.base.Field(i)
			ok, _ := implementsInterface(field.Type, JavaClassNameGetterInterfaceType)
			if ok {
				value := rv.Field(i)
				if value.CanInterface() {
					className := ""
					if safeIsNil(value) {
						// dereferencing pointer to T
						value = decAllocReflectValue(reflect.New(field.Type))
						// malloc *T
						value = reflect.New(reflect.TypeOf(value.Interface()))
						className = value.Interface().(JavaClassNameGetter).GetJavaClassName()
					} else {
						className = value.Interface().(JavaClassNameGetter).GetJavaClassName()
					}
					if _, err := tr.Register(className, value.Interface()); err != nil {
						panic("bug: failed to register class")
					}
				}
			}
		}
	} else if kind == reflect.Slice || kind == reflect.Array {
		// recursive register the slice element type
		ut.ebase = ut.base.Elem()
		ok, _ := implementsInterface(ut.ebase, JavaClassNameGetterInterfaceType)
		if ok {
			value := reflect.New(ut.ebase)
			className := ""
			className = value.Interface().(JavaClassNameGetter).GetJavaClassName()
			if _, err := tr.Register(className, value.Interface()); err != nil {
				panic("bug: failed to register class")
			}
		}
	}

	return true, nil
}

// implementsInterface reports whether the type implements the concrete interface.
// It also returns the number of indirections required to get to the
// implementation.
func implementsInterface(typ, interfaceType reflect.Type) (success bool, indir int8) {
	if typ == nil {
		return
	}
	rt := typ
	// The type might be a pointer and we need to keep
	// dereferencing to the base type until we find an implementation.
	for {
		if rt.Implements(interfaceType) {
			return true, indir
		}
		if p := rt; p.Kind() == reflect.Ptr {
			indir++
			if indir > 100 { // insane number of indirections
				return false, 0
			}
			rt = p.Elem()
			continue
		}
		break
	}
	// No luck yet, but if this is a base type (non-pointer), the pointer might satisfy.
	if typ.Kind() != reflect.Ptr {
		// Not a pointer, but does the pointer work?
		if reflect.PtrTo(typ).Implements(interfaceType) {
			return true, -1
		}
	}
	return false, 0
}
