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

package cmd

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/sofastack/sofa-hessian-go/javaobject"
	hessian "github.com/sofastack/sofa-hessian-go/sofahessian"

	"github.com/sofastack/sofa-common-go/helper/easyreader"
	"github.com/spf13/cobra"
	"github.com/valyala/fastjson"
)

var (
	fromsofaVerbose bool
	fromsofaVersion int
	fromsofaFormat  string
	// fromsofaWriteFile   string
	// fromsofaWriteFormat string
	jsonparserPool fastjson.ParserPool
)

func init() {
	fs := fromsofaCmd.Flags()
	fs.BoolVar(&fromsofaVerbose, "verbose", false, "Show verbose logs")
	fs.IntVar(&fromsofaVersion, "version", 3, "Set the hessian version (1 or 3 or 4)")
	// fs.StringVar(&fromjsonWriteFile, "file", "", "Write to the file (default STDOUT)")
	// fs.StringVar(&fromjsonWriteFormat, "format", "hex", "Write to the file (default HEX)")
	fs.StringVar(&fromsofaFormat, "format", "hex", "Set the output format (hex or bin)")
}

var fromsofaCmd = &cobra.Command{
	Use:   "fromsofa interface method arguments",
	Short: "fromsofa read from json then transcode to sofarequest hessian encoding",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if err := transcodeToSofaRequest(args); err != nil {
			log.Fatal(err)
		}
	},
}

func transcodeToSofaRequest(args []string) error {
	inter := args[0]
	method := args[1]
	argument := args[2]

	p := jsonparserPool.Get()
	defer jsonparserPool.Put(p)

	reader, err := easyreader.EasyRead(easyreader.
		NewOption().
		SetDefaultFormat(easyreader.BinFormat),
		argument)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	v, err := p.ParseBytes(data)
	if err != nil {
		return err
	}

	obj := v.GetObject()
	if obj == nil {
		return fmt.Errorf("expect json object but got %s", v.Type().String())
	}

	var req javaobject.SofaRPCRequest

	d, err := jsonToSofaHessian(inter, method, getHessianVersion(fromsofaVersion), &req, obj)
	if err != nil {
		return err
	}

	if fromsofaFormat == "hex" {
		fmt.Printf(hex.EncodeToString(d))
	} else {
		fmt.Printf(string(d))
	}

	return nil
}

func jsonToSofaHessian(inter, method string, version hessian.Version,
	req *javaobject.SofaRPCRequest, obj *fastjson.Object) ([]byte, error) {
	a := obj.Get("arguments")
	if a == nil {
		return nil, fmt.Errorf("expect jsonobject.arguments but not")
	}

	arguments, err := a.Array()
	if err != nil {
		return nil, fmt.Errorf("expect jsonobject.arguments is array but got %s", a.Type().String())
	}

	s := obj.Get("signatures")
	if s != nil {
		var sig []*fastjson.Value
		sig, err = s.Array()
		if err != nil {
			return nil, fmt.Errorf("expect jsonobject.signatures is array but got %s", a.Type().String())
		}
		for i := range sig {
			var s []byte
			s, err = sig[i].StringBytes()
			if err != nil {
				return nil, err
			}
			// short life cycle
			req.MethodArgSigs = append(req.MethodArgSigs, string(s))
		}

	} else {
		for i := range arguments {
			arg := arguments[i]
			switch arg.Type() {
			case fastjson.TypeObject:
				obj := arg.GetObject()
				if obj == nil {
					req.MethodArgSigs = append(req.MethodArgSigs, "null")
					break
				}
				v := obj.Get("$class")
				if v == nil {
					req.MethodArgSigs = append(req.MethodArgSigs, "unknow class")
				} else {
					class := v.GetStringBytes()
					req.MethodArgSigs = append(req.MethodArgSigs, string(class))
				}

			case fastjson.TypeArray:
				req.MethodArgSigs = append(req.MethodArgSigs, "[")

			case fastjson.TypeNull:
				req.MethodArgSigs = append(req.MethodArgSigs, "null")

			case fastjson.TypeString:
				req.MethodArgSigs = append(req.MethodArgSigs, "string")

			case fastjson.TypeNumber:
				f64 := arg.GetFloat64()
				if f64 == float64(int(f64)) {
					req.MethodArgSigs = append(req.MethodArgSigs, "long")
				} else {
					req.MethodArgSigs = append(req.MethodArgSigs, "double")
				}

			case fastjson.TypeTrue:
				fallthrough
			case fastjson.TypeFalse:
				req.MethodArgSigs = append(req.MethodArgSigs, "boolean")
			}
		}
	}

	req.TargetAppName = ""
	req.TargetServiceUniqueName = inter
	req.MethodName = method

	ectx := hessian.AcquireHessianEncodeContext().SetVersion(version)
	encoder := hessian.AcquireHessianEncoder(ectx)
	jctx := hessian.AcquireJSONContext().SetMaxDepth(32)
	defer func() {
		hessian.ReleaseJSONContext(jctx)
		hessian.ReleaseHessianEncoder(encoder)
		hessian.ReleaseHessianEncodeContext(ectx)
	}()

	// write sofarequest
	if err = encoder.Encode(req); err != nil {
		hessian.ReleaseHessianEncodeContext(ectx)
		hessian.ReleaseHessianEncoder(encoder)
		return nil, err
	}

	fmt.Println("cc", hex.EncodeToString(encoder.Bytes()))

	for i := range arguments {
		if err = encoder.EncodeFastJSON(jctx, arguments[i]); err != nil {
			break
		}
	}

	return encoder.Bytes(), err
}
