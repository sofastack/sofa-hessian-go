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
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeNilMap(t *testing.T) {
	t.Run("test nil map", func(t *testing.T) {
		dst, err := EncodeMapToHessian4V2(&EncodeContext{}, nil, nil)
		require.Nil(t, err)
		require.Equal(t, []byte("N"), dst)
	})
}

func TestEncodeMapExample(t *testing.T) {
	m := make(map[int32]string, 256)
	m[1] = "fee"
	m[16] = "fie"
	m[256] = "foe"

	// map = new HashMap();
	// map.put(new Integer(1), "fee");
	// map.put(new Integer(16), "fie");
	// map.put(new Integer(256), "foe");

	// ---

	// H           # untyped map (HashMap for Java)
	//   x91       # 1
	//   x03 fee   # "fee"

	//   xa0       # 16
	//   x03 fie   # "fie"

	//   xc9 x00   # 256
	//   x03 foe   # "foe"

	//   Z

	b := bytes.NewBuffer(nil)
	b.WriteByte('H')
	b.WriteByte(0x91)
	b.WriteByte(0x03)
	b.WriteString("fee")
	b.WriteByte(0xa0)
	b.WriteByte(0x03)
	b.WriteString("fie")
	b.WriteByte(0xc9)
	b.WriteByte(0x00)
	b.WriteByte(0x03)
	b.WriteString("foe")
	b.WriteString("Z")

	dst, err := EncodeToHessian4V2(&EncodeContext{
		objectrefs: NewEncodeObjectrefs(),
		less: func(ki, kj, vi, vj interface{}) bool {
			ii := ki.(int32)
			ij := kj.(int32)
			return ii <= ij
		},
	}, nil, m)
	require.Nil(t, err)
	require.Equal(t, hex.EncodeToString(b.Bytes()), hex.EncodeToString(dst))

	// public class Car implements Serializable {
	// 	String color = "aquamarine";
	// 	String model = "Beetle";
	// 	int mileage = 65536;
	//   }

	//   ---
	//   M
	// 	x13 com.caucho.test.Car  # type

	// 	x05 color                # color field
	// 	x0a aquamarine

	// 	x05 model                # model field
	// 	x06 Beetle

	// 	x07 mileage              # mileage field
	// 	I x00 x01 x00 x00
	// 	Z
	b.Reset()
	b.WriteByte('M')
	b.WriteByte(0x13)
	b.WriteString("com.caucho.test.Car")
	b.WriteByte(0x05)
	b.WriteString("color")
	b.WriteByte(0x0a)
	b.WriteString("aquamarine")
	b.WriteByte(0x05)
	b.WriteString("model")
	b.WriteByte(0x06)
	b.WriteString("Beetle")
	b.WriteByte(0x07)
	b.WriteString("mileage")
	// http://hessian.caucho.com/doc/hessian-serialization.html##type-map
	// Maybe it's wrong?
	b.WriteByte('I')
	b.WriteByte(0x00)
	b.WriteByte(0x01)
	b.WriteByte(0x00)
	b.WriteByte(0x00)
	b.WriteByte('Z')

	x := map[string]interface{}{
		"color":   "aquamarine",
		"model":   "Beetle",
		"mileage": int32(65536),
	}

	y := map[string]int{
		"color":   1,
		"model":   2,
		"mileage": 3,
	}

	c := CarMap(x)
	dst, err = EncodeToHessian4V2(&EncodeContext{
		typerefs:   NewEncodeTyperefs(),
		objectrefs: NewEncodeObjectrefs(),
		tracer:     NewDummyTracer(),
		less: func(ki, kj, vi, vj interface{}) bool {
			ii := ki.(string)
			ij := kj.(string)

			return y[ii] <= y[ij]
		},
	}, nil, c)
	require.Nil(t, err)
	// require.Equal(t, hex.EncodeToString(b.Bytes()), hex.EncodeToString(dst))
}

type CarMap map[string]interface{}

func (c CarMap) GetJavaClassName() string {
	return "com.caucho.test.Car"
}
