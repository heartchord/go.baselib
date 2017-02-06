package goblazer

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"fmt"

	"github.com/henrylee2cn/mahonia"
)

type tabCell struct {
	content string
}

func newTabCell(content string) *tabCell {
	tc := new(tabCell)
	tc.content = content
	return tc
}

type tabOffset struct {
	offset int
	length int
}

// TabFile opens a "xxx.tab" file to resolve it.
type TabFile struct {
	rows int
	cols int
	tabs []*tabCell
}

// Load is
func (f *TabFile) Load(path string, code string) bool {
	var ok bool
	var fi *os.File
	var err error
	var buff []byte
	var size int

	if fi, err = os.Open(path); err != nil {
		fmt.Println(err)
		return false
	}
	defer fi.Close()

	if buff, err = ioutil.ReadAll(fi); err != nil {
		return false
	}

	if ok = strings.EqualFold(code, "utf8"); !ok {
		var s string
		mdecoder := mahonia.NewDecoder(strings.ToUpper(code))

		if s, ok = mdecoder.ConvertStringOK(string(buff)); !ok {
			return false
		}
		buff = []byte(s)
	}

	size = len(buff)
	if size > 0 {
		f.createTabOffsets(buff, size)
	}

	return true
}

// Save is
func (f *TabFile) Save(path string, code string) bool {
	var ok, isutf8 bool
	var err error
	var buff bytes.Buffer
	var str string
	var fi *os.File

	for i := 0; i < f.rows; i++ {
		for j := 0; j < f.cols; j++ {
			idx := i*f.cols + j
			if j < f.cols-1 {
				buff.WriteString(fmt.Sprintf("%s\t", f.tabs[idx].content))
			} else {
				buff.WriteString(fmt.Sprintf("%s\r\n", f.tabs[idx].content))
			}
		}
	}

	str = buff.String()
	isutf8 = strings.EqualFold(code, "utf8")
	if !isutf8 {
		code = strings.ToUpper(code)
		mencoder := mahonia.NewEncoder(code)
		if str, ok = mencoder.ConvertStringOK(str); !ok {
			return false
		}
	}

	if fi, err = os.Create(path); err != nil {
		return false
	}
	defer fi.Close()

	if _, err = fi.WriteString(str); err != nil {
		return false
	}

	return true
}

// Reset is
func (f *TabFile) Reset() {
	f.rows = 0
	f.cols = 0
	f.tabs = nil
}

// PrintTabInfo is
func (f *TabFile) PrintTabInfo() {
	fmt.Printf("rows = %d\n", f.rows)
	fmt.Printf("cols = %d\n", f.cols)

	for i := 0; i < f.rows; i++ {
		fmt.Printf("line %d: \n", i)
		for j := 0; j < f.cols; j++ {
			idx := i*f.cols + j
			fmt.Printf("tab[%d][%d] = %s, ", i, j, f.tabs[idx].content)
		}
		fmt.Println()
	}
}

// GetIntByIntIdx is
func (f *TabFile) GetIntByIntIdx(row int, col int, dflt int) int {
	ret := dflt
	if s, ok := f.GetCell(row, col); ok {
		if v, err := strconv.Atoi(s); err == nil {
			ret = v
		}
	}
	return ret
}

// GetIntByStrIdx is
func (f *TabFile) GetIntByStrIdx(row string, col string, dflt int) int {
	return f.GetIntByIntIdx(f.FindRow(row), f.FindCol(col), dflt)
}

// GetIntByMixIdx is
func (f *TabFile) GetIntByMixIdx(row int, col string, dflt int) int {
	return f.GetIntByIntIdx(row, f.FindCol(col), dflt)
}

// GetByteByIntIdx is
func (f *TabFile) GetByteByIntIdx(row int, col int, dflt byte) byte {
	ret := dflt
	if s, ok := f.GetCell(row, col); ok {
		if v, err := strconv.ParseInt(s, 10, 8); err == nil {
			ret = byte(v)
		}
	}
	return ret
}

// GetByteByStrIdx is
func (f *TabFile) GetByteByStrIdx(row string, col string, dflt byte) byte {
	return f.GetByteByIntIdx(f.FindRow(row), f.FindCol(col), dflt)
}

// GetByteByMixIdx is
func (f *TabFile) GetByteByMixIdx(row int, col string, dflt byte) byte {
	return f.GetByteByIntIdx(row, f.FindCol(col), dflt)
}

// GetInt8ByIntIdx is
func (f *TabFile) GetInt8ByIntIdx(row int, col int, dflt int8) int8 {
	ret := dflt
	if s, ok := f.GetCell(row, col); ok {
		if v, err := strconv.ParseInt(s, 10, 8); err == nil {
			ret = int8(v)
		}
	}
	return ret
}

// GetInt8ByStrIdx is
func (f *TabFile) GetInt8ByStrIdx(row string, col string, dflt int8) int8 {
	return f.GetInt8ByIntIdx(f.FindRow(row), f.FindCol(col), dflt)
}

// GetInt8ByMixIdx is
func (f *TabFile) GetInt8ByMixIdx(row int, col string, dflt int8) int8 {
	return f.GetInt8ByIntIdx(row, f.FindCol(col), dflt)
}

// GetInt16ByIntIdx is
func (f *TabFile) GetInt16ByIntIdx(row int, col int, dflt int16) int16 {
	ret := dflt
	if s, ok := f.GetCell(row, col); ok {
		if v, err := strconv.ParseInt(s, 10, 16); err == nil {
			ret = int16(v)
		}
	}
	return ret
}

// GetInt16ByStrIdx is
func (f *TabFile) GetInt16ByStrIdx(row string, col string, dflt int16) int16 {
	return f.GetInt16ByIntIdx(f.FindRow(row), f.FindCol(col), dflt)
}

// GetInt16ByMixIdx is
func (f *TabFile) GetInt16ByMixIdx(row int, col string, dflt int16) int16 {
	return f.GetInt16ByIntIdx(row, f.FindCol(col), dflt)
}

// GetInt32ByIntIdx is
func (f *TabFile) GetInt32ByIntIdx(row int, col int, dflt int32) int32 {
	ret := dflt
	if s, ok := f.GetCell(row, col); ok {
		if v, err := strconv.ParseInt(s, 10, 32); err == nil {
			ret = int32(v)
		}
	}
	return ret
}

// GetInt32ByStrIdx is
func (f *TabFile) GetInt32ByStrIdx(row string, col string, dflt int32) int32 {
	return f.GetInt32ByIntIdx(f.FindRow(row), f.FindCol(col), dflt)
}

// GetInt32ByMixIdx is
func (f *TabFile) GetInt32ByMixIdx(row int, col string, dflt int32) int32 {
	return f.GetInt32ByIntIdx(row, f.FindCol(col), dflt)
}

// GetInt64ByIntIdx is
func (f *TabFile) GetInt64ByIntIdx(row int, col int, dflt int64) int64 {
	ret := dflt
	if s, ok := f.GetCell(row, col); ok {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			ret = v
		}
	}
	return ret
}

// GetInt64ByStrIdx is
func (f *TabFile) GetInt64ByStrIdx(row string, col string, dflt int64) int64 {
	return f.GetInt64ByIntIdx(f.FindRow(row), f.FindCol(col), dflt)
}

// GetInt64ByMixIdx is
func (f *TabFile) GetInt64ByMixIdx(row int, col string, dflt int64) int64 {
	return f.GetInt64ByIntIdx(row, f.FindCol(col), dflt)
}

// GetStrByIntIdx is
func (f *TabFile) GetStrByIntIdx(row int, col int, dflt string) string {
	ret := dflt
	if s, ok := f.GetCell(row, col); ok {
		ret = s
	}
	return ret
}

// GetStrByStrIdx is
func (f *TabFile) GetStrByStrIdx(row string, col string, dflt string) string {
	return f.GetStrByIntIdx(f.FindRow(row), f.FindCol(col), dflt)
}

// GetStrByMixIdx is
func (f *TabFile) GetStrByMixIdx(row int, col string, dflt string) string {
	return f.GetStrByIntIdx(row, f.FindCol(col), dflt)
}

// GetFloat32ByIntIdx is
func (f *TabFile) GetFloat32ByIntIdx(row int, col int, dflt float32) float32 {
	ret := dflt
	if s, ok := f.GetCell(row, col); ok {
		if v, err := strconv.ParseFloat(s, 32); err == nil {
			ret = float32(v)
		}
	}
	return ret
}

// GetFloat32ByStrIdx is
func (f *TabFile) GetFloat32ByStrIdx(row string, col string, dflt float32) float32 {
	return f.GetFloat32ByIntIdx(f.FindRow(row), f.FindCol(col), dflt)
}

// GetFloat32ByMixIdx is
func (f *TabFile) GetFloat32ByMixIdx(row int, col string, dflt float32) float32 {
	return f.GetFloat32ByIntIdx(row, f.FindCol(col), dflt)
}

// GetFloat64ByIntIdx is
func (f *TabFile) GetFloat64ByIntIdx(row int, col int, dflt float64) float64 {
	ret := dflt
	if s, ok := f.GetCell(row, col); ok {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			ret = v
		}
	}
	return ret
}

// GetFloat64ByStrIdx is
func (f *TabFile) GetFloat64ByStrIdx(row string, col string, dflt float64) float64 {
	return f.GetFloat64ByIntIdx(f.FindRow(row), f.FindCol(col), dflt)
}

// GetFloat64ByMixIdx is
func (f *TabFile) GetFloat64ByMixIdx(row int, col string, dflt float64) float64 {
	return f.GetFloat64ByIntIdx(row, f.FindCol(col), dflt)
}

// FindCol is
func (f *TabFile) FindCol(col string) int {
	ret := -1
	for i := 0; i < f.cols; i++ {
		if s, ok := f.GetCell(0, i); ok {
			if strings.EqualFold(col, s) {
				ret = i
				break
			}
		}
	}
	return ret
}

// FindRow is
func (f *TabFile) FindRow(row string) int {
	ret := -1
	for i := 0; i < f.rows; i++ {
		if s, ok := f.GetCell(i, 0); ok {
			if strings.EqualFold(row, s) {
				ret = i
				break
			}
		}
	}
	return ret
}

// GetCell is
func (f *TabFile) GetCell(row int, col int) (string, bool) {
	if row >= f.rows || col >= f.cols || row < 0 || col < 0 {
		return "", false
	}

	idx := row*f.cols + col
	return f.tabs[idx].content, true
}

// SetIntByIntIdx is
func (f *TabFile) SetIntByIntIdx(row int, col int, val int) bool {
	return f.SetCell(row, col, strconv.Itoa(val))
}

// SetIntByStrIdx is
func (f *TabFile) SetIntByStrIdx(row string, col string, val int) bool {
	return f.SetIntByIntIdx(f.FindRow(row), f.FindCol(col), val)
}

// SetIntByMixIdx is
func (f *TabFile) SetIntByMixIdx(row int, col string, val int) bool {
	return f.SetIntByIntIdx(row, f.FindCol(col), val)
}

// SetByteByIntIdx is
func (f *TabFile) SetByteByIntIdx(row int, col int, val byte) bool {
	return f.SetCell(row, col, strconv.FormatInt(int64(val), 10))
}

// SetByteByStrIdx is
func (f *TabFile) SetByteByStrIdx(row string, col string, val byte) bool {
	return f.SetByteByIntIdx(f.FindRow(row), f.FindCol(col), val)
}

// SetByteByMixIdx is
func (f *TabFile) SetByteByMixIdx(row int, col string, val byte) bool {
	return f.SetByteByIntIdx(row, f.FindCol(col), val)
}

// SetInt8ByIntIdx is
func (f *TabFile) SetInt8ByIntIdx(row int, col int, val int8) bool {
	return f.SetCell(row, col, strconv.FormatInt(int64(val), 10))
}

// SetInt8ByStrIdx is
func (f *TabFile) SetInt8ByStrIdx(row string, col string, val int8) bool {
	return f.SetInt8ByIntIdx(f.FindRow(row), f.FindCol(col), val)
}

// SetInt8ByMixIdx is
func (f *TabFile) SetInt8ByMixIdx(row int, col string, val int8) bool {
	return f.SetInt8ByIntIdx(row, f.FindCol(col), val)
}

// SetInt16ByIntIdx is
func (f *TabFile) SetInt16ByIntIdx(row int, col int, val int16) bool {
	return f.SetCell(row, col, strconv.FormatInt(int64(val), 10))
}

// SetInt16ByStrIdx is
func (f *TabFile) SetInt16ByStrIdx(row string, col string, val int16) bool {
	return f.SetInt16ByIntIdx(f.FindRow(row), f.FindCol(col), val)
}

// SetInt16ByMixIdx is
func (f *TabFile) SetInt16ByMixIdx(row int, col string, val int16) bool {
	return f.SetInt16ByIntIdx(row, f.FindCol(col), val)
}

// SetInt32ByIntIdx is
func (f *TabFile) SetInt32ByIntIdx(row int, col int, val int32) bool {
	return f.SetCell(row, col, strconv.FormatInt(int64(val), 10))
}

// SetInt32ByStrIdx is
func (f *TabFile) SetInt32ByStrIdx(row string, col string, val int32) bool {
	return f.SetInt32ByIntIdx(f.FindRow(row), f.FindCol(col), val)
}

// SetInt32ByMixIdx is
func (f *TabFile) SetInt32ByMixIdx(row int, col string, val int32) bool {
	return f.SetInt32ByIntIdx(row, f.FindCol(col), val)
}

// SetInt64ByIntIdx is
func (f *TabFile) SetInt64ByIntIdx(row int, col int, val int64) bool {
	return f.SetCell(row, col, strconv.FormatInt(val, 10))
}

// SetInt64ByStrIdx is
func (f *TabFile) SetInt64ByStrIdx(row string, col string, val int64) bool {
	return f.SetInt64ByIntIdx(f.FindRow(row), f.FindCol(col), val)
}

// SetInt64ByMixIdx is
func (f *TabFile) SetInt64ByMixIdx(row int, col string, val int64) bool {
	return f.SetInt64ByIntIdx(row, f.FindCol(col), val)
}

// SetStrByIntIdx is
func (f *TabFile) SetStrByIntIdx(row int, col int, val string) bool {
	return f.SetCell(row, col, val)
}

// SetStrByStrIdx is
func (f *TabFile) SetStrByStrIdx(row string, col string, val string) bool {
	return f.SetStrByIntIdx(f.FindRow(row), f.FindCol(col), val)
}

// SetStrByMixIdx is
func (f *TabFile) SetStrByMixIdx(row int, col string, val string) bool {
	return f.SetStrByIntIdx(row, f.FindCol(col), val)
}

// SetFloat32ByIntIdx is
func (f *TabFile) SetFloat32ByIntIdx(row int, col int, val float32) bool {
	return f.SetCell(row, col, strconv.FormatFloat(float64(val), 'f', 30, 32))
}

// SetFloat32ByStrIdx is
func (f *TabFile) SetFloat32ByStrIdx(row string, col string, val float32) bool {
	return f.SetFloat32ByIntIdx(f.FindRow(row), f.FindCol(col), val)
}

// SetFloat32ByMixIdx is
func (f *TabFile) SetFloat32ByMixIdx(row int, col string, val float32) bool {
	return f.SetFloat32ByIntIdx(row, f.FindCol(col), val)
}

// SetFloat64ByIntIdx is
func (f *TabFile) SetFloat64ByIntIdx(row int, col int, val float64) bool {
	return f.SetCell(row, col, strconv.FormatFloat(val, 'f', 62, 64))
}

// SetFloat64ByStrIdx is
func (f *TabFile) SetFloat64ByStrIdx(row string, col string, val float64) bool {
	return f.SetFloat64ByIntIdx(f.FindRow(row), f.FindCol(col), val)
}

// SetFloat64ByMixIdx is
func (f *TabFile) SetFloat64ByMixIdx(row int, col string, val float64) bool {
	return f.SetFloat64ByIntIdx(row, f.FindCol(col), val)
}

// SetCell is
func (f *TabFile) SetCell(row int, col int, val string) bool {
	if row >= f.rows || col >= f.cols || row < 0 || col < 0 {
		return false
	}

	idx := row*f.cols + col
	f.tabs[idx].content = val
	return true
}

func (f *TabFile) createTabOffsets(buff []byte, size int) {
	f.getRowsAndColumns(buff, size)
	f.createTabOffsetLinks(buff, size)
}

func (f *TabFile) createTabOffsetLinks(buff []byte, size int) {
	var offset, length, start, end int
	var str string

	for i := 0; i < f.rows; i++ {
		for j := 0; j < f.cols; j++ {
			idx := i*f.cols + j

			// 读取数据
			start = offset
			length = 0
			for offset < size {
				v := buff[offset]
				if v == '\t' || v == '\r' || v == '\n' {
					break
				}
				offset++
				length++
			}
			end = start + length
			str = string(buff[start:end])
			f.tabs[idx] = newTabCell(str)

			// 跳过tab
			v := buff[offset]
			if v == '\t' {
				offset++
				continue
			}

			if v == '\r' || v == '\n' {
				// 如果没有填满本行
				for k := j + 1; k < f.cols; k++ {
					idx := i*f.cols + k
					f.tabs[idx] = newTabCell("")
				}

				// 读到行尾，跳过本行的所有换行回车
				if offset+2 < size && buff[offset] == 0x0D && buff[offset+1] == 0x0A {
					// \r\n
					offset += 2
				} else { // \r or \n
					offset++
				}

				break
			}
		}
	}
}

func (f *TabFile) getRowsAndColumns(buff []byte, size int) {
	var rows, cols, offset int

	for offset < size {
		v := buff[offset]

		if v == '\t' { // '\t' = 0x09
			cols++
		} else if v == '\r' || v == '\n' { // '\r' = 0x0D, '\n' = 0x0A
			cols++
			rows++
			f.cols = Max(f.cols, cols)
			cols = 0

			if offset+2 <= size && buff[offset] == '\r' && buff[offset+1] == '\n' { // '\r\n'
				offset += 2
				continue
			}
		}

		offset++
		continue
	}

	f.rows = rows
	f.tabs = make([]*tabCell, f.rows*f.cols)
}
