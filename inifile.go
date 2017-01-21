package goblazer

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// IniFileKeyNode :
type IniFileKeyNode struct {
	ID    uint32 // id of key-value node
	Name  string // key name of key-value node
	Value string // key value of key-value node
}

func newIniFileKeyNode(id uint32, name string, value string) *IniFileKeyNode {
	n := new(IniFileKeyNode)
	n.ID = id
	n.Name = name
	n.Value = value
	return n
}

// IniFileSecNode :
type IniFileSecNode struct {
	ID       uint32                     // id of section
	Name     string                     // section name
	KeyNodes map[uint32]*IniFileKeyNode // all keys of a section
}

func newIniFileSecNode(id uint32, name string) *IniFileSecNode {
	n := new(IniFileSecNode)
	n.ID = id
	n.Name = name
	n.KeyNodes = make(map[uint32]*IniFileKeyNode)
	return n
}

// IniFile :
type IniFile struct {
	SecNodes map[uint32]*IniFileSecNode
	offset   int64
}

// NewIniFile :
func NewIniFile() *IniFile {
	f := new(IniFile)
	f.SecNodes = make(map[uint32]*IniFileSecNode)
	return f
}

// Load :
func (f *IniFile) Load(filePath string) bool {
	var ok bool
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

	size = int64(len(buff))
	if info.Size() != size {
		return false
	}

	f.createLinks(buff, size)
	return true
}

// PrintSections :
func (f *IniFile) PrintSections() {
	for secID, secNode := range f.SecNodes {
		fmt.Printf("[Section Node %d] ID = %d, Name = %s\n", secID, secNode.ID, secNode.Name)

		for keyID, keyNode := range f.SecNodes[secID].KeyNodes {
			fmt.Printf("<Section Node %d> ID = %d, Name = %s, Value = %s\n", keyID, keyNode.ID, keyNode.Name, keyNode.Value)
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
		id = SimpleHashString2ID(&key)
		if keyNode, ok := secNode.KeyNodes[id]; ok {
			keyNode.Value = ""
		}
	}
}

// RemoveKey :
func (f *IniFile) RemoveKey(sec string, key string) {
	id := f.formatSectionName(&sec)
	if secNode, ok := f.SecNodes[id]; ok {
		id = SimpleHashString2ID(&key)
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

// GetInt :
func (f *IniFile) GetInt(sec string, key string, dflt int) int {
	if s, ok := f.getKeyValue(sec, key); ok {
		if v, err := strconv.Atoi(s); err == nil {
			return v
		}
	}
	return dflt
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

// GetInt64 :
func (f *IniFile) GetInt64(sec string, key string, dflt int64) int64 {
	if s, ok := f.getKeyValue(sec, key); ok {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			return v
		}
	}
	return dflt
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
		secNode = newIniFileSecNode(id, sec)
		f.SecNodes[id] = secNode
	} else {
		// 如果Section Node存在
	}

	// 查找对应的Key Node
	id = SimpleHashString2ID(&key)
	if keyNode, ok = secNode.KeyNodes[id]; !ok {
		// 如果Key Node不存在
		keyNode = newIniFileKeyNode(id, key, val)
		secNode.KeyNodes[id] = keyNode
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
		*sec = JoinStrings([]string{"]", *sec})
	}

	return SimpleHashString2ID(sec)
}

func (f *IniFile) getKeyValue(sec string, key string) (string, bool) {
	id := f.formatSectionName(&sec)
	if secNode, ok := f.SecNodes[id]; ok {
		if keyNode, ok := secNode.KeyNodes[id]; ok {
			return keyNode.Value, true
		}
	}

	return "", false
}
