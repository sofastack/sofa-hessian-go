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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type Car struct {
	Color string `hessian:"color"`
	Model string `hessian:"model"`
}

func (c Car) GetJavaClassName() string {
	return "example.Car"
}

type Color struct {
	Name string `hessian:"name"`
}

func (c Color) GetJavaClassName() string {
	return "hessian.Main$Color"
}

type DemoCar struct {
	A       string `hessian:"a"`
	C       string `hessian:"c"`
	B       string `hessian:"b"`
	Model   string `hessian:"model"`
	Color   string `hessian:"color"`
	Mileage int32  `hessian:"mileage"`
}

func (c DemoCar) GetJavaClassName() string {
	return "hessian.demo.Car"
}

func TestEncodeNilObject(t *testing.T) {
	t.Run("test nil object", func(t *testing.T) {
		dst, err := EncodeObjectToHessian4V2(&EncodeContext{}, nil, nil)
		require.Nil(t, err)
		require.Equal(t, []byte("N"), dst)
	})
}

func TestEncodeObjectSanity(t *testing.T) {
	t.Run("test red green blue golden", func(t *testing.T) {
		for _, n := range []string{
			"red",
			"green",
			"blue",
		} {
			s := fmt.Sprintf("%s.bin", n)
			fn := filepath.Join("testdata", "enum", s)
			data, err := ioutil.ReadFile(fn)
			require.Nil(t, err)
			color := &Color{
				Name: strings.ToUpper(n),
			}
			require.Nil(t, err)
			o := &EncodeContext{
				classrefs:  NewEncodeClassrefs(),
				typerefs:   NewEncodeTyperefs(),
				objectrefs: NewEncodeObjectrefs(),
			}
			dst, err := EncodeToHessian4V2(o, nil, color)
			require.Nil(t, err)
			require.Equal(t, data, dst, fmt.Sprintf("%s", fn))
		}
	})

	t.Run("read one car list", func(t *testing.T) {
		car1 := DemoCar{
			A:       "a",
			C:       "c",
			B:       "b",
			Model:   "model 1",
			Color:   "aquamarine",
			Mileage: 65536,
		}

		o := &EncodeContext{
			classrefs:  NewEncodeClassrefs(),
			typerefs:   NewEncodeTyperefs(),
			objectrefs: NewEncodeObjectrefs(),
			tracer:     &DummyTracer{},
		}
		dst, err := EncodeToHessian4V2(o, nil, []interface{}{car1})
		require.Nil(t, err)
		require.Equal(t, hex.EncodeToString(readFile(t, "testdata/map/one_car_list.bin")), hex.EncodeToString(dst))
	})

	t.Run("read two car list", func(t *testing.T) {
		car1 := DemoCar{
			A:       "a",
			C:       "c",
			B:       "b",
			Model:   "model 1",
			Color:   "aquamarine",
			Mileage: 65536,
		}
		car2 := DemoCar{
			A:       "a",
			C:       "c",
			B:       "b",
			Model:   "model 2",
			Color:   "aquamarine",
			Mileage: 65536,
		}

		o := &EncodeContext{
			classrefs:  NewEncodeClassrefs(),
			typerefs:   NewEncodeTyperefs(),
			objectrefs: NewEncodeObjectrefs(),
			tracer:     &DummyTracer{},
		}
		dst, err := EncodeToHessian4V2(o, nil, []interface{}{car1, car2})
		require.Nil(t, err)
		require.Equal(t, hex.EncodeToString(readFile(t, "testdata/map/two_car_list.bin")), hex.EncodeToString(dst))
	})

	t.Run("read three car list", func(t *testing.T) {
		car1 := DemoCar{
			A:       "a",
			C:       "c",
			B:       "b",
			Model:   "model 1",
			Color:   "aquamarine",
			Mileage: 65536,
		}
		car2 := DemoCar{
			A:       "a",
			C:       "c",
			B:       "b",
			Model:   "model 2",
			Color:   "aquamarine",
			Mileage: 65536,
		}
		car3 := DemoCar{
			A:       "a",
			C:       "c",
			B:       "b",
			Model:   "model 3",
			Color:   "aquamarine",
			Mileage: 65536,
		}

		o := &EncodeContext{
			classrefs:  NewEncodeClassrefs(),
			typerefs:   NewEncodeTyperefs(),
			objectrefs: NewEncodeObjectrefs(),
			tracer:     &DummyTracer{},
		}
		dst, err := EncodeToHessian4V2(o, nil, []interface{}{car1, car2, car3})
		require.Nil(t, err)
		require.Equal(t, hex.EncodeToString(readFile(t, "testdata/map/car_list.bin")), hex.EncodeToString(dst))
	})

	t.Run("java.util.concurrent.atomic.AtomicLong", func(t *testing.T) {
		x := AtomicLong{0}
		o := &EncodeContext{
			classrefs:  NewEncodeClassrefs(),
			typerefs:   NewEncodeTyperefs(),
			objectrefs: NewEncodeObjectrefs(),
			tracer:     &DummyTracer{},
		}
		dst, err := EncodeToHessian4V2(o, nil, x)
		require.Nil(t, err)

		require.Equal(t, hex.EncodeToString(readFile(t, "testdata/object/AtomicLong0.bin")), hex.EncodeToString(dst))
	})
}

type AtomicLong struct {
	Value int64 `hessian:"value"`
}

func (al AtomicLong) GetJavaClassName() string {
	return "java.util.concurrent.atomic.AtomicLong"
}

func TestEncodeObject(t *testing.T) {
	car1 := Car{
		"red", "corvette",
	}
	o := &EncodeContext{
		classrefs:  NewEncodeClassrefs(),
		typerefs:   NewEncodeTyperefs(),
		objectrefs: NewEncodeObjectrefs(),
	}
	dst, err := EncodeToHessian4V2(o, nil, car1)
	require.Nil(t, err)
	// 	C                        # object definition (#0)
	//   x0b example.Car        # type is example.Car
	//   x92                    # two fields
	//   x05 color              # color field name
	//   x05 model              # model field name

	// O                        # object def (long form)
	//   x90                    # object definition #0
	//   x03 red                # color field value
	//   x08 corvette           # model field value

	// x60                      # object def #0 (short form)
	//   x05 green              # color field value
	//   x05 civic              # model field value
	require.Equal(t, "430b6578616d706c652e4361729205636f6c6f72056d6f64656c600372656408636f727665747465",
		hex.EncodeToString(dst))
	car2 := &Car{
		"red", "corvette",
	}
	o = &EncodeContext{
		classrefs:  NewEncodeClassrefs(),
		typerefs:   NewEncodeTyperefs(),
		objectrefs: NewEncodeObjectrefs(),
	}
	dst, err = EncodeToHessian4V2(o, nil, car2)
	require.Nil(t, err)
	require.Equal(t, "430b6578616d706c652e4361729205636f6c6f72056d6f64656c600372656408636f727665747465",
		hex.EncodeToString(dst))
}

func TestEncodeRefObjectTo(t *testing.T) {
	car1 := Car{
		"red", "corvette",
	}
	car2 := Car{
		"green", "civic",
	}
	o := &EncodeContext{
		classrefs:  NewEncodeClassrefs(),
		typerefs:   NewEncodeTyperefs(),
		objectrefs: NewEncodeObjectrefs(),
		tracer:     &DummyTracer{},
	}
	dst, err := EncodeToHessian4V2(o, nil, car1)
	require.Nil(t, err)
	dst, err = EncodeToHessian4V2(o, dst, car2)
	require.Nil(t, err)
	// class Car {
	// 	String color;
	// 	String model;
	//   }

	//   out.writeObject(new Car("red", "corvette"));
	//   out.writeObject(new Car("green", "civic"));

	//   ---

	//   C                        # object definition (#0)
	// 	x0b example.Car        # type is example.Car
	// 	x92                    # two fields
	// 	x05 color              # color field name
	// 	x05 model              # model field name

	//   O                        # object def (long form)
	// 	x90                    # object definition #0
	// 	x03 red                # color field value
	// 	x08 corvette           # model field value

	//   x60                      # object def #0 (short form)
	// 	x05 green              # color field value
	// 	x05 civic              # model field value

	b := bytes.NewBuffer(nil)
	b.WriteByte('C')
	b.WriteByte(0x0b)
	b.WriteString("example.Car")
	b.WriteByte(0x92)
	b.WriteByte(0x05)
	b.WriteString("color")
	b.WriteByte(0x05)
	b.WriteString("model")
	b.WriteByte(0x60)
	b.WriteByte(0x03)
	b.WriteString("red")
	b.WriteByte(0x08)
	b.WriteString("corvette")
	b.WriteByte(0x60)
	b.WriteByte(0x05)
	b.WriteString("green")
	b.WriteByte(0x05)
	b.WriteString("civic")

	require.Equal(t, hex.EncodeToString(b.Bytes()), hex.EncodeToString(dst))
}
