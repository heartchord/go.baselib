package goblazer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/henrylee2cn/mahonia"
)

// IniFileKeyNode :
type IniFileKeyNode struct {
	ID    uint32 // id of key-value node
	Name  string // key name of key-value node
	Value string // key value of key-value node
	SeqNo uint32 // sequence number
}

func newIniFileKeyNode(id uint32, name string, value string, no uint32) *IniFileKeyNode {
	n := new(IniFileKeyNode)
	n.ID = id
	n.Name = name
	n.Value = value
	n.SeqNo = no
	return n
}

// KeyNodesMap :
type KeyNodesMap map[uint32]*IniFileKeyNode

// KeyNodesMapSorter :
type KeyNodesMapSorter []*IniFileKeyNode

// NewKeyNodesMapSorter :
func NewKeyNodesMapSorter(m KeyNodesMap) KeyNodesMapSorter {
	ms := make(KeyNodesMapSorter, 0, len(m))

	for _, v := range m {
		ms = append(ms, v)
	}

	return ms
}

func (ms KeyNodesMapSorter) Len() int {
	return len(ms)
}

func (ms KeyNodesMapSorter) Less(i, j int) bool {
	return ms[i].SeqNo < ms[j].SeqNo
}

func (ms KeyNodesMapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

// IniFileSecNode :
type IniFileSecNode struct {
	ID       uint32      // id of section
	Name     string      // section name
	SeqNo    uint32      // sequence number
	KeyNodes KeyNodesMap // all keys of a section
}

func newIniFileSecNode(id uint32, name string, no uint32) *IniFileSecNode {
	n := new(IniFileSecNode)
	n.ID = id
	n.Name = name
	n.SeqNo = no
	n.KeyNodes = make(map[uint32]*IniFileKeyNode)
	return n
}

// SecNodesMap :
type SecNodesMap map[uint32]*IniFileSecNode

// SecNodesMapSorter :
type SecNodesMapSorter []*IniFileSecNode

// NewSecNodesMapSorter :
func NewSecNodesMapSorter(m SecNodesMap) SecNodesMapSorter {
	ms := make(SecNodesMapSorter, 0, len(m))

	for _, v := range m {
		ms = append(ms, v)
	}

	return ms
}

func (ms SecNodesMapSorter) Len() int {
	return len(ms)
}

func (ms SecNodesMapSorter) Less(i, j int) bool {
	return ms[i].SeqNo < ms[j].SeqNo
}

func (ms SecNodesMapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

// IniFile :
type IniFile struct {
	SecNodes     SecNodesMap
	offset       int64
	seqNoCounter uint32
}

// NewIniFile :
func NewIniFile() *IniFile {
	f := new(IniFile)
	f.SecNodes = make(map[uint32]*IniFileSecNode)
	return f
}

// Load :
func (f *IniFile) Load(filePath string, code string) bool {
	var ok, isutf8 bool
	var fi *os.File
	var err error
	var buff []byte
	var info os.FileInfo
	var size int64

	if info, err = os.Stat(filePath); err != nil {
		return false
	}

	if ok = IsFileExisted(filePath); !ok {
		return false
	}

	if fi, err = os.Open(filePath); err != nil {
		return false
	}
	defer fi.Close()

	if buff, err = ioutil.ReadAll(fi); err != nil {
		return false
	}

	isutf8 = strings.EqualFold(code, "utf8")
	if !isutf8 {
		code = strings.ToUpper(code)
		mdecoder := mahonia.NewDecoder(code)
		buff = []byte(mdecoder.ConvertString(string(buff)))
	}

	size = int64(len(buff))
	if isutf8 && info.Size() != size {
		return false
	}

	f.createLinks(buff, size)
	return true
}

// Save :
func (f *IniFile) Save(filePath string, code string) bool {
	var ok, isutf8 bool
	var err error
	var buff bytes.Buffer
	var str string
	var fi *os.File

	secms := NewSecNodesMapSorter(f.SecNodes)
	sort.Sort(secms)

	for _, sec := range secms {
		buff.WriteString(fmt.Sprintf("%s\r\n", sec.Name))

		keyms := NewKeyNodesMapSorter(sec.KeyNodes)
		sort.Sort(keyms)

		for _, key := range keyms {
			buff.WriteString(fmt.Sprintf("%s%s%s\r\n", key.Name, "=", key.Value))
		}

		buff.WriteString("\r\n")
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

	if fi, err = os.Create(filePath); err != nil {
		return false
	}
	defer fi.Close()

	if _, err = fi.WriteString(str); err != nil {
		return false
	}

	return true
}

// PrintSections :
func (f *IniFile) PrintSections() {
	for secID, secNode := range f.SecNodes {
		fmt.Printf("[Section Node %d] ID = %d, Name = %s\n", secID, secNode.ID, secNode.Name)

		for keyID, keyNode := range f.SecNodes[secID].KeyNodes {
			fmt.Printf("<Key Node %d> ID = %d, Name = %s, Value = %s\n", keyID, keyNode.ID, keyNode.Name, keyNode.Value)
		}
	}
}

// Clear :
func (f *IniFile) Clear() {
	f.removeLinks()
}

// IsSectionExisted :
func (f *IniFile) IsSectionExisted(sec string) bool {
	id := f.formatSectionName(&sec)
	_, ok := f.SecNodes[id]
	return ok
}

// ClearSection :
func (f *IniFile) ClearSection(sec string) {
	var ok bool
	var secNode *IniFileSecNode

	id := f.formatSectionName(&sec)
	if secNode, ok = f.SecNodes[id]; !ok {
		return
	}

	secNode.KeyNodes = make(map[uint32]*IniFileKeyNode)
}

// RemoveSection :
func (f *IniFile) RemoveSection(sec string) {
	id := f.formatSectionName(&sec)
	if _, ok := f.SecNodes[id]; ok {
		delete(f.SecNodes, id)
	}
}

// ClearKey :
func (f *IniFile) ClearKey(sec string, key string) {
	id := f.formatSectionName(&sec)
	if secNode, ok := f.SecNodes[id]; ok {
		id = SimpleHashString2ID(key)
		if keyNode, ok := secNode.KeyNodes[id]; ok {
			keyNode.Value = ""
		}
	}
}

// RemoveKey :
func (f *IniFile) RemoveKey(sec string, key string) {
	id := f.formatSectionName(&sec)
	if secNode, ok := f.SecNodes[id]; ok {
		id = SimpleHashString2ID(key)
		if _, ok := secNode.KeyNodes[id]; ok {
			delete(secNode.KeyNodes, id)
		}
	}
}

// GetString :
func (f *IniFile) GetString(sec string, key string, dflt string) string {
	if s, ok := f.getKeyValue(sec, key); ok {
		return s
	}
	return dflt
}

// SetString :
func (f *IniFile) SetString(sec string, key string, val string) bool {
	return f.setKeyValue(sec, key, val)
}

// GetStrings :
func (f *IniFile) GetStrings(sec string, key string, sep string) []string {
	var ret []string
	if s, ok := f.getKeyValue(sec, key); ok {
		ret = strings.Split(s, sep)
	}
	return ret
}

// SetStrings :
func (f *IniFile) SetStrings(sec string, key string, val []string, sep string) bool {
	s := strings.Join(val, sep)
	return f.setKeyValue(sec, key, s)
}

// GetInt :
func (f *IniFile) GetInt(sec string, key string, dflt int) int {
	if s, ok := f.getKeyValue(sec, key); ok {
		if v, err := strconv.Atoi(s); err == nil {
			return v
		}
	}
	return dflt
}

// SetInt :
func (f *IniFile) SetInt(sec string, key string, val int) bool {
	s := strconv.Itoa(val)
	return f.setKeyValue(sec, key, s)
}

// GetInts :
func (f *IniFile) GetInts(sec string, key string, sep string) []int {
	var ret []int
	if s, ok := f.getKeyValue(sec, key); ok {
		strs := strings.Split(s, sep)
		ret = StrSliceToIntSlice(strs)
	}
	return ret
}

// SetInts :
func (f *IniFile) SetInts(sec string, key string, val []int, sep string) bool {
	strs := IntSliceToStrSlice(val)
	s := strings.Join(strs, sep)
	return f.setKeyValue(sec, key, s)
}

// GetInt32 :
func (f *IniFile) GetInt32(sec string, key string, dflt int32) int32 {
	if s, ok := f.getKeyValue(sec, key); ok {
		if v, err := strconv.ParseInt(s, 10, 32); err == nil {
			return int32(v)
		}
	}
	return dflt
}

// SetInt32 is
func (f *IniFile) SetInt32(sec string, key string, val int32) bool {
	s := strconv.FormatInt(int64(val), 10)
	return f.setKeyValue(sec, key, s)
}

// GetInt32s :
func (f *IniFile) GetInt32s(sec string, key string, sep string) []int32 {
	var ret []int32
	if s, ok := f.getKeyValue(sec, key); ok {
		strs := strings.Split(s, sep)
		ret = StrSliceToInt32Slice(strs)
	}
	return ret
}

// SetInt32s :
func (f *IniFile) SetInt32s(sec string, key string, val []int32, sep string) bool {
	strs := Int32SliceToStrSlice(val)
	s := strings.Join(strs, sep)
	return f.setKeyValue(sec, key, s)
}

// GetInt64 :
func (f *IniFile) GetInt64(sec string, key string, dflt int64) int64 {
	if s, ok := f.getKeyValue(sec, key); ok {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			return v
		}
	}
	return dflt
}

// SetInt64 is
func (f *IniFile) SetInt64(sec string, key string, val int64) bool {
	s := strconv.FormatInt(val, 10)
	return f.setKeyValue(sec, key, s)
}

// GetInt64s :
func (f *IniFile) GetInt64s(sec string, key string, sep string) []int64 {
	var ret []int64
	if s, ok := f.getKeyValue(sec, key); ok {
		strs := strings.Split(s, sep)
		ret = StrSliceToInt64Slice(strs)
	}
	return ret
}

// SetInt64s :
func (f *IniFile) SetInt64s(sec string, key string, val []int64, sep string) bool {
	strs := Int64SliceToStrSlice(val)
	s := strings.Join(strs, sep)
	return f.setKeyValue(sec, key, s)
}

// GetFloat32 :
func (f *IniFile) GetFloat32(sec string, key string, dflt float32) float32 {
	if s, ok := f.getKeyValue(sec, key); ok {
		if v, err := strconv.ParseFloat(s, 32); err == nil {
			return float32(v)
		}
	}
	return dflt
}

// SetFloat32 is
func (f *IniFile) SetFloat32(sec string, key string, val float32) bool {
	s := strconv.FormatFloat(float64(val), 'f', 30, 32)
	return f.setKeyValue(sec, key, s)
}

// GetFloat32s :
func (f *IniFile) GetFloat32s(sec string, key string, sep string) []float32 {
	var ret []float32
	if s, ok := f.getKeyValue(sec, key); ok {
		strs := strings.Split(s, sep)
		ret = StrSliceToFloat32Slice(strs)
	}
	return ret
}

// SetFloat32s :
func (f *IniFile) SetFloat32s(sec string, key string, val []float32, sep string) bool {
	strs := Float32SliceToStrSlice(val)
	s := strings.Join(strs, sep)
	return f.setKeyValue(sec, key, s)
}

// GetFloat64 :
func (f *IniFile) GetFloat64(sec string, key string, dflt float64) float64 {
	if s, ok := f.getKeyValue(sec, key); ok {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return v
		}
	}
	return dflt
}

// SetFloat64 is
func (f *IniFile) SetFloat64(sec string, key string, val float64) bool {
	s := strconv.FormatFloat(val, 'f', 62, 64)
	return f.setKeyValue(sec, key, s)
}

// GetFloat64s :
func (f *IniFile) GetFloat64s(sec string, key string, sep string) []float64 {
	var ret []float64
	if s, ok := f.getKeyValue(sec, key); ok {
		strs := strings.Split(s, sep)
		ret = StrSliceToFloat64Slice(strs)
	}
	return ret
}

// SetFloat64s :
func (f *IniFile) SetFloat64s(sec string, key string, val []float64, sep string) bool {
	strs := Float64SliceToStrSlice(val)
	s := strings.Join(strs, sep)
	return f.setKeyValue(sec, key, s)
}

// GetBool :
func (f *IniFile) GetBool(sec string, key string, dflt bool) bool {
	if s, ok := f.getKeyValue(sec, key); ok {
		return IsTrueString(s)
	}
	return dflt
}

// SetBool :
func (f *IniFile) SetBool(sec string, key string, val bool) bool {
	s := GetBoolString(val)
	return f.setKeyValue(sec, key, s)
}

// GetBools :
func (f *IniFile) GetBools(sec string, key string, sep string) []bool {
	var ret []bool
	if s, ok := f.getKeyValue(sec, key); ok {
		strs := strings.Split(s, sep)
		ret = StrSliceToBoolSlice(strs)
	}
	return ret
}

// SetBools :
func (f *IniFile) SetBools(sec string, key string, val []bool, sep string) bool {
	strs := BoolSliceToStrSlice(val)
	s := strings.Join(strs, sep)
	return f.setKeyValue(sec, key, s)
}

func (f *IniFile) createLinks(buff []byte, size int64) {
	var ok bool
	var len, start, end int64
	var str, sec, key, val string

	// 清空缓冲偏移
	f.offset = 0

	for f.offset < size {
		start = f.offset

		len = f.readLine(buff, size)
		if len < 0 { // 文件读完
			break
		}

		if len == 0 { // 空行
			continue
		}

		end = start + len
		str = strings.TrimSpace(string(buff[start:end]))
		if f.isKeyChar(str[0]) { // key - value
			if sec == "" { // 没有sec忽略所有
				continue
			}

			if key, val, ok = f.splitKeyValue(str); ok {
				f.setKeyValue(sec, key, val)
			}
			continue
		}

		// section处理
		if str[0] == '[' {
			sec = str
		}
	}
}

func (f *IniFile) removeLinks() {
	f.SecNodes = make(map[uint32]*IniFileSecNode)
	f.seqNoCounter = 0
}

func (f *IniFile) readLine(buff []byte, size int64) int64 {
	var len int64

	if f.offset >= size { // 文件读完
		return -1
	}

	for f.offset < size && buff[f.offset] != 0x0D && buff[f.offset] != 0x0A {
		f.offset++
		len++
	}

	if f.offset+2 <= size && buff[f.offset] == 0x0D && buff[f.offset+1] == 0x0A {
		// windows
		f.offset += 2
	} else if f.offset < size {
		// linux
		f.offset++
	}

	return len
}

func (f *IniFile) isKeyChar(ch byte) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}

	if ch >= 'A' && ch <= 'Z' {
		return true
	}

	if ch >= 'a' && ch <= 'z' {
		return true
	}

	return false
}

func (f *IniFile) splitKeyValue(s string) (string, string, bool) {
	i := strings.IndexAny(s, "=")

	if i < 0 { // 没有找到分隔符
		return "", "", false
	}

	key := strings.TrimSpace(s[0:i])
	val := strings.TrimSpace(s[i+1:])
	return key, val, true
}

func (f *IniFile) setKeyValue(sec string, key string, val string) bool {
	var id uint32
	var ok bool
	var secNode *IniFileSecNode
	var keyNode *IniFileKeyNode

	// 查找对应的Section Node
	id = f.formatSectionName(&sec)
	if secNode, ok = f.SecNodes[id]; !ok {
		// 如果Section Node不存在，创建一个
		secNode = newIniFileSecNode(id, sec, f.seqNoCounter)
		f.SecNodes[id] = secNode
		f.seqNoCounter++
	} else {
		// 如果Section Node存在
	}

	// 查找对应的Key Node
	id = SimpleHashString2ID(key)
	if keyNode, ok = secNode.KeyNodes[id]; !ok {
		// 如果Key Node不存在
		keyNode = newIniFileKeyNode(id, key, val, f.seqNoCounter)
		secNode.KeyNodes[id] = keyNode
		f.seqNoCounter++
	} else {
		// 如果Key Node存在，覆盖旧值
		secNode.KeyNodes[id].Value = val
	}

	return true
}

func (f *IniFile) formatSectionName(sec *string) uint32 {
	if (*sec)[0] != '[' {
		*sec = JoinStrings([]string{"[", *sec})
	}

	l := len(*sec)
	if (*sec)[l-1] != ']' {
		*sec = JoinStrings([]string{*sec, "]"})
	}

	return SimpleHashString2ID(*sec)
}

func (f *IniFile) getKeyValue(sec string, key string) (string, bool) {
	id := f.formatSectionName(&sec)
	if secNode, ok := f.SecNodes[id]; ok {
		id := SimpleHashString2ID(key)
		if keyNode, ok := secNode.KeyNodes[id]; ok {
			return keyNode.Value, true
		}
	}

	return "", false
}
