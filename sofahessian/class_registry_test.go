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
	"reflect"
	"testing"

	"github.com/sofastack/sofa-hessian-go/javaobject"
)

type MyStruct struct {
}

type MySlice struct {
	P []int
}

func TestJavaObject(t *testing.T) {
	var gr ClassRegistry
	gr.RegisterJavaClass(&javaobject.TBRemotingConnectionRequest{})
}

func TestSafeIsNil(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name   string
		args   args
		wantRv bool
	}{
		{name: "struct", args: args{v: reflect.ValueOf(MyStruct{})}, wantRv: false},
		{name: "slice", args: args{v: reflect.ValueOf(MySlice{}.P)}, wantRv: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRv := safeIsNil(tt.args.v); gotRv != tt.wantRv {
				t.Errorf("safeIsNil() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}
