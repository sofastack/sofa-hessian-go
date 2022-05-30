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
	"fmt"
	"testing"

	"github.com/sofastack/sofa-hessian-go/javaobject"
	"github.com/stretchr/testify/require"
)

func TestDecodeObject(t *testing.T) {
	t.Run("should read enum Color", func(t *testing.T) {
		testDecodeObject(t,
			readFile(t, "testdata/enum/red.bin"),
			&JavaObject{
				class:  "hessian.Main$Color",
				names:  []string{"name"},
				values: []interface{}{"RED"},
			},
		)

		testDecodeObject(t,
			readFile(t, "testdata/enum/green.bin"),
			&JavaObject{
				class:  "hessian.Main$Color",
				names:  []string{"name"},
				values: []interface{}{"GREEN"},
			},
		)

		testDecodeObject(t,
			readFile(t, "testdata/enum/blue.bin"),
			&JavaObject{
				class:  "hessian.Main$Color",
				names:  []string{"name"},
				values: []interface{}{"BLUE"},
			},
		)
	})

	t.Run("should decode and encode to ConnectionRequest", func(t *testing.T) {
		ctx := &HessianConnectionRequestContext{
			Id:   101,
			This: &HessianConnectionRequest{},
		}
		ctxp := &ctx
		ctxpp := &ctxp
		ctxppp := &ctxpp
		ctxpppp := &ctxppp
		hr := &HessianConnectionRequest{Ctx: ctxpppp}
		ctx.This = hr
		chr := &HessianConnectionRequest{}

		testDecodeToObject(t,
			readFile(t, "testdata/object/ConnectionRequest.bin"),
			chr,
			hr,
		)
	})

	t.Run("should decode and encode ConnectionRequest", func(t *testing.T) {
		ctx := &HessianConnectionRequestContext{
			Id:   101,
			This: &HessianConnectionRequest{},
		}
		ctxp := &ctx
		ctxpp := &ctxp
		ctxppp := &ctxpp
		ctxpppp := &ctxppp
		hr := &HessianConnectionRequest{Ctx: ctxpppp}
		ctx.This = hr

		testDecodeObject(t,
			readFile(t, "testdata/object/ConnectionRequest.bin"),
			hr,
		)
	})

	t.Run("java.util.concurrent.atomic.AtomicLong", func(t *testing.T) {
		testDecodeObject(t,
			readFile(t, "testdata/object/AtomicLong0.bin"),
			&javaobject.JavaUtilConcurrentAtomicLong{0},
		)

		testDecodeObject(t,
			readFile(t, "testdata/object/AtomicLong1.bin"),
			&javaobject.JavaUtilConcurrentAtomicLong{1},
		)
	})

	t.Run("should encode and decode array field", func(t *testing.T) {
		type A struct {
			A []int32
		}
		dst, err := EncodeObjectHessian4V2(&EncodeContext{
			classrefs:  NewEncodeClassrefs(),
			objectrefs: NewEncodeObjectrefs(),
		}, A{
			A: []int32{4, 5, 6},
		})
		require.Nil(t, err)
		obj, err := DecodeObjectHessian4V2(NewDecodeContext().SetTracer(NewDummyTracer()), bufio.NewReader(bytes.NewReader(dst)))
		require.Nil(t, err)
		require.Equal(t, obj, &JavaObject{
			class:  "A",
			names:  []string{"A"},
			values: []interface{}{[]interface{}{int32(4), int32(5), int32(6)}},
		})
	})
}

type HessianConnectionRequestContext struct {
	Id   int32                     `hessian:"id"`
	This *HessianConnectionRequest `hessian:"this$0"`
}

func (ctx HessianConnectionRequestContext) GetJavaClassName() string {
	return "hessian.ConnectionRequest$RequestContext"
}

type HessianConnectionRequest struct {
	Ctx *****HessianConnectionRequestContext `hessian:"ctx"`
}

func (ctx HessianConnectionRequest) GetJavaClassName() string {
	return "hessian.ConnectionRequest"
}

func testDecodeObject(t *testing.T, b []byte, x interface{}) {
	hr := HessianConnectionRequest{}
	cr := NewClassRegistry()
	cr.RegisterJavaClass(hr)
	cr.RegisterJavaClass(javaobject.JavaUtilConcurrentAtomicLong{0})
	o := NewDecodeContext().SetTracer(NewDummyTracer()).SetClassRegistry(cr)
	y, err := DecodeObjectHessian4V2(o, bufio.NewReader(bytes.NewReader(b)))
	require.Nil(t, err)
	require.Equal(t, x, y)
}

func testDecodeToObject(t *testing.T, b []byte, x, y interface{}) {
	hr := HessianConnectionRequest{}
	cr := NewClassRegistry()
	cr.RegisterJavaClass(hr)
	cr.RegisterJavaClass(javaobject.JavaUtilConcurrentAtomicLong{0})
	o := NewDecodeContext().SetTracer(NewDummyTracer()).SetClassRegistry(cr)
	err := DecodeObjectToHessian4V2(o, bufio.NewReader(bytes.NewReader(b)), x)
	require.Nil(t, err)
	require.Equal(t, x, y)
}

func TestDecodeHessian1Object(t *testing.T) {
	t.Run("should read hessian1 object", func(t *testing.T) {
		cr := &ConnectionRequest{
			Ctx:                nil,
			FromAppKey:         "",
			ToAppKey:           "",
			EncryptedToken:     nil,
			ApplicationRequest: nil,
		}

		ctx := &RequestContext{This: cr, Id: 134}
		cr.Ctx = ctx
		connectionRequestBytes, _ := hex.DecodeString("4d74002a636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e526571756573745300036374784d740039636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e526571756573742452657175657374436f6e7465787453000269644c000000000000008653000674686973243052000000007a7a636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e4865617274426561744d74002c636f6d2e74616f62616f2e72656d6f74696e672e696d706c2e436f6e6e656374696f6e486561727442656174530009636c69656e7455726c5300476261722d706f6f6c2e746573742e616c697061792e6e65743a31323230303f703d3133265f53455249414c495a45545950453d3126763d342e30266170705f6e616d653d6261727a")
		fmt.Printf(hex.Dump(connectionRequestBytes))

		testDecodeHessian1Object(t, connectionRequestBytes, cr)
	})
}

func testDecodeHessian1Object(t *testing.T, b []byte, x *ConnectionRequest) {
	hr := ConnectionRequest{}
	cr := NewClassRegistry()
	cr.RegisterJavaClass(hr)
	hc := RequestContext{}
	cr.RegisterJavaClass(hc)
	cr.RegisterJavaClass(javaobject.JavaUtilConcurrentAtomicLong{0})
	o := NewDecodeContext().SetVersion(HessianV1).SetTracer(NewDummyTracer()).SetClassRegistry(cr)
	y, err := DecodeHessianV1(o, bufio.NewReader(bytes.NewReader(b)))

	yC := y.(*ConnectionRequest)

	require.Nil(t, err)
	require.Equal(t, x.Ctx.Id, yC.Ctx.Id)
}
