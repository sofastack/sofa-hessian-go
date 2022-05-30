# 目录

-   [概要](#概要)
-   [类型系统](#类型系统)
    -   [基础类型](#基础类型)
    -   [复合类型](#复合类型)
    -   [引用类型](#引用类型)
-   [编码规范](#编码规范)
-   [解码规范](#解码规范)
-   [JSON 转换](#json转换)
-   [API](#api)
    -   [encode](#encode)
    -   [decode](#decode)
    -   [json](#json)
-   [CLI](#cli)
    -   [install](#install)
    -   [decode](#decode)
    -   [fromjson](#fromjson)
-   [性能测试](#性能测试)
-   [TODO](#todo)

# 概要

`sofa-hessian-go` 是 [hessian 2.0/1.0 serialization protocol](http://hessian.caucho.com/doc/hessian-serialization.html) 的 Golang 实现，包括 1.0 协议以及 2.0 协议的 `java3.x` 和 `java4.x` 版本，同时还提供了 [JSON](https://json.org) 到 [hessian](http://hessian.caucho.com/doc/hessian-serialization.html) 类型系统互相转换的设计。

在此非常感谢 [node-modules/hessian.js](https://github.com/node-modules/hessian.js) 提供的 golden files(测试数据), 两者底层复用了同样的测试数据集。

# 类型系统

hessian 的类型系统由 8 种基础类型和 3 种复合类型，以及 3 种 引用类型组成。

## 基础类型

-   binary
-   bool
-   string
-   int32
-   int64
-   float64
-   null
-   date (64bit)

## 复合类型

-   list
-   map
-   object

## 引用类型

-   class reference: represents the definition of class
-   type reference: represents the name of class
-   object reference: represents the instance of object or list or map.

# 编码规范

sofa-hessian-go 遵循以下编码规范，将 go 的类型转化为 hessian 的类型。

-   uint8/int8 => int32
-   uint16/int => int32
-   uint32/int32 => int32
-   uint64/int64 => int64
-   uint/int => int64
-   bool => bool
-   string => string
-   []byte => binary
-   nil => null
-   time.Time => date
-   map[interface{}]interface{} => map
-   []interface{} => list
-   struct => object

# 解码规范

sofa-hessian-go 遵循以下解码规范，将 hessian 的类型转化为 go 的类型。

-   int32 => int32
-   int64 => int64
-   bool => bool
-   string => string
-   null => nil
-   date => time.Time
-   map

    -   untyped map => map[interface{}]interface{}
    -   typed map => \*JavaMap{class: "balaba", map[interface{}]interface{}}

-   list

    -   untyped list => []interface{}
    -   typed list => \*JavaList{class: "balaba", []interface{}}

-   object =>
    -   concrete object => go concrete struct
    -   generic object => \*JavaObject{class: "balaba", JavaObject{}}

# json 转换

JSON 的类型系统比 HESSIAN 的类型系统更精简，从理论上是可以在一定人为约束条件下做到 JSON 和 HESSIAN 的类型转换。

## JSON 到 HESSIAN 的类型转换

给定一个 HESSIAN 类型总是存在一种 JSON 类型可以等价描述出来。但是在实际实现中只有一种情况例外，即 HESSIAN 可以以较小的代价描述循环引用的数据结构，JSON 虽然可以描述，但通常导致的结果就是无限递归导致爆栈。 JSON 天然无法处理循环引用的数据结构，不过 sofa-hessian-go 在实际的实现中以较小的代价跟踪循环引用的问题，以不完整的数据结构代替了爆栈。

### 基础类型

#### NULL 类型

给定一个 hessian null 类型总是可以用 json null 类型描述即 `json.null => hessian.bool`

#### bool 类型

给定一个 hessian bool 类型总是可以用 json bool 类型描述即 `json.bool => hessian.bool`

#### int 类型

给定一个 hessian int 类型无法直接用 json number 类型描述，但是我们可以通过 json object 包装来表示即

```
json.object {
	"$class": "int", // or "java.lang.Integer"
	"$": number
} => hessian.int
```

#### long 类型

给定一个 hessian long 类型总是可以用 json number 类型描述即 `json.number => hessian.long`

#### double 类型

给定一个 hessian double 类型总是可以用 json object 类型描述即

```
json.object {
	"$class": "double", // or "java.lang.Double"
	"$": number
} => hessian.double
```

#### binary 类型

给定一个 hessian binary 类型总是可以用 json object 类型描述即

```
json.object {
	"$class": "bytes",
	"$": "base64 encoding string"
} => hessian.binary
```

#### date 类型

给定一个 hessian date 类型总是可以用 json object 类型描述即

```
json.object {
	"$class": "date", // or "java.util.Date"
	"$": number
} => hessian.date
```

### 复合类型

#### map 类型

给定一个 hessian map 类型总是可以用 json object 类型描述即 `json.object => hessian.map`

#### list 类型

给定一个 hessian map 类型总是可以用 json array 类型描述即 `json.array => hessian.list`

#### object 类型

给定一个 hessian object 类型总是可以用 json object 类型描述即

```
json.object {
	"$class": "classname",
	"$": json element
} => hessian.object
```

## HESSIAN 到 JSON 的类型转换

给定一个 JSON 类型总是存在一种 HESSIAN 类型可以等价描述出来.

### 基础类型

#### number

给定一个 json number 类型总是可以用 hessian int/long/double 类型描述即 `hessian.int|long/double => json.number`

#### string

给定一个 json string 类型总是可以用 hessian string 类型描述即 `hessian.string => json.string`

#### bool

给定一个 json bool 类型总是可以用 hessian bool 类型描述即 `hessian.bool => json.bool`

#### null

给定一个 json null 类型总是可以用 hessian null 类型描述即 `hessian.null => json.null`

### 复合类型

#### object

给定一个 json object 类型总是可以用 hessian map 类型描述即 `hessian.map => json.object`

#### array

给定一个 json array 类型总是可以用 hessian array 类型描述即 `hessian.array => json.array`

## 附录

### JSON 类型系统

```
json
    element

value
    object
    array
    string
    number
    "true"
    "false"
    "null"

```

### HESSIAN 类型系统

```
#starting production
top        ::= value

#main production
value      ::= null
           ::= binary
           ::= boolean
           ::= class-def value
           ::= date
           ::= double
           ::= int
           ::= list
           ::= long
           ::= map
           ::= object
           ::= ref
           ::= string
```

# API

## encode

查看 [examples/sofahessian_examples_test.go](/examples/sofahessian_examples_test.go#L13)

## decode

查看 [examples/sofahessian_examples_test.go](/examples/sofahessian_examples_test.go#L14)

## json

查看 [examples/sofahessian_examples_test.go](/examples/sofahessian_examples_test.go#L68)

# CLI

## install

make hessian

## decode

```bash
bin/hessian decode 4fbc636f6d2e616c697061792e736f66612e7270632e636f72652e726571756573742e536f666152657175657374950d7461726765744170704e616d650a6d6574686f644e616d651774617267657453657276696365556e697175654e616d650c7265717565737450726f70730d6d6574686f64417267536967736f904e0873617948656c6c6f1048656c6c6f536572766963653a312e304d0870726f746f636f6c04626f6c74117270635f74726163655f636f6e746578744d09736f6661527063496401300473616d700566616c73650b73797350656e4174747273000d736f666143616c6c6572496463000c736f666143616c6c65724970000b736f6661547261636549641e3061306665383633313537313034363337383735383130303138363232300c736f666150656e4174747273000e736f666143616c6c65725a6f6e65000d736f666143616c6c6572417070007a7a567400075b737472696e676e02106a6176612e6c616e672e537472696e67046c6f6e677a05776f726c64e1 --version=3

&hessian.JavaObject{
  class: "com.alipay.sofa.rpc.core.request.SofaRequest",
  names: []string{
    "targetAppName",
    "methodName",
    "targetServiceUniqueName",
    "requestProps",
    "methodArgSigs",
  },
  values: []interface {}{
    nil,
    "sayHello",
    "HelloService:1.0",
    map[interface {}]interface {}{
      "protocol":          "bolt",
      "rpc_trace_context": map[interface {}]interface {}{
        "sofaCallerZone": "",
        "sofaCallerApp":  "",
        "sofaRpcId":      "0",
        "sysPenAttrs":    "",
        "sofaCallerIdc":  "",
        "sofaPenAttrs":   "",
        "samp":           "false",
        "sofaCallerIp":   "",
        "sofaTraceId":    "0a0fe8631571046378758100186220",
      },
    },
    &hessian.JavaList{
      class: "[string",
      value: []interface {}{
        "java.lang.String",
        "long",
      },
    },
  },
}
"world"
1
```

## fromjson

```bash
bin/hessian fromjson '{"$class": "com.alipay.sofa.rpc.core.request.SofaRequest", "$": {"1": "2"}}' --format hex

43302c636f6d2e616c697061792e736f66612e7270632e636f72652e726571756573742e536f666152657175657374910131600132

bin/hessian decode 43302c636f6d2e616c697061792e736f66612e7270632e636f72652e726571756573742e536f666152657175657374910131600132
&sofahessian.JavaObject{
  class: "com.alipay.sofa.rpc.core.request.SofaRequest",
  names: []string{
    "1",
  },
  values: []interface {}{
    "2",
  },
}
```

# 性能测试

```
─λ make bench
go test -benchmem -run="^$" -bench ^Benchmark ./...
?   	gitlab.alipay-inc.com/sofa-go/sofa-hessian-go/sofahessianv1	[no test files]
goos: darwin
goarch: amd64
pkg: gitlab.alipay-inc.com/sofa-go/sofa-hessian-go/sofahessianv2
BenchmarkDecodeBinary-8             	  558362	      2180 ns/op	15041.85 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecodeBool-8               	33203718	        35.8 ns/op	  27.91 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecodeDate-8               	14283891	        73.1 ns/op	 123.15 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecodeFloat64-8            	15707289	        74.1 ns/op	  40.49 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecodeInt32-8              	15519128	        82.9 ns/op	  60.31 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecodeInt64-8              	15663241	        78.2 ns/op	 115.09 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecodeList-8               	  751725	      2969 ns/op	   6.06 MB/s	     478 B/op	      19 allocs/op
BenchmarkDecodeMap-8                	  217478	      5185 ns/op	  14.08 MB/s	     949 B/op	      30 allocs/op
BenchmarkDecodeNil-8                	27432326	        42.6 ns/op	  23.46 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecodeObject-8             	  965149	      1498 ns/op	  22.03 MB/s	     548 B/op	      13 allocs/op
BenchmarkDecodeString-8             	     818	   1333529 ns/op	 147.44 MB/s	    1370 B/op	       0 allocs/op
BenchmarkEncodeBinary-8             	  392306	      2989 ns/op	21944.34 MB/s	       0 B/op	       0 allocs/op
BenchmarkEncodeBool-8               	143134738	         8.50 ns/op	 117.64 MB/s	       0 B/op	       0 allocs/op
BenchmarkEncodeFloat64-8            	100000000	        11.0 ns/op	 819.48 MB/s	       0 B/op	       0 allocs/op
BenchmarkEncodeInt64-8              	100000000	        10.2 ns/op	 490.08 MB/s	       0 B/op	       0 allocs/op
BenchmarkEncodeInt32-8              	100000000	        10.7 ns/op	 465.65 MB/s	       0 B/op	       0 allocs/op
BenchmarkEncodeList-8               	 7413168	       161 ns/op	  12.45 MB/s	      40 B/op	       2 allocs/op
BenchmarkEncodeMap-8                	13931235	        82.9 ns/op	  24.12 MB/s	       8 B/op	       1 allocs/op
BenchmarkEncodeNil-8                	197458849	         6.12 ns/op	 163.33 MB/s	       0 B/op	       0 allocs/op
BenchmarkEncodeObject-8             	 3964602	       303 ns/op	  46.15 MB/s	      64 B/op	       3 allocs/op
BenchmarkEncodeRef-8                	88544038	        13.5 ns/op	 222.70 MB/s	       0 B/op	       0 allocs/op
BenchmarkEncodeString-8             	     716	   1662180 ns/op	 160.43 MB/s	      16 B/op	       1 allocs/op
BenchmarkEncodeAndDecodeJSON-8      	   51471	     22695 ns/op	  44.68 MB/s	    2913 B/op	      37 allocs/op
BenchmarkEncodeAndDecodeHessianV2-8   	  139880	      7979 ns/op	  56.77 MB/s	    3104 B/op	      57 allocs/op
PASS
ok  	gitlab.alipay-inc.com/sofa-go/sofa-hessian-go/sofahessianv2	34.443s
?   	gitlab.alipay-inc.com/sofa-go/sofa-hessian-go/javaobject	[no test files]
```

# TODO

1. hessian-generator: 通过编译时代码生成 Encode 和 Decode 方法，减少反射带来的负担。
