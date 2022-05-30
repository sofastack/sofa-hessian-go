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
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sofastack/sofa-hessian-go/javaobject"
	hessian "github.com/sofastack/sofa-hessian-go/sofahessian"

	"github.com/k0kubun/pp"
	"github.com/kr/pretty"
	"github.com/sofastack/sofa-common-go/helper/easyreader"
	"github.com/spf13/cobra"
	jp "github.com/tidwall/pretty"
)

func init() {
	scheme := pp.ColorScheme{
		Integer:         pp.Green,
		Float:           pp.Black | pp.BackgroundWhite | pp.Bold,
		String:          pp.Red,
		Bool:            pp.Yellow,
		StringQuotation: pp.White,
		EscapedChar:     pp.Magenta,
		FieldName:       pp.Cyan,
		PointerAdress:   pp.White,
		Nil:             pp.Blue,
		Time:            pp.Magenta,
		StructName:      pp.Green | pp.Bold,
		ObjectLength:    pp.White,
	}

	// Register it for usage
	pp.SetColorScheme(scheme)
}

var (
	decodeOutputType           string
	decodeVerbose              bool
	decodeShowtype             bool
	decodeVersion              int
	decodeLoadJavaBuiltinClass bool
	decodeLoadSofaRPCClass     bool
	decodeTee                  string
	decodeFormat               string
)

func init() {
	fs := decodeCmd.Flags()
	fs.BoolVar(&decodeVerbose, "verbose", false, "Show verbose logs")
	fs.BoolVar(&decodeShowtype, "showtype", false, "Show type")
	fs.BoolVar(&decodeLoadJavaBuiltinClass, "loadjavabuiltin", true, "Load java builtin class")
	fs.BoolVar(&decodeLoadSofaRPCClass, "loadsofarpc", true, "Load sofa rpc class")
	fs.StringVarP(&decodeFormat, "in-format", "i", "hex", "Set the input format")
	fs.StringVar(&decodeTee, "tee", "hex", "output the raw bytes to stdout (format hex or bin)")
	fs.IntVar(&decodeVersion, "version", 3, "Set the hessian version (1 or 3 or 4)")
	fs.StringVar(&decodeOutputType, "out-format", "go", "Set the output format")
}

var decodeCmd = &cobra.Command{
	Use:   "decode <input>",
	Short: "decode input to pretty stdout",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// nolint
		if decodeLoadJavaBuiltinClass {
			hessian.RegisterJavaClass(javaobject.JavaLangStackTraceElement{})
			hessian.RegisterJavaClass(javaobject.JavaLangStackTraceElements{})
		}

		// nolint
		if decodeLoadSofaRPCClass {
			hessian.RegisterJavaClass(&javaobject.SofaRPCRequest{})
			hessian.RegisterJavaClass(&javaobject.SofaRPCResponse{})
			hessian.RegisterJavaClass(&javaobject.SofaRPCServerException{})
		}

		var o *easyreader.Option
		if decodeFormat == "hex" {
			o = easyreader.NewOption().
				SetDefaultFormat(easyreader.HexFormat)
		} else {
			o = easyreader.NewOption().
				SetDefaultFormat(easyreader.BinFormat)
		}

		reader, err := easyreader.EasyRead(
			o,
			args[0])
		if err != nil {
			log.Fatal(err)
		}

		dctx := hessian.NewDecodeContext()
		if decodeVerbose {
			dctx.SetTracer(hessian.NewStderrTracer())
		}

		dctx.SetVersion(getHessianVersion(decodeVersion))

		raw := bytes.NewBuffer(nil)
		br := bufio.NewReaderSize(io.TeeReader(reader, raw), 1024*1024)
		decoder := hessian.NewDecoder(dctx, br)

		for {
			obj, err := decoder.Decode()
			if err != nil {
				if err == io.EOF {
					if decodeTee == "hex" {
						fmt.Printf(hex.EncodeToString(raw.Bytes()))
					} else {
						fmt.Printf(raw.String())
					}
					return
				}
				log.Fatal(err)
			}
			output(obj)
		}
	},
}

func output(obj interface{}) {
	switch strings.ToLower(decodeOutputType) {
	case "json":
		outputJSON(obj)
	case "go":
		fallthrough
	default:
		outputGo(obj)
	}
}

func outputJSON(obj interface{}) {
	jctx := hessian.NewJSONEncodeContext().SetMaxPtrCycles(1)
	if decodeVerbose {
		jctx.SetTracer(hessian.NewStdoutTracer())
	}
	dst, err := hessian.EncodeToJSON(jctx, nil, obj)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(os.Stderr, string(jp.Color(jp.Pretty(dst), nil)))
}

type Context struct{}

func outputGo(obj interface{}) {
	if decodeShowtype {
		// nolint
		pretty.Fprintf(os.Stderr, pretty.Sprint(obj))
	} else {
		// nolint
		pp.Fprintln(os.Stderr, obj)
	}
}
