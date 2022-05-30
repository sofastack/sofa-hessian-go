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
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type SomeArrayList []string

func (s SomeArrayList) GetJavaClassName() string {
	return "hessian.demo.SomeArrayList"
}

type IntArray []int

func (i IntArray) GetJavaClassName() string {
	return "[int"
}

type StringArray []string

func (i StringArray) GetJavaClassName() string {
	return "[string"
}

func TestEncodeToSliceStruct(t *testing.T) {
	blue := &Color{
		Name: "BLUE",
	}

	red := &Color{
		Name: "RED",
	}

	green := &Color{
		Name: "GREEN",
	}
	ecf := NewEncodeClassrefs()
	o := &EncodeContext{
		classrefs:  ecf,
		typerefs:   NewEncodeTyperefs(),
		objectrefs: NewEncodeObjectrefs(),
		tracer:     NewDummyTracer(),
	}
	dst, err := EncodeToHessian4V2(o, nil, []*Color{
		blue, red, green,
	})
	require.Nil(t, err)

	s := fmt.Sprintf("%s.bin", "lists")
	fn := filepath.Join("testdata", "enum", s)
	data, err := ioutil.ReadFile(fn)
	require.Nil(t, err)
	require.Equal(t, data, dst, fmt.Sprintf("%s", fn))
}

func TestEncodeList(t *testing.T) {
	t.Run("testdata/list/untyped_list.bin", func(t *testing.T) {
		testEncodeList(t, []interface{}{
			1,
			2,
			"foo",
		}, "testdata/list/untyped_list.bin")
	})

	t.Run("testdata/list/untyped_[].bin", func(t *testing.T) {
		testEncodeList(t, []interface{}{}, "testdata/list/untyped_[].bin")
	})

	t.Run("test list reference", func(t *testing.T) {
		o := &EncodeContext{
			classrefs:  NewEncodeClassrefs(),
			typerefs:   NewEncodeTyperefs(),
			objectrefs: NewEncodeObjectrefs(),
			tracer:     &DummyTracer{},
		}
		x := []string{"foo"}
		dst, err := EncodeToHessian4V2(o, nil, x)
		require.Nil(t, err)
		dst, err = EncodeToHessian4V2(o, dst, x)
		require.Nil(t, err)
		require.Equal(t, "7903666f6f5190", hex.EncodeToString(dst))
	})

	t.Run("should write and read typed fixed-length list", func(t *testing.T) {
		x := []string{"ok", "some list"}
		z := SomeArrayList(x)
		testEncodeList(t, z, "testdata/list/typed_list.bin")
	})

	t.Run("should write and read typed fixed-length list", func(t *testing.T) {
		x := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
		z := SomeArrayList(x)
		testEncodeList(t, z, "testdata/list/typed_list_8.bin")
	})

	t.Run("should write and read typed fixed-length list", func(t *testing.T) {
		x := []int{1, 2, 3}
		z := IntArray(x)
		testEncodeList(t, z, "testdata/list/[int.bin")
	})

	t.Run("should write and read typed fixed-length list", func(t *testing.T) {
		x := []string{"1", "@", "3"}
		z := StringArray(x)
		testEncodeList(t, z, "testdata/list/[string.bin")
	})

	t.Run("test nil", func(t *testing.T) {
		dst, err := EncodeListToHessian4V2(&EncodeContext{}, nil, nil)
		require.Nil(t, err)
		require.Equal(t, []byte("N"), dst)
	})
}

func testEncodeList(t *testing.T, i interface{}, filename string) {
	dst, err := EncodeToHessian4V2(&EncodeContext{
		classrefs:  NewEncodeClassrefs(),
		typerefs:   NewEncodeTyperefs(),
		objectrefs: NewEncodeObjectrefs(),
		tracer:     &DummyTracer{},
	}, nil, i)
	require.Nil(t, err)
	require.Equal(t, hex.EncodeToString(readFile(t, filename)), hex.EncodeToString(dst))
}
