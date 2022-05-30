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
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	hessian "github.com/sofastack/sofa-hessian-go/sofahessian"

	"github.com/sofastack/sofa-common-go/helper/easyreader"
	"github.com/spf13/cobra"
)

var (
	fromjsonVerbose     bool
	fromjsonVersion     int
	fromjsonWriteFile   string
	fromjsonWriteFormat string
)

func init() {
	fs := fromjsonCmd.Flags()
	fs.BoolVar(&fromjsonVerbose, "verbose", false, "Show verbose logs")
	fs.IntVar(&fromjsonVersion, "version", 3, "Set the hessian version (1 or 3 or 4)")
	fs.StringVar(&fromjsonWriteFile, "file", "", "Write to the file (default STDOUT)")
	fs.StringVar(&fromjsonWriteFormat, "format", "hex", "Write to the file (default HEX)")
}

var fromjsonCmd = &cobra.Command{
	Use:   "fromjson <input>",
	Short: "fromjson read from json then transcode to hessian encoding",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
			jw  bytes.Buffer
			ew  bytes.Buffer
		)
		jtracer := hessian.NewWriterTracer(&jw)
		etracer := hessian.NewWriterTracer(&ew)

		ectx := hessian.NewEncodeContext()
		if fromjsonVerbose {
			ectx.SetTracer(etracer)
		}

		ectx.SetVersion(getHessianVersion(fromjsonVersion))
		encoder := hessian.NewEncoder(ectx)
		jctx := hessian.NewJSONEncodeContext()

		if fromjsonVerbose {
			jctx.SetTracer(jtracer)
		}

		// nolint
		for i := range args {
			reader, err := easyreader.EasyRead(easyreader.NewOption().SetDefaultFormat(easyreader.BinFormat), args[i])
			if err != nil {
				log.Fatal(err)
			}

			data, err := ioutil.ReadAll(reader)
			if err != nil {
				log.Fatal(err)
			}

			if err = encoder.EncodeStreamingJSONBytes(jctx, data); err != nil {
				log.Fatal(err)
			}
		}

		if fromjsonVerbose {
			fmt.Println(string(jw.Bytes()))
			fmt.Println(string(ew.Bytes()))
		}

		var writer io.Writer
		if fromjsonWriteFile == "" {
			writer = os.Stdout
		} else {
			writer, err = os.OpenFile(fromjsonWriteFile, os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		b := encoder.Bytes()

		// nolint
		switch strings.ToLower(fromjsonWriteFormat) {
		case "hexp":
			writer.Write([]byte("0x"))
			fallthrough
		case "hex":
			writer.Write([]byte(hex.EncodeToString(b)))
		default:
			writer.Write(b)
		}
	},
}
