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
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeList(t *testing.T) {
	t.Run("should read untyped array", func(t *testing.T) {
		testDecodeList(t,
			[]byte{0x7e, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96},
			[]interface{}{
				int32(1), int32(2), int32(3),
				int32(4), int32(5), int32(6),
			},
		)
	})
	t.Run("should read untyped list", func(t *testing.T) {
		testDecodeList(t,
			readFile(t, "testdata/list/untyped_list.bin"),
			[]interface{}{
				int64(1), int64(2), "foo",
			},
		)

		testDecodeList(t,
			readFile(t, "testdata/list/untyped_list_8.bin"),
			[]interface{}{
				"1", "2", "3", "4", "5", "6", "7", "8",
			},
		)

		testDecodeList(t,
			readFile(t, "testdata/list/untyped_[].bin"),
			[]interface{}{},
		)

		testDecodeList(t,
			readFile(t, "testdata/list/untyped__String_foo,bar.bin"),
			[]interface{}{
				"foo",
				"bar",
			},
		)
	})

	t.Run("should write and read typed fixed-length list", func(t *testing.T) {
		testDecodeList(t,
			readFile(t, "testdata/list/typed_list.bin"),
			&JavaList{
				class: "hessian.demo.SomeArrayList",
				value: []interface{}{
					"ok",
					"some list",
				},
			},
		)
		testDecodeList(t,
			readFile(t, "testdata/list/typed_list_8.bin"),
			&JavaList{
				class: "hessian.demo.SomeArrayList",
				value: []interface{}{
					"1", "2", "3", "4", "5", "6", "7", "8",
				},
			},
		)

		testDecodeList(t,
			readFile(t, "testdata/list/[int.bin"),
			&JavaList{
				class: "[int",
				value: []interface{}{
					int64(1),
					int64(2),
					int64(3),
				},
			},
		)

		testDecodeList(t,
			readFile(t, "testdata/list/[string.bin"),
			&JavaList{
				class: "[string",
				value: []interface{}{
					"1",
					"@",
					"3",
				},
			},
		)
	})

	t.Run("should read enum lists", func(t *testing.T) {
		testDecodeList(t,
			readFile(t, "testdata/enum/lists.bin"),
			[]interface{}{
				&JavaObject{
					class:  "hessian.Main$Color",
					names:  []string{"name"},
					values: []interface{}{"BLUE"},
				},
				&JavaObject{
					class:  "hessian.Main$Color",
					names:  []string{"name"},
					values: []interface{}{"RED"},
				},
				&JavaObject{
					class:  "hessian.Main$Color",
					names:  []string{"name"},
					values: []interface{}{"GREEN"},
				},
			},
		)
	})
}

func testDecodeList(t *testing.T, b []byte, x interface{}) {
	o := NewDecodeContext().SetTracer(NewDummyTracer()).
		SetClassrefs(NewDecodeClassRefs()).
		SetTyperefs(NewDecodeTypeRefs()).
		SetObjectrefs(NewDecodeObjectRefs())
	y, err := DecodeListHessian4V2(o, bufio.NewReader(bytes.NewReader(b)))
	require.Nil(t, err)
	require.Equal(t, x, y)
}
