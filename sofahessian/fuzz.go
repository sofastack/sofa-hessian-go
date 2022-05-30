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

//go:build fuzz
// +build fuzz

package sofahessian

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"reflect"
)

func Fuzz(data []byte) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(hex.EncodeToString(data))
			panic(r)
		}
	}()
	dctx := NewDecodeContext()
	ectx := NewEncodeContext()
	br := bufio.NewReader(bytes.NewReader(data))
	v, err := DecodeHessian4V2(dctx, br)
	if err != nil {
		return 0
	}

	if fv, ok := v.(float64); ok && math.IsNaN(fv) {
		return 0
	}

	dst, err := EncodeHessian4V2(ectx, v)
	if err != nil {
		fmt.Println(hex.EncodeToString(data))
		panic(err)
	}

	vv, err := DecodeHessian4V2(dctx, bufio.NewReader(bytes.NewReader(dst)))
	if err != nil {
		if err == ErrDecodeMapUnhashable {
			return 0
		}

		fmt.Println(hex.EncodeToString(data))
		panic(err)
	}

	if mv, ok := v.(map[interface{}]interface{}); ok {
		if ov, ok := vv.(map[interface{}]interface{}); ok {
			_ = mv
			_ = ov
			return 0
		}
		log.Fatal("expect map[interface{}]interface{}")
	}

	if !reflect.DeepEqual(v, vv) {
		fmt.Println(hex.EncodeToString(data))
		fmt.Println(hex.EncodeToString(dst))
		panic("not equal")
	}

	return 1
}
