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
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeString(t *testing.T) {
	t.Run("should read and write utf8 string as java", func(t *testing.T) {
		for _, i := range []int{
			32767,
			32768,
			32769,
			65534,
			65535,
			65536,
			65537,
		} {
			testDecodeString(t,
				readFile(t, fmt.Sprintf("testdata/string/utf8_%d.bin", i)),
				strings.Repeat("锋", i),
			)
		}
	})

	t.Run("should read short strings", func(t *testing.T) {
		testDecodeString(t, []byte{0x00}, "")
		testDecodeString(t, []byte{0x05, 'h', 'e', 'l', 'l', 'o'}, "hello")
		testDecodeString(t, []byte{0x01, 0xc3, 0x83}, "\u00c3")
		x := append([]byte{}, 0x09)
		x = append(x, "hello, 中文"...)
		testDecodeString(t, x, "hello, 中文")
	})

	t.Run("should read hello in long form", func(t *testing.T) {
		testDecodeString(t, []byte{'S', 0x00, 0x05, 'h', 'e', 'l', 'l', 'o'}, "hello")
	})

	t.Run("should read split into two chunks: s and short strings", func(t *testing.T) {
		testDecodeString(t, []byte{
			0x52, 0x00, 0x07,
			'h', 'e', 'l', 'l', 'o', ',', ' ',
			0x05, 'w', 'o', 'r', 'l', 'd',
		}, "hello, world")
	})

	t.Run("should read golden files", func(t *testing.T) {
		for _, i := range [][2]string{
			{"empty.bin", ""},
			{"foo.bin", "foo"},
			{"chinese.bin", "中文 Chinese"},
		} {
			testDecodeString(t,
				readFile(t, fmt.Sprintf("testdata/string/%s", i[0])),
				i[1],
			)
		}

		testDecodeString(t,
			readFile(t, "testdata/string/text4k.bin"),
			string(readFile(t, "testdata/string/4k.txt")),
		)

		testDecodeString(t,
			readFile(t, "testdata/string/large_string_32767.bin"),
			strings.Repeat("\x41", 32767),
		)

		testDecodeString(t,
			readFile(t, "testdata/string/large_string_32768.bin"),
			strings.Repeat("\x41", 32768),
		)

		testDecodeString(t,
			readFile(t, "testdata/string/large_string_65534.bin"),
			strings.Repeat("\x41", 65534),
		)

		testDecodeString(t,
			readFile(t, "testdata/string/large_string_32769.bin"),
			strings.Repeat("\x41", 32769),
		)

		testDecodeString(t,
			readFile(t, "testdata/string/large_string_65534.bin"),
			strings.Repeat("\x41", 65534),
		)

		testDecodeString(t,
			readFile(t, "testdata/string/large_string_65535.bin"),
			strings.Repeat("\x41", 65535),
		)

		testDecodeString(t,
			readFile(t, "testdata/string/large_string_65536.bin"),
			strings.Repeat("\x41", 65536),
		)

		testDecodeString(t,
			readFile(t, "testdata/string/large_string_65537.bin"),
			strings.Repeat("\x41", 65537),
		)
	})
}

func testDecodeString(t *testing.T, b []byte, s string) {
	o := &DecodeContext{
		tracer: NewDummyTracer(),
	}

	x, err := DecodeStringHessian4V2(o, bufio.NewReader(bytes.NewReader(b)))
	require.Nil(t, err)
	require.Equal(t, s, string(x))
}
