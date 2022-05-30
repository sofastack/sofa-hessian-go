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
	"encoding/binary"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type BytesBufioReader struct {
	r  *bytes.Reader
	br *bufio.Reader
}

func (bbr *BytesBufioReader) Reset(b []byte) {
	bbr.r.Reset(b)
	bbr.br.Reset(bbr.r)
}

func (bbr *BytesBufioReader) GetBufioReader() *bufio.Reader {
	return bbr.br
}

func getInterfaceName(i interface{}) string {
	g, ok := i.(JavaClassNameGetter)
	if ok {
		return g.GetJavaClassName()
	}

	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}

func b2s(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func s2b(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func allocAtLeast(dst []byte, length int) []byte {
	dc := cap(dst)
	n := len(dst) + length
	if dc < n {
		dst = dst[:dc]
		dst = append(dst, make([]byte, n-dc)...)
	}
	dst = dst[:n]
	return dst
}

func readAtLeastBytesFromReader(reader io.Reader, length int, buf []byte) error {
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		return err
	}
	if n < length {
		return ErrDecodeBufferNotEnough
	}

	return nil
}

func readUint8FromReader(reader io.Reader) (uint8, error) {
	p := acquireU8()
	n, err := io.ReadFull(reader, (*p)[7:])
	if err != nil {
		releaseU8(p)
		return 0, err
	}

	if n < 1 {
		releaseU8(p)
		return 0, ErrDecodeBufferNotEnough
	}

	u8 := p[7]
	releaseU8(p)

	return u8, nil
}

func readUint16FromReader(reader io.Reader) (uint16, error) {
	p := acquireU8()
	n, err := io.ReadFull(reader, (*p)[6:])
	if err != nil {
		releaseU8(p)
		return 0, err
	}

	if n < 2 {
		releaseU8(p)
		return 0, ErrDecodeBufferNotEnough
	}

	u16 := binary.BigEndian.Uint16(p[6:])
	releaseU8(p)

	return u16, nil
}

func readInt16FromReader(reader io.Reader) (int16, error) {
	u16, err := readUint16FromReader(reader)
	return int16(u16), err
}

func readInt32FromReader(reader io.Reader) (int32, error) {
	u32, err := readUint32FromReader(reader)
	return int32(u32), err
}

func readUTF8StringFromReader(reader *bufio.Reader, dst []byte, length int) ([]byte, error) {
	for i := 0; i < length; i++ {
		r, s, err := reader.ReadRune()
		if err != nil {
			return dst, err
		}
		dst = appendRune(dst, uint32(r), s)
	}
	return dst, nil
}

func readLenAndUTF8StringFromReader(reader *bufio.Reader, dst []byte) ([]byte, error) {
	u16, err := readUint16FromReader(reader)
	if err != nil {
		return dst, err
	}

	if u16 == 0 {
		return dst, nil
	}

	return readUTF8StringFromReader(reader, dst, int(u16))
}

func readUint32FromReader(reader io.Reader) (uint32, error) {
	p := acquireU8()
	n, err := io.ReadFull(reader, (*p)[4:])
	if err != nil {
		releaseU8(p)
		return 0, err
	}

	if n < 4 {
		releaseU8(p)
		return 0, ErrDecodeBufferNotEnough
	}

	u32 := binary.BigEndian.Uint32(p[4:])
	releaseU8(p)

	return u32, nil
}

func readUint64FromReader(reader io.Reader) (uint64, error) {
	// TODO(detailyang): until golang do not escape the p to heap because of the io.Reader interface.
	p := acquireU8()
	n, err := io.ReadFull(reader, (*p)[:])
	if err != nil {
		releaseU8(p)
		return 0, err
	}

	if n < 8 {
		releaseU8(p)
		return 0, ErrDecodeBufferNotEnough
	}

	u64 := binary.BigEndian.Uint64((*p)[:])
	releaseU8(p)

	return u64, nil
}

type utf8Helper struct {
	offsetSlice []int
	countSlice  []int
}

func (u *utf8Helper) GetOffsetSlice() []int {
	return u.offsetSlice
}

func (u *utf8Helper) GetCountSlice() []int {
	return u.countSlice
}

func (u *utf8Helper) grow(n int) {
	co := cap(u.offsetSlice)
	cc := cap(u.countSlice)
	if co < n {
		u.offsetSlice = u.offsetSlice[:co]
		u.offsetSlice = append(u.offsetSlice, make([]int, n-co)...)
	}
	if cc < n {
		u.countSlice = u.countSlice[:cc]
		u.countSlice = append(u.countSlice, make([]int, n-cc)...)
	}
	u.offsetSlice = u.offsetSlice[:n]
	u.countSlice = u.countSlice[:n]
}

func (u *utf8Helper) Reset() {
	u.offsetSlice = u.offsetSlice[:0]
	u.countSlice = u.countSlice[:0]
}

func (u *utf8Helper) Write(s string) (int, error) {
	u.grow(len(s))
	os := u.offsetSlice
	cs := u.countSlice

	var (
		b      = bytes.NewBuffer(s2b(s))
		err    error
		size   int
		offset = 0
		count  = 0
	)

	for {
		_, size, err = b.ReadRune()
		if err != nil {
			break
		}
		for i := 0; i < size; i++ {
			os[offset+i] = count
		}
		cs[count] = offset
		count++
		offset += size
	}

	u.offsetSlice = os
	u.countSlice = cs

	if err == io.EOF {
		return count, nil
	}

	return count, err
}

func appendRune(p []byte, r uint32, size int) []byte {
	const (
		t1    = 0   // 0b00000000
		tx    = 128 // 0b10000000
		t2    = 192 // 0b11000000
		t3    = 224 // 0b11100000
		t4    = 240 // 0b11110000
		t5    = 248 // 0b11111000
		maskx = 63  // 0b00111111
		mask2 = 31  // 0b00011111
		mask3 = 15  // 0b00001111
		mask4 = 7   // 0b00000111
	)

	switch size {
	case 1:
		p = append(p, byte(r))
	case 2:
		p = append(p, t2|byte(r>>6), tx|byte(r)&maskx)
	case 3:
		p = append(p, t3|byte(r>>12), tx|byte(r>>6)&maskx, tx|byte(r)&maskx)
	default:
		p = append(p, t4|byte(r>>18), tx|byte(r>>12)&maskx, tx|byte(r>>6)&maskx, tx|byte(r)&maskx)
	}

	return p
}

func escapeString(dst []byte, s string) []byte {
	if !hasSpecialChars(s) {
		// Fast path - nothing to escape.
		dst = append(dst, '"')
		dst = append(dst, s...)
		dst = append(dst, '"')
		return dst
	}

	// Slow path.
	return strconv.AppendQuote(dst, s)
}

func hasSpecialChars(s string) bool {
	if strings.IndexByte(s, '"') >= 0 || strings.IndexByte(s, '\\') >= 0 {
		return true
	}
	for i := 0; i < len(s); i++ {
		if s[i] < 0x20 {
			return true
		}
	}
	return false
}
