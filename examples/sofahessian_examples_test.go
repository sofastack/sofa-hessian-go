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

package examples

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/sofastack/sofa-hessian-go/sofahessian"
)

func ExampleEncodeHessian() {
	ectx := sofahessian.NewEncodeContext()
	ectx.
		SetVersion(sofahessian.Hessian3xV2)

	encoder := sofahessian.NewEncoder(ectx)
	encoder.EncodeInt32(1)
	encoder.EncodeInt64(1)
	encoder.EncodeFloat64(1.23)
	encoder.EncodeString("example")
	encoder.EncodeBool(true)
	encoder.EncodeNil()
	encoder.EncodeBinary([]byte{0x01, 0x02})
	encoder.Encode(2)
	encoder.Encode(3)
	encoder.Encode(5)
	d := encoder.Bytes()
	fmt.Println(hex.EncodeToString(d))
	// Output: 91e1443ff3ae147ae147ae076578616d706c65544e220102e2e3e5
}

func ExampleDecodeHessian() {
	dctx := sofahessian.NewDecodeContext()
	dctx.
		SetVersion(sofahessian.Hessian3xV2)
	d, err := hex.DecodeString("91e1443ff3ae147ae147ae076578616d706c65544e220102e2e3e5")
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(bytes.NewReader(d))
	decoder := sofahessian.NewDecoder(dctx, reader)
	i32, err := decoder.DecodeInt32()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(i32)

	i64, err := decoder.DecodeInt64()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(i64)

	f64, err := decoder.DecodeFloat64()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(f64)
	// Output: 1
	// 1
	// 1.23
}

func ExampleJSON2Hessian() {
	ectx := sofahessian.NewEncodeContext()
	ectx.
		SetVersion(sofahessian.Hessian3xV2)
	encoder := sofahessian.NewEncoder(ectx)
	jctx := sofahessian.NewJSONEncodeContext()
	data := []byte(`{"$class": "com.alipay.sofa.rpc.core.request.SofaRequest", "$": {"1": "2"}}`)
	if err := encoder.EncodeStreamingJSONBytes(jctx, data); err != nil {
		log.Fatal(err)
	}
	dctx := sofahessian.NewDecodeContext()
	dctx.
		SetVersion(sofahessian.Hessian3xV2)

	reader := bufio.NewReader(bytes.NewReader(encoder.Bytes()))
	decoder := sofahessian.NewDecoder(dctx, reader)
	o, err := decoder.Decode()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(o)
	// Output: &{com.alipay.sofa.rpc.core.request.SofaRequest [1] [2]}
}
