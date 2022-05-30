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
	"encoding/json"
	"testing"
	"time"

	"github.com/sofastack/sofa-hessian-go/javaobject"
	"github.com/stretchr/testify/require"
)

type structfoo struct {
	A int64 `json:"A"`
}

type complex struct {
	Bool         bool
	Boolp        *bool
	Boolpp       **bool
	Int32        int32
	Int32p       *int32
	Int32pp      **int32
	Int64        int64
	Int64p       *int64
	Int64pp      **int64
	Float64      float64
	Float64p     *float64
	Float64pp    **float64
	String       string
	Stringp      *string
	Stringpp     **string
	Time         time.Time
	Timep        *time.Time
	Timepp       **time.Time
	Binary       []byte
	Binaryp      *[]byte
	Binarypp     **[]byte
	Interface    interface{}
	Interfacep   *interface{}
	Interfacepp  **interface{}
	Foo          structfoo
	Foop         *structfoo
	Foopp        **structfoo
	Map          map[string]string
	Mapp         *map[string]string
	Mappp        **map[string]string
	UntypedSlice []string
	Slice        []interface{}
	Slicep       *[]interface{}
	Slicepp      **[]interface{}
}

func (c complex) GetJavaClassName() string { return "test.comple" }

type TestTypedMap map[interface{}]interface{}

func (t TestTypedMap) GetJavaClassName() string {
	return "test.typedmap"
}

type RequestContext struct {
	Id   int64              `hessian:"id"`
	This *ConnectionRequest `hessian:"this$0"`
}

func (c RequestContext) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionRequest$RequestContext"
}

type ConnectionRequest struct {
	Ctx                *RequestContext `hessian:"ctx"`
	FromAppKey         string          `hessian:"fromAppKey"`
	ToAppKey           string          `hessian:"toAppKey"`
	EncryptedToken     []byte          `hessian:"encryptedToken"`
	ApplicationRequest interface{}     `hessian:"-"`
}

func (c ConnectionRequest) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionRequest"
}

func TestEncodeCodeGolden(t *testing.T) {
	t.Run("should encode decode success", func(t *testing.T) {
		connRequest := &ConnectionRequest{}
		connRequest.Ctx = &RequestContext{}
		connRequest.Ctx.Id = 1
		connRequest.Ctx.This = connRequest
		connRequest.ApplicationRequest = "aa"
		encoder := NewEncoder(NewEncodeContext().SetVersion(Hessian3xV2).SetMaxDepth(8))
		err := encoder.Encode(connRequest)
		require.Nil(t, err)
		buffer2 := bytes.NewReader(encoder.Bytes())
		decoder := NewDecoder(NewDecodeContext().SetVersion(Hessian3xV2), buffer2)
		d, err := decoder.Decode()
		if err != nil {
			t.Fatal(err)
		}
		_ = d
	})
}

func TestEncodeDecode(t *testing.T) {
	t.Run("should support circle data structures", func(t *testing.T) {
		connResponse := &ConnectionResponse{}
		connResponse.Ctx = &ResponseContext{}
		connResponse.Ctx.Id = 1
		connResponse.Ctx.This = connResponse
		testHessian4V2EncodeDecode(t, connResponse, connResponse)
		testHessian3V2EncodeDecode(t, connResponse, connResponse)
		testHessianV1EncodeDecode(t, connResponse, connResponse)
	})

	t.Run("shoudl encode and decode typed list", func(t *testing.T) {
		x := []javaobject.JavaLangStackTraceElement{
			{
				DeclaringClass: "a",
				MethodName:     "b",
				FileName:       "c",
				LineNumber:     333,
			},
			{
				DeclaringClass: "1a",
				MethodName:     "1b",
				FileName:       "1c",
				LineNumber:     1333,
			},
		}
		y := javaobject.JavaLangStackTraceElements(x)

		testHessian4V2EncodeDecode(t, y, y)
		testHessian3V2EncodeDecode(t, y, y)
		testHessianV1EncodeDecode(t, y, y)
	})

	t.Run("should test encode and decode concrete type", func(t *testing.T) {
		ctx := &HessianConnectionRequestContext{
			Id: 101,
		}
		ctxp := &ctx
		ctxpp := &ctxp
		ctxppp := &ctxpp
		ctxpppp := &ctxppp

		hr := &HessianConnectionRequest{Ctx: ctxpppp}

		testHessian3V2EncodeDecode(t, hr, hr)
		testHessian4V2EncodeDecode(t, hr, hr)
		testHessianV1EncodeDecode(t, hr, hr)
	})

	t.Run("should test object with complex object", func(t *testing.T) {
		type Bar struct {
			C int32
		}

		type Foo struct {
			A     int32
			B     Bar
			M     map[interface{}]interface{}
			Int32 []int32
		}

		f := Foo{
			A: 1,
			B: Bar{
				C: 2,
			},
			M: map[interface{}]interface{}{
				"a": 1,
				"b": "c",
				"z": []int32{1, 2, 3},
			},
			Int32: []int32{1, 2, 3},
		}

		y := &JavaObject{
			class: "Foo",
			names: []string{"A", "B", "M", "Int32"},
			values: []interface{}{
				int32(1),
				&JavaObject{
					class: "Bar",
					names: []string{"C"},
					values: []interface{}{
						int32(2),
					},
				},
				map[interface{}]interface{}{
					"z": []interface{}{
						int32(1),
						int32(2),
						int32(3),
					},
					"a": int64(1),
					"b": "c",
				},
				[]interface{}{
					int32(1),
					int32(2),
					int32(3),
				},
			},
		}

		testHessian4V2EncodeDecode(t, f, y)
		testHessian3V2EncodeDecode(t, f, y)
		testHessianV1EncodeDecode(t, f, y)
	})

	t.Run("should test string with utf8", func(t *testing.T) {
		x := "cc😊cc"
		testHessian4V2EncodeDecode(t, x, x)
		testHessian3V2EncodeDecode(t, x, x)
		testHessianV1EncodeDecode(t, x, x)
	})

	t.Run("should test date", func(t *testing.T) {
		ms := time.Now().UnixNano() / 1000 / 1000
		x := time.Unix(0, ms*1000*1000)
		testHessian4V2EncodeDecode(t, x, x)
		testHessian3V2EncodeDecode(t, x, x)
		testHessianV1EncodeDecode(t, x, x)
	})

	t.Run("should test float32", func(t *testing.T) {
		testHessian3V2EncodeDecode(t, float32(13.0), float64(13.0))
		testHessian4V2EncodeDecode(t, float32(13.0), float64(13.0))
		testHessianV1EncodeDecode(t, float32(13.0), float64(13.0))
	})

	t.Run("should test float64", func(t *testing.T) {
		testHessian4V2EncodeDecode(t, float64(13.0), float64(13.0))
		testHessian3V2EncodeDecode(t, float64(13.0), float64(13.0))
		testHessianV1EncodeDecode(t, float64(13.0), float64(13.0))
	})

	t.Run("should test uint8", func(t *testing.T) {
		testHessian4V2EncodeDecode(t, uint8(13), int32(13))
		testHessian3V2EncodeDecode(t, uint8(13), int32(13))
		testHessianV1EncodeDecode(t, uint8(13), int32(13))
	})

	t.Run("should test int8", func(t *testing.T) {
		testHessian4V2EncodeDecode(t, int8(13), int32(13))
		testHessian3V2EncodeDecode(t, int8(13), int32(13))
		testHessianV1EncodeDecode(t, int8(13), int32(13))
	})

	t.Run("should test uint16", func(t *testing.T) {
		testHessian4V2EncodeDecode(t, uint16(13), int32(13))
		testHessian3V2EncodeDecode(t, uint16(13), int32(13))
		testHessianV1EncodeDecode(t, uint16(13), int32(13))
	})

	t.Run("should test int16", func(t *testing.T) {
		testHessian4V2EncodeDecode(t, int16(13), int32(13))
		testHessian3V2EncodeDecode(t, int16(13), int32(13))
		testHessianV1EncodeDecode(t, int16(13), int32(13))
	})

	t.Run("should test uint32", func(t *testing.T) {
		testHessian4V2EncodeDecode(t, uint32(13), int32(13))
		testHessian3V2EncodeDecode(t, uint32(13), int32(13))
		testHessianV1EncodeDecode(t, uint32(13), int32(13))
	})

	t.Run("should test int32", func(t *testing.T) {
		testHessian4V2EncodeDecode(t, int32(13), int32(13))
		testHessian3V2EncodeDecode(t, int32(13), int32(13))
		testHessianV1EncodeDecode(t, int32(13), int32(13))
	})

	t.Run("should test uint64", func(t *testing.T) {
		for _, i := range []int{1, 2047, 262143, 2147483647, 3147483647} {
			testHessian4V2EncodeDecode(t, uint64(i), int64(i))
			testHessian3V2EncodeDecode(t, uint64(i), int64(i))
			testHessianV1EncodeDecode(t, uint64(i), int64(i))
		}
	})

	t.Run("should test int64", func(t *testing.T) {
		testHessian4V2EncodeDecode(t, int64(13), int64(13))
		testHessian3V2EncodeDecode(t, int64(13), int64(13))
		testHessianV1EncodeDecode(t, int64(13), int64(13))
	})

	t.Run("should test binary", func(t *testing.T) {
		for _, x := range [][]byte{
			[]byte("abcd"),
			[]byte("abcd😎😎😎😎😎😎😎😎😎😎😎"),
		} {
			testHessian4V2EncodeDecode(t, x, x)
			testHessian3V2EncodeDecode(t, x, x)
			testHessianV1EncodeDecode(t, x, x)
		}
	})

	t.Run("should test bool", func(t *testing.T) {
		testHessian4V2EncodeDecode(t, true, true)
		testHessian4V2EncodeDecode(t, false, false)
		testHessian3V2EncodeDecode(t, true, true)
		testHessian3V2EncodeDecode(t, false, false)
		testHessianV1EncodeDecode(t, true, true)
		testHessianV1EncodeDecode(t, false, false)
	})

	t.Run("should test nil", func(t *testing.T) {
		testHessian4V2EncodeDecode(t, nil, nil)
		testHessian3V2EncodeDecode(t, nil, nil)
		testHessianV1EncodeDecode(t, nil, nil)
	})

	t.Run("should test typedmap", func(t *testing.T) {
		x := map[interface{}]interface{}{
			"a": "1",
			"b": "1",
			"c": "1",
			"d": "1",
		}
		testHessian4V2EncodeDecode(t, x, x)
		testHessian3V2EncodeDecode(t, x, x)
		testHessianV1EncodeDecode(t, x, x)
		testEncodeToJSON(t, x, `{"a":"1","b":"1","c":"1","d":"1"}`)
	})

	t.Run("should test untypedmap", func(t *testing.T) {
		x := map[interface{}]interface{}{
			"a": "1",
			"b": "1",
			"c": "1",
			"d": "1",
		}
		z := TestTypedMap(x)
		y := NewJavaMap("test.typedmap", x)
		testHessian4V2EncodeDecode(t, z, y)
		testHessian3V2EncodeDecode(t, z, y)
		testEncodeToJSON(t, x, `{"a":"1","b":"1","c":"1","d":"1"}`)
	})

	t.Run("should test list", func(t *testing.T) {
		x := []interface{}{
			"a", "b", int32(1), int64(2), nil,
		}
		testHessian4V2EncodeDecode(t, x, x)
		testHessian3V2EncodeDecode(t, x, x)
		testHessianV1EncodeDecode(t, x, x)
		testEncodeToJSON(t, x, `["a","b",1,2,null]`)
	})

	t.Run("should test object", func(t *testing.T) {
		x := DemoCar{
			A:       "cccc",
			C:       "666",
			B:       "haha",
			Color:   "333",
			Mileage: 3762323,
		}
		y := &JavaObject{
			class: "hessian.demo.Car",
			names: []string{"a", "c", "b", "model", "color", "mileage"},
			values: []interface{}{
				"cccc",
				"666",
				"haha",
				"",
				"333",
				int32(3762323),
			},
		}
		testHessian4V2EncodeDecode(t, x, y)
		testHessian3V2EncodeDecode(t, x, y)
		testHessianV1EncodeDecode(t, x, y)
		testEncodeToJSON(t, x, `{"$":{"a":"cccc","b":"haha","c":"666","color":"333","mileage":3762323,"model":""},"$class":"hessian.demo.Car"}`)
	})

	t.Run("test complex struct", func(t *testing.T) {
		Bool := false
		Boolp := &Bool
		Boolpp := &Boolp
		Int32 := int32(13234)
		Int32p := &Int32
		Int32pp := &Int32p
		Int64 := int64(16434)
		Int64p := &Int64
		Int64pp := &Int64p
		Float64 := float64(16434.2)
		Float64p := &Float64
		Float64pp := &Float64p

		String := "1.2.3😎😎😎😎😎😎😎😎😎😎😎"
		Stringp := &String
		Stringpp := &Stringp

		ms := time.Now().UnixNano() / 1000 / 1000
		Time := time.Unix(0, ms*1000*1000)
		Timep := &Time
		Timepp := &Timep

		Binary := []byte("🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻🐻")
		Binaryp := &Binary
		Binarypp := &Binaryp

		Foo := structfoo{
			A: 132323,
		}
		Foop := &Foo
		Foopp := &Foop

		Map := map[string]string{
			"a": "b",
			"4": "5",
		}
		Mapp := &Map
		Mappp := &Mapp

		UntypedSlice := []string{"a", "1"}

		Slice := []interface{}{
			"a",
			1,
			"c",
		}
		Slicep := &Slice
		Slicepp := &Slicep

		c := complex{
			Bool:         Bool,
			Boolp:        Boolp,
			Boolpp:       Boolpp,
			Int32:        Int32,
			Int32p:       Int32p,
			Int32pp:      Int32pp,
			Int64:        Int64,
			Int64p:       Int64p,
			Int64pp:      Int64pp,
			Float64:      Float64,
			Float64p:     Float64p,
			Float64pp:    Float64pp,
			String:       String,
			Stringp:      Stringp,
			Stringpp:     Stringpp,
			Time:         Time,
			Timep:        Timep,
			Timepp:       Timepp,
			Binary:       Binary,
			Binaryp:      Binaryp,
			Binarypp:     Binarypp,
			Foo:          Foo,
			Foop:         Foop,
			Foopp:        Foopp,
			Map:          Map,
			Mapp:         Mapp,
			Mappp:        Mappp,
			Slice:        Slice,
			Slicep:       Slicep,
			Slicepp:      Slicepp,
			UntypedSlice: UntypedSlice,
		}

		testHessian3V2EncodeDecodeByJSONEqual(t, c, c)
		testEncodeDecodeByJSONEqual(t, c, c)

		y := []interface{}{c, c, nil, "abcd", int32(0), false, []byte("abc")}
		testEncodeDecodeByJSONEqual(t, y, y)

		testEncodeDecodeByJSONEqual(t,
			[]interface{}{c, c, "1", nil},
			[]interface{}{c, c, "1", nil},
		)
	})
}

func testHessian4V2EncodeDecode(t *testing.T, x, y interface{}) {
	xb, err := EncodeHessian4V2(NewEncodeContext().SetTracer(NewDummyTracer()), x)
	require.Nil(t, err)
	cr := NewClassRegistry()
	cr.RegisterJavaClass(javaobject.JavaLangStackTraceElements{})
	cr.RegisterJavaClass(complex{})
	cr.RegisterJavaClass(ConnectionResponse{})
	cr.RegisterJavaClass(HessianConnectionRequest{})
	z, err := DecodeHessian4V2(NewDecodeContext().SetTracer(NewDummyTracer()).
		SetClassRegistry(cr),
		bufio.NewReader(bytes.NewReader(xb)))
	require.Nil(t, err)
	require.Equal(t, y, z)
}

func testHessian3V2EncodeDecode(t *testing.T, x, y interface{}) {
	xb, err := EncodeHessian3V2(NewEncodeContext().SetTracer(NewStdoutTracer()), x)
	require.Nil(t, err)
	cr := NewClassRegistry()
	cr.RegisterJavaClass(javaobject.JavaLangStackTraceElements{})
	cr.RegisterJavaClass(complex{})
	cr.RegisterJavaClass(ConnectionResponse{})
	cr.RegisterJavaClass(HessianConnectionRequest{})
	z, err := DecodeHessian3V2(NewDecodeContext().SetTracer(NewStdoutTracer()).
		SetClassRegistry(cr),
		bufio.NewReader(bytes.NewReader(xb)))
	require.Nil(t, err)
	require.Equal(t, y, z)
}

func testHessianV1EncodeDecode(t *testing.T, x, y interface{}) {
	xb, err := EncodeHessianV1(NewEncodeContext().SetTracer(NewStdoutTracer()), x)
	require.Nil(t, err)
	cr := NewClassRegistry()
	cr.RegisterJavaClass(javaobject.JavaLangStackTraceElements{})
	cr.RegisterJavaClass(complex{})
	cr.RegisterJavaClass(ConnectionResponse{})
	cr.RegisterJavaClass(HessianConnectionRequest{})
	z, err := DecodeHessianV1(NewDecodeContext().SetTracer(NewStdoutTracer()).
		SetClassRegistry(cr),
		bufio.NewReader(bytes.NewReader(xb)))
	require.Nil(t, err)
	require.Equal(t, y, z)
}

func testHessian3V2EncodeDecodeByJSONEqual(t *testing.T, x, y interface{}) {
	cr := NewClassRegistry()
	cr.RegisterJavaClass(complex{})
	ectx := NewEncodeContext().SetTracer(NewDummyTracer())
	xb, err := EncodeHessian3V2(ectx, x)
	require.Nil(t, err)
	dctx := NewDecodeContext()
	z, err := DecodeHessian3V2(dctx.
		SetClassRegistry(cr).
		SetTracer(NewDummyTracer()),
		bufio.NewReader(bytes.NewReader(xb)))
	require.Nil(t, err)

	g1, err := json.MarshalIndent(z, "", "    ")
	require.Nil(t, err)
	g2, err := json.MarshalIndent(y, "", "    ")
	require.Nil(t, err)
	require.Equal(t, string(g2), string(g1))
}

func testEncodeDecodeByJSONEqual(t *testing.T, x, y interface{}) {
	ectx := NewEncodeContext().SetTracer(NewDummyTracer())
	xb, err := EncodeHessian4V2(ectx, x)
	require.Nil(t, err)
	cr := NewClassRegistry()
	cr.RegisterJavaClass(complex{})
	dctx := NewDecodeContext()
	z, err := DecodeHessian4V2(dctx.
		SetClassRegistry(cr).
		SetTracer(NewDummyTracer()),
		bufio.NewReader(bytes.NewReader(xb)))

	require.Nil(t, err)
	g1, err := json.MarshalIndent(z, "", "    ")
	require.Nil(t, err)
	g2, err := json.MarshalIndent(y, "", "    ")
	require.Nil(t, err)
	require.Equal(t, string(g2), string(g1))
}

func TestEncode1CodeGolden(t *testing.T) {
	connectionRequestBytes, _ := hex.DecodeString("4d74002a636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e526571756573745300036374784d740039636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e526571756573742452657175657374436f6e7465787453000269644c000000000000008653000674686973243052000000007a7a636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e4865617274426561744d74002c636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e486561727442656174530009636c69656e7455726c5300476261722d706f6f6c2e746573742e616c697061792e6e65743a31323230303f703d3133265f53455249414c495a45545950453d3126763d342e30266170705f6e616d653d6261727a")
	o := NewDecodeContext().SetVersion(HessianV1)
	y, err := DecodeHessianV1(o, bufio.NewReader(bytes.NewReader(connectionRequestBytes)))
	require.Nil(t, err)
	t.Log(y)
	yC := y.(*JavaObject)
	hr := ConnectionRequest{}
	require.Equal(t, yC.class, hr.GetJavaClassName())
}

func TestEncodeAndDecodeHessian1Object(t *testing.T) {
	hr := ConnectionRequest{}
	cr := NewClassRegistry()
	cr.RegisterJavaClass(hr)
	hc := RequestContext{}
	cr.RegisterJavaClass(hc)
	cr.RegisterJavaClass(javaobject.JavaUtilConcurrentAtomicLong{0})
	o := NewEncodeContext().SetVersion(HessianV1).SetTracer(NewDummyTracer())

	crq := &ConnectionRequest{
		Ctx:                nil,
		FromAppKey:         "",
		ToAppKey:           "",
		EncryptedToken:     nil,
		ApplicationRequest: nil,
	}

	ctx := &RequestContext{This: crq, Id: 134}
	crq.Ctx = ctx

	encoder := NewEncoder(o)
	encoder.Encode(crq)

	do := NewDecodeContext().SetVersion(HessianV1).SetTracer(NewDummyTracer()).SetClassRegistry(cr)
	i := encoder.Bytes()

	decoder := NewDecoder(do, bufio.NewReader(bytes.NewReader(i)))
	y, err := decoder.Decode()
	if err != nil {
		t.Fatal(err)
	} else {
		yC := y.(*ConnectionRequest)
		require.Equal(t, crq.Ctx.Id, yC.Ctx.Id)
	}
}

type L struct {
	A M `hessian:"a"`
	B M `hessian:"b"`
}

func (l L) GetJavaClassName() string {
	return "com.alipay.hessian.dump.L"
}

type M struct {
	M string `hessian:"m"`
}

func (m M) GetJavaClassName() string {
	return "com.alipay.hessian.dump.M"
}

func TestEq(t *testing.T) {
	i := M{M: "a"}
	i2 := M{M: "a"}
	l := L{A: i2, B: i}

	m := M{M: "a"}
	cr := NewClassRegistry()
	l1 := L{}
	cr.RegisterJavaClass(l1)
	cr.RegisterJavaClass(m)

	data, _ := hex.DecodeString("4FA9636F6D2E616C697061792E6865737369616E2E64756D702E4C92016101626F904FA9636F6D2E616C697061792E6865737369616E2E64756D702E4D91016D6F9101616F910161")
	do := NewDecodeContext().SetVersion(Hessian3xV2).SetTracer(NewDummyTracer()).SetClassRegistry(cr)
	decoder := NewDecoder(do, bufio.NewReader(bytes.NewReader(data)))
	obj, err := decoder.Decode()
	_ = obj
	require.Nil(t, err)

	o := NewEncodeContext().SetVersion(Hessian3xV2).SetTracer(NewDummyTracer())
	encoder := NewEncoder(o)

	encoder.Encode(l)
	require.Equal(t, data, encoder.Bytes())
}
