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
	"fmt"
	"io"
	"os"
	"strings"
)

type Tracer interface {
	OnTraceStart(event string)
	OnTraceStop(event string)
}

type DummyTracer struct{}

func NewDummyTracer() *DummyTracer { return &DummyTracer{} }

func (t *DummyTracer) OnTraceStart(event string) { _ = event }

func (t *DummyTracer) OnTraceStop(event string) { _ = event }

type StdoutTracer struct {
	depth int
}

func NewStdoutTracer() *StdoutTracer { return &StdoutTracer{} }

func (t *StdoutTracer) OnTraceStart(event string) {
	t.depth++
	fmt.Printf("%s<start %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
}

func (t *StdoutTracer) OnTraceStop(event string) {
	fmt.Printf("%s<stop %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
	t.depth--
}

type StderrTracer struct {
	depth int
}

func NewStderrTracer() *StderrTracer { return &StderrTracer{} }

func (t *StderrTracer) OnTraceStart(event string) {
	t.depth++
	fmt.Fprintf(os.Stderr, "%s<start %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
}

func (t *StderrTracer) OnTraceStop(event string) {
	fmt.Fprintf(os.Stderr, "%s<stop %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
	t.depth--
}

type WriterTracer struct {
	w     io.Writer
	depth int
}

func NewWriterTracer(w io.Writer) *WriterTracer {
	return &WriterTracer{
		w:     w,
		depth: 0,
	}
}

func (t *WriterTracer) OnTraceStart(event string) {
	t.depth++
	fmt.Fprintf(t.w, "%s<start %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
}

func (t *WriterTracer) OnTraceStop(event string) {
	fmt.Fprintf(t.w, "%s<stop %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
	t.depth--
}
