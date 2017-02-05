package goblazer

import (
	"io/ioutil"
	"os"
	"strings"

	"fmt"

	"github.com/henrylee2cn/mahonia"
)

type tabOffset struct {
	offset uint32
	length uint32
}

// TabFile opens a "xxx.tab" file to resolve it.
type TabFile struct {
	rows    uint32
	cols    uint32
	buff    []byte
	size    uint32
	offsets []tabOffset
}

// Load is
func (f *TabFile) Load(path string, code string) bool {
	var ok bool
	var fi *os.File
	var err error

	if fi, err = os.Open(path); err != nil {
		fmt.Println(err)
		return false
	}
	defer fi.Close()

	if f.buff, err = ioutil.ReadAll(fi); err != nil {
		return false
	}

	if ok = strings.EqualFold(code, "utf8"); !ok {
		var s string
		mdecoder := mahonia.NewDecoder(strings.ToUpper(code))

		if s, ok = mdecoder.ConvertStringOK(string(f.buff)); !ok {
			return false
		}
		f.buff = []byte(s)
	}

	f.size = uint32(len(f.buff))
	if f.size > 0 {
		f.createTabOffsets()
	}

	return true
}

func (f *TabFile) createTabOffsets() {
	f.getRowsAndColumns()
	f.createTabOffsetLinks()
}

//linux,unix: \r\n
//windows : \n
//Mac OS ： \r

func (f *TabFile) createTabOffsetLinks() {
	offset := uint32(0)
	length := uint32(0)

	for i := uint32(0); i < f.rows; i++ {
		for j := uint32(0); j < f.cols; j++ {
			idx := i*f.cols + j
			f.offsets[idx].offset = offset

			length = 0
			for offset < f.size {
				v := f.buff[offset]

				if v == 0x09 || v == 0x0D || v == 0x0A {
					break
				}

				offset++
				length++

			}

			f.offsets[idx].length = length

			// 跳过tab
			v := f.buff[offset]
			if v == 0x09 {
				offset++
				continue
			}

			if v == 0x0D || v == 0x0A {
				// 读到行尾，跳过本行的所有换行回车
				if offset+2 < f.size && f.buff[offset] == 0x0D && f.buff[offset+1] == 0x0A {
					// \r\n
					offset += 2
				} else { // \r or \n
					offset++
				}

				// 如果没有填满本行
				for k := j + 1; k < f.cols; k++ {
					idx := i*f.cols + k
					f.offsets[idx].length = 0
					f.offsets[idx].offset = offset
				}
				break
			}
		}
	}

	fmt.Println(f.offsets)
}

func (f *TabFile) getRowsAndColumns() {
	offset := uint32(0)

	// 如果文件非空，默认至少有一行一列
	if f.size > 0 {
		f.rows = 1
		f.cols = 1
	}

	// 读第一行决定有多少列
	for offset < f.size {
		v := f.buff[offset]

		if v == 0x0D || v == 0x0A { // 回车符或换行符
			break
		}

		if v == 0x09 { // Tab符
			f.cols++
		}

		offset++
	}

	// 跳过第一行的回车换行
	if offset+2 <= f.size && f.buff[offset] == 0x0D && f.buff[offset+1] == 0x0A {
		// \r\n
		offset += 2
	} else if offset+1 <= f.size {
		// \n
		offset++
	} else {
		panic("Unexpected character!")
	}

	// 读取有多少行
	for offset < f.size {
		v := f.buff[offset]

		// 跳过非回车换行字符
		if v != 0x0D && v != 0x0A {
			offset++
			continue
		}

		// 行数增加
		f.rows++

		// 跳过回车换行
		if offset+2 <= f.size && f.buff[offset] == 0x0D && f.buff[offset+1] == 0x0A {
			// \r\n
			offset += 2
		} else if offset+1 <= f.size {
			// \n
			offset++
		} else {
			panic("Unexpected character")
		}
	}

	f.offsets = make([]tabOffset, f.rows*f.cols)

	fmt.Println(f.rows)
	fmt.Println(f.cols)
}
