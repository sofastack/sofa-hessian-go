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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeToJSON(t *testing.T) {
	t.Run("should test gold data", func(t *testing.T) {
		jo := &JavaObject{
			class: "com.alipay.remoting.rpc.exception.RpcServerException",
			names: []string{
				"detailMessage",
				"cause",
				"stackTrace",
				"suppressedExceptions",
			},
			values: []interface{}{
				"[Server]OriginErrorMsg: com.alipay.remoting.exception.DeserializationException: Content of request is null. AdditionalErrorMsg: null",
				&JavaObject{},
				&JavaList{
					class: "[java.lang.StackTraceElement",
					value: []interface{}{
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{
								"declaringClass",
								"methodName",
								"fileName",
								"lineNumber",
							},
							values: []interface{}{
								"com.alipay.sofa.rpc.codec.bolt.SofaRpcSerialization",
								"deserializeContent",
								"SofaRpcSerialization.java",
								241,
							},
						},
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{
								"declaringClass",
								"methodName",
								"fileName",
								"lineNumber",
							},
							values: []interface{}{
								"com.alipay.remoting.rpc.protocol.RpcRequestCommand",
								"deserializeContent",
								"RpcRequestCommand.java",
								144,
							},
						},
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{
								"declaringClass",
								"methodName",
								"fileName",
								"lineNumber",
							},
							values: []interface{}{
								"com.alipay.remoting.rpc.RpcCommand",
								"deserialize",
								"RpcCommand.java",
								117,
							},
						},
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{
								"declaringClass",
								"methodName",
								"fileName",
								"lineNumber",
							},
							values: []interface{}{
								"com.alipay.remoting.rpc.RpcCommand",
								"deserialize",
								"RpcCommand.java",
								138,
							},
						},
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{
								"declaringClass",
								"methodName",
								"fileName",
								"lineNumber",
							},
							values: []interface{}{
								"com.alipay.remoting.rpc.protocol.RpcRequestProcessor",
								"deserializeRequestCommand",
								"RpcRequestProcessor.java",
								267,
							},
						},
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{
								"declaringClass",
								"methodName",
								"fileName",
								"lineNumber",
							},
							values: []interface{}{
								"com.alipay.remoting.rpc.protocol.RpcRequestProcessor",
								"doProcess",
								"RpcRequestProcessor.java",
								142,
							},
						},
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{
								"declaringClass",
								"methodName",
								"fileName",
								"lineNumber",
							},
							values: []interface{}{
								"com.alipay.remoting.rpc.protocol.RpcRequestProcessor$ProcessTask",
								"run",
								"RpcRequestProcessor.java",
								366,
							},
						},
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{
								"declaringClass",
								"methodName",
								"fileName",
								"lineNumber",
							},
							values: []interface{}{
								"java.util.concurrent.ThreadPoolExecutor",
								"runWorker",
								"ThreadPoolExecutor.java",
								1149,
							},
						},
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{
								"declaringClass",
								"methodName",
								"fileName",
								"lineNumber",
							},
							values: []interface{}{
								"java.util.concurrent.ThreadPoolExecutor$Worker",
								"run",
								"ThreadPoolExecutor.java",
								624,
							},
						},
						&JavaObject{
							class: "java.lang.StackTraceElement",
							names: []string{
								"declaringClass",
								"methodName",
								"fileName",
								"lineNumber",
							},
							values: []interface{}{
								"java.lang.Thread",
								"run",
								"Thread.java",
								748,
							},
						},
					},
				},
				&JavaList{
					class: "java.util.Collections$UnmodifiableRandomAccessList",
					value: []interface{}{},
				},
			},
		}
		jo.values[1] = jo
		testEncodeToJSON(t, jo, `{"$":{"cause":{"$":"JavaObject","$class":"CyclicReference"},"detailMessage":"[Server]OriginErrorMsg: com.alipay.remoting.exception.DeserializationException: Content of request is null. AdditionalErrorMsg: null","stackTrace":{"$":[{"$":{"declaringClass":"com.alipay.sofa.rpc.codec.bolt.SofaRpcSerialization","fileName":"SofaRpcSerialization.java","lineNumber":241,"methodName":"deserializeContent"},"$class":"java.lang.StackTraceElement"},{"$":{"declaringClass":"com.alipay.remoting.rpc.protocol.RpcRequestCommand","fileName":"RpcRequestCommand.java","lineNumber":144,"methodName":"deserializeContent"},"$class":"java.lang.StackTraceElement"},{"$":{"declaringClass":"com.alipay.remoting.rpc.RpcCommand","fileName":"RpcCommand.java","lineNumber":117,"methodName":"deserialize"},"$class":"java.lang.StackTraceElement"},{"$":{"declaringClass":"com.alipay.remoting.rpc.RpcCommand","fileName":"RpcCommand.java","lineNumber":138,"methodName":"deserialize"},"$class":"java.lang.StackTraceElement"},{"$":{"declaringClass":"com.alipay.remoting.rpc.protocol.RpcRequestProcessor","fileName":"RpcRequestProcessor.java","lineNumber":267,"methodName":"deserializeRequestCommand"},"$class":"java.lang.StackTraceElement"},{"$":{"declaringClass":"com.alipay.remoting.rpc.protocol.RpcRequestProcessor","fileName":"RpcRequestProcessor.java","lineNumber":142,"methodName":"doProcess"},"$class":"java.lang.StackTraceElement"},{"$":{"declaringClass":"com.alipay.remoting.rpc.protocol.RpcRequestProcessor$ProcessTask","fileName":"RpcRequestProcessor.java","lineNumber":366,"methodName":"run"},"$class":"java.lang.StackTraceElement"},{"$":{"declaringClass":"java.util.concurrent.ThreadPoolExecutor","fileName":"ThreadPoolExecutor.java","lineNumber":1149,"methodName":"runWorker"},"$class":"java.lang.StackTraceElement"},{"$":{"declaringClass":"java.util.concurrent.ThreadPoolExecutor$Worker","fileName":"ThreadPoolExecutor.java","lineNumber":624,"methodName":"run"},"$class":"java.lang.StackTraceElement"},{"$":{"declaringClass":"java.lang.Thread","fileName":"Thread.java","lineNumber":748,"methodName":"run"},"$class":"java.lang.StackTraceElement"}],"$class":"[java.lang.StackTraceElement"},"suppressedExceptions":{"$":[],"$class":"java.util.Collections$UnmodifiableRandomAccessList"}},"$class":"com.alipay.remoting.rpc.exception.RpcServerException"}`)
	})
	t.Run("should truncate encode circular data", func(t *testing.T) {
		type Circular struct {
			A    int
			This *Circular
		}
		c := &Circular{
			A: 123,
		}
		c.This = c
		testEncodeToJSON(t, c, `{"$":{"A":123,"This":{"$":"*sofahessian.Circular","$class":"CyclicReference"}},"$class":"Circular"}`)
	})

	t.Run("should encode nested map", func(t *testing.T) {
		x := map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"b": "c",
				},
			},
		}
		y := map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"b": "c",
					"d": map[string]interface{}{
						"1": 1,
						"2": 2.2,
						"3": "cc",
					},
				},
			},
		}
		testEncodeToJSON(t, x, `{"a":{"b":{"b":"c"}}}`)
		testEncodeToJSON(t, y, `{"a":{"b":{"b":"c","d":{"1":1,"2":2.2,"3":"cc"}}}}`)
	})

	t.Run("shoud encode int8 to json", func(t *testing.T) {
		for _, i := range []int8{-128, -3, -4, 1, 3, 5, 7, 9} {
			testEncodeToJSON(t, int8(i), fmt.Sprintf("%d", i))
		}
	})

	t.Run("shoud encode uint8 to json", func(t *testing.T) {
		for _, i := range []uint8{1, 3, 5, 7, 9} {
			testEncodeToJSON(t, uint8(i), fmt.Sprintf("%d", i))
		}
	})

	t.Run("shoud encode int16 to json", func(t *testing.T) {
		for _, i := range []int16{-128, -3, -4, 1, 3, 5, 7, 9} {
			testEncodeToJSON(t, int16(i), fmt.Sprintf("%d", i))
		}
	})

	t.Run("shoud encode uint16 to json", func(t *testing.T) {
		for _, i := range []uint16{1, 3, 5, 7, 9} {
			testEncodeToJSON(t, uint16(i), fmt.Sprintf("%d", i))
		}
	})

	t.Run("shoud encode uint32 to json", func(t *testing.T) {
		for _, i := range []uint32{1, 3, 5, 7, 9, 128, 256, 512, 1024} {
			testEncodeToJSON(t, uint32(i), fmt.Sprintf("%d", i))
		}
	})

	t.Run("shoud encode int32 to json", func(t *testing.T) {
		for _, i := range []int32{-128, -3, -4, 1, 3, 5, 7, 9, 128, 256, 512, 1024} {
			testEncodeToJSON(t, int32(i), fmt.Sprintf("%d", i))
		}
	})

	t.Run("shoud encode int64 to json", func(t *testing.T) {
		for _, i := range []int64{-128, -3, -4, 1, 3, 5, 7, 9, 128, 256, 512, 1024} {
			testEncodeToJSON(t, int64(i), fmt.Sprintf("%d", i))
		}
	})

	t.Run("shoud encode int to json", func(t *testing.T) {
		for _, i := range []int{-128, -3, -4, 1, 3, 5, 7, 9, 128, 256, 512, 1024} {
			testEncodeToJSON(t, int(i), fmt.Sprintf("%d", i))
		}
	})

	t.Run("shoud encode int64 to json", func(t *testing.T) {
		for _, i := range []uint{1, 3, 5, 7, 9, 128, 256, 512, 1024} {
			testEncodeToJSON(t, uint(i), fmt.Sprintf("%d", i))
		}
	})

	t.Run("shoud encode bool to json", func(t *testing.T) {
		for _, i := range []bool{true, false} {
			testEncodeToJSON(t, i, fmt.Sprintf("%t", i))
		}
	})

	t.Run("shoud encode nil to json", func(t *testing.T) {
		testEncodeToJSON(t, nil, "null")
	})

	t.Run("shoud encode string to json", func(t *testing.T) {
		testEncodeToJSON(t, `ab"cd`, `"ab\"cd"`)
	})

	t.Run("shoud encode []byte to json", func(t *testing.T) {
		testEncodeToJSON(t, []byte("abcd"), "\"YWJjZA==\"")
	})

	t.Run("should encode map[string]interface{} to json", func(t *testing.T) {
		testEncodeToJSON(t, map[string]interface{}{
			"a": 1,
			"b": "cd",
		}, `{"a":1,"b":"cd"}`)
	})

	t.Run("should encode []interface{} to json", func(t *testing.T) {
		testEncodeToJSON(t, []interface{}{
			1,
			"2",
			3.2,
		}, `[1,"2",3.2]`)
	})

	t.Run("should encode javaobject to json", func(t *testing.T) {
		testEncodeToJSON(t, &JavaObject{
			class:  "test",
			names:  []string{"a", "b"},
			values: []interface{}{1, "ab"},
		}, `{"$":{"a":1,"b":"ab"},"$class":"test"}`)
	})

	t.Run("should encode javamap to json", func(t *testing.T) {
		testEncodeToJSON(t, &JavaMap{
			class: "test",
			m: map[interface{}]interface{}{
				"a": 1,
				"b": "ab",
			},
		}, `{"$":{"a":1,"b":"ab"},"$class":"test"}`)
	})

	t.Run("should encode javalist to json", func(t *testing.T) {
		testEncodeToJSON(t, &JavaList{
			class: "test",
			value: []interface{}{1, 2, 3, 4, "abcd", "defg", 2.3},
		}, `{"$":[1,2,3,4,"abcd","defg",2.3],"$class":"test"}`)
	})

	t.Run("should encode struct to json", func(t *testing.T) {
		type Foo struct {
			A string
		}

		testEncodeToJSON(t, &Foo{
			A: "abcd",
		}, `{"$":{"A":"abcd"},"$class":"Foo"}`)

		testEncodeToJSON(t, &DemoCar{
			A:       "a",
			C:       "c",
			B:       "b",
			Model:   "model 1",
			Color:   "aquamarine",
			Mileage: 65536,
		}, `{"$":{"a":"a","b":"b","c":"c","color":"aquamarine","mileage":65536,"model":"model 1"},"$class":"hessian.demo.Car"}`)
	})
}

func TestEncodeJSONToHessian(t *testing.T) {
	t.Run("should avoid infinite recursive", func(t *testing.T) {
		json := `{"a":{"a":{"a":{"a":1}}}}`
		ectx := NewEncodeContext().SetTracer(NewDummyTracer())
		encoder := NewEncoder(ectx)
		jctx := NewJSONEncodeContext().SetTracer(NewDummyTracer()).SetMaxDepth(2)
		err := encoder.EncodeJSONString(jctx, json)
		require.Equal(t, err, ErrEncodeTooManyNestedJSON)
	})

	t.Run("should encode number to int64", func(t *testing.T) {
		testEncodeJSONToHessian(t, "1", int64(1))
	})

	t.Run("should encode number to float64", func(t *testing.T) {
		testEncodeJSONToHessian(t, "1.2", float64(1.2))
	})

	t.Run("should encode string to string", func(t *testing.T) {
		testEncodeJSONToHessian(t, `"abcdefg"`, "abcdefg")
	})

	t.Run("should encode true to true", func(t *testing.T) {
		testEncodeJSONToHessian(t, `true`, true)
	})

	t.Run("should encode false to false", func(t *testing.T) {
		testEncodeJSONToHessian(t, `false`, false)
	})

	t.Run("should encode null to nil", func(t *testing.T) {
		testEncodeJSONToHessian(t, `null`, nil)
	})

	t.Run("should encode [1,2,3]", func(t *testing.T) {
		testEncodeJSONToHessian(t, `[1,2,3]`, []interface{}{int64(1), int64(2), int64(3)})
	})

	t.Run(`should encode map`, func(t *testing.T) {
		testEncodeJSONToHessian(t, `{"a": "b"}`, map[interface{}]interface{}{"a": "b"})
	})

	t.Run("should encode $class to object", func(t *testing.T) {
		testEncodeJSONToHessian(t,
			string(readFile(t, "testdata/json/request.json")),
			&JavaObject{
				class: "com.alipay.sofa.rpc.core.request.SofaRequest",
				names: []string{"targetAppName", "methodName", "targetServiceUniqueName", "requestProps", "methodArgSigs"},
				values: []interface{}{
					"", "sayHello", "HelloService:1.0",
					map[interface{}]interface{}{
						"protocol": "bolt",
						"rpc_trace_context": map[interface{}]interface{}{
							"samp":           "false",
							"sofaCallerApp":  "",
							"sofaCallerIdc":  "",
							"sofaCallerIp":   "",
							"sofaCallerZone": "",
							"sofaPenAttrs":   "",
							"sofaRpcId":      "0",
							"sofaTraceId":    "0a0fe8631571046378758100186220",
							"sysPenAttrs":    "",
						},
					},
					[]interface{}{"java.lang.String", "long"},
				},
			})
	})
}

func testEncodeToJSON(t *testing.T, obj interface{}, expect string) {
	jctx := NewJSONEncodeContext().
		SetMaxPtrCycles(1).
		SetTracer(NewDummyTracer()).
		SetLessFunc(func(keyi, keyj, valuei, valuej interface{}) bool {
			i := keyi.(string)
			j := keyj.(string)
			return i <= j
		})
	dst, err := EncodeToJSON(jctx, nil, obj)
	require.Nil(t, err)
	require.Equal(t, expect, string(dst))
}

func testEncodeJSONToHessian(t *testing.T, json string, x interface{}) {
	ectx := NewEncodeContext().SetTracer(NewDummyTracer())
	dctx := NewDecodeContext().SetTracer(NewDummyTracer())
	encoder := NewEncoder(ectx)
	jctx := NewJSONEncodeContext().SetTracer(NewDummyTracer())
	err := encoder.EncodeJSONString(jctx, json)
	require.Nil(t, err)

	decoder := NewDecoder(dctx, bufio.NewReader(bytes.NewReader(encoder.Bytes())))
	y, err := decoder.Decode()
	require.Nil(t, err)
	require.Equal(t, x, y)
}
