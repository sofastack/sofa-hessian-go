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

// IdentityIntMap implements the identity map of java in hessian with Golang.
type IdentityIntMap struct {
	m map[interface{}]int
}

func NewIdentityIntMap() *IdentityIntMap {
	return &IdentityIntMap{
		m: make(map[interface{}]int, 256),
	}
}

func (m *IdentityIntMap) Size() int {
	return len(m.m)
}

func (m *IdentityIntMap) Get(key interface{}) int {
	value, ok := m.m[key]
	if ok {
		return value
	}

	return -1
}

func (m *IdentityIntMap) Put(key interface{}, value int, replaced bool) int {
	oldvalue, ok := m.m[key]
	if ok {
		if replaced {
			m.m[key] = value
			return value
		}
		return oldvalue
	}

	m.m[key] = value
	return value
}
