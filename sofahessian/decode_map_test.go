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
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeMap(t *testing.T) {
	hashmapbuf := []byte{
		0x48,
		0x01, // 1
		0x31,
		0x03,
		'f', 'e', 'e',
		0x02, // 16
		0x31,
		0x36,
		0x03,
		'f', 'i', 'e', // 'fie'
		0x03,
		0x32, // 256
		0x35,
		0x36,
		0x03,
		'f', 'o', 'e', // 'foe'
		'Z',
	}

	t.Run("should read car list", func(t *testing.T) {
		testDecodeMap(t,
			readFile(t, "testdata/map/car_list.bin"),
			[]interface{}{
				&JavaObject{
					class: "hessian.demo.Car",
					names: []string{"a", "c", "b", "model", "color", "mileage"},
					values: []interface{}{
						"a",
						"c",
						"b",
						"model 1",
						"aquamarine",
						int32(65536),
					},
				},
				&JavaObject{
					class: "hessian.demo.Car",
					names: []string{"a", "c", "b", "model", "color", "mileage"},
					values: []interface{}{
						"a",
						"c",
						"b",
						"model 2",
						"aquamarine",
						int32(65536),
					},
				},
				&JavaObject{
					class: "hessian.demo.Car",
					names: []string{"a", "c", "b", "model", "color", "mileage"},
					values: []interface{}{
						"a",
						"c",
						"b",
						"model 3",
						"aquamarine",
						int32(65536),
					},
				},
			},
		)

		cr := NewClassRegistry()
		cr.RegisterJavaClass(DemoCar{})

		d, err := DecodeHessian4V2(NewDecodeContext().
			SetTracer(NewDummyTracer()).
			SetClassRegistry(cr),
			bufio.NewReader(bytes.NewReader(readFile(t, "testdata/map/car_list.bin"))))
		require.Nil(t, err)
		require.Equal(t, []interface{}{
			&DemoCar{A: "a", C: "c", B: "b", Model: "model 1", Color: "aquamarine", Mileage: 65536},
			&DemoCar{A: "a", C: "c", B: "b", Model: "model 2", Color: "aquamarine", Mileage: 65536},
			&DemoCar{A: "a", C: "c", B: "b", Model: "model 3", Color: "aquamarine", Mileage: 65536},
		}, d)
	})

	t.Run("should read two car list", func(t *testing.T) {
		testDecodeMap(t,
			readFile(t, "testdata/map/two_car_list.bin"),
			[]interface{}{
				&JavaObject{
					class: "hessian.demo.Car",
					names: []string{"a", "c", "b", "model", "color", "mileage"},
					values: []interface{}{
						"a",
						"c",
						"b",
						"model 1",
						"aquamarine",
						int32(65536),
					},
				},
				&JavaObject{
					class: "hessian.demo.Car",
					names: []string{"a", "c", "b", "model", "color", "mileage"},
					values: []interface{}{
						"a",
						"c",
						"b",
						"model 2",
						"aquamarine",
						int32(65536),
					},
				},
			},
		)
	})

	t.Run("should read untyped map", func(t *testing.T) {
		testDecodeMap(t, hashmapbuf, map[interface{}]interface{}{
			"1":   "fee",
			"16":  "fie",
			"256": "foe",
		})
	})

	t.Run("should read a circular java object", func(t *testing.T) {
		testDecodeMap(t,
			readFile(t, "testdata/map/car.bin"),
			&JavaObject{
				class:  "hessian.demo.Car",
				names:  []string{"a", "c", "b", "model", "color", "mileage"},
				values: []interface{}{"a", "c", "b", "Beetle", "aquamarine", int32(65536)},
			},
		)
		jo := &JavaObject{
			class:  "hessian.demo.Car",
			names:  []string{"model", "color", "mileage", "self", "prev"},
			values: []interface{}{"Beetle", "aquamarine", int32(65536), interface{}(nil), interface{}(nil)},
		}
		jo.values[3] = jo

		testDecodeMap(t,
			readFile(t, "testdata/map/car1.bin"),
			jo,
		)

		testDecodeMap(t,
			readFile(t, "testdata/map/foo_empty.bin"),
			map[interface{}]interface{}{"foo": ""},
		)

		testDecodeMap(t,
			readFile(t, "testdata/map/foo_bar.bin"),
			map[interface{}]interface{}{
				"123": int32(456), "foo": "bar",
				"zero": int32(0), "中文key": "中文哈哈value",
			},
		)
	})

	t.Run("should decode map with type", func(t *testing.T) {
		testDecodeMap(t,
			readFile(t, "testdata/map/hashtable.bin"),
			&JavaMap{
				class: "java.util.Hashtable",
				m:     map[interface{}]interface{}{"foo": "bar", "中文key": "中文哈哈value"},
			},
		)
	})

	t.Run("should read java.util.HashMap", func(t *testing.T) {
		b, err := hex.DecodeString(`43302c636f6d2e74616f62616f2e63756e2e74726164652e726573756c746d6f64656c2e526573756c744d6f64656c940773756363657373086572726f724d7367096572726f72436f6465046461746160544e4e4843303d636f6d2e74616f62616f2e63756e2e74726164652e7472616465706c6174666f726d2e706172616d2e5061636b6167655175657279506172616d73564f93066974656d49640a616374697669747949640a6469766973696f6e4964614c000002006a3b9b5204313233343daef1e2613c411f04313233343daef13e895f5a`)
		require.Nil(t, err)

		key1 := &JavaObject{
			class: "com.taobao.cun.trade.tradeplatform.param.PackageQueryParamsVO",
			names: []string{"itemId", "activityId", "divisionId"},
			values: []interface{}{
				int64(2200805546834),
				"1234",
				int64(110321),
			},
		}

		key2 := &JavaObject{
			class: "com.taobao.cun.trade.tradeplatform.param.PackageQueryParamsVO",
			names: []string{"itemId", "activityId", "divisionId"},
			values: []interface{}{
				int64(16671),
				"1234",
				int64(110321),
			},
		}
		ym := map[interface{}]interface{}{
			key1: int64(2),
			key2: int64(166239),
		}
		jo := &JavaObject{
			class: "com.taobao.cun.trade.resultmodel.ResultModel",
			names: []string{"success", "errorMsg", "errorCode", "data"},
			values: []interface{}{
				bool(true),
				nil,
				nil,
				ym,
			},
		}
		d, err := DecodeHessian4V2(NewDecodeContext().SetTracer(NewDummyTracer()),
			bufio.NewReader(bytes.NewReader(b)))
		require.Nil(t, err)
		djo, ok := d.(*JavaObject)
		require.Equal(t, true, ok)
		require.Equal(t, jo.class, djo.class)
		require.Equal(t, jo.names, djo.names)
		require.Equal(t, jo.values[:3], djo.values[:3])
		x, ok := jo.values[3].(map[interface{}]interface{})
		require.Equal(t, true, ok)
		for i := range x {
			require.Equal(t, ym[i], x[i])
		}
	})
}

func testDecodeMap(t *testing.T, b []byte, i interface{}) {
	d, err := DecodeHessian4V2(NewDecodeContext().SetTracer(NewDummyTracer()),
		bufio.NewReader(bytes.NewReader(b)))
	require.Nil(t, err)
	require.Equal(t, i, d)
}
