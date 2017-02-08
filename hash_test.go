package goblazer

import "testing"

func Test_SimpleHashString2ID_1(t *testing.T) {
	b := []byte{0xB9, 0xFE, 0xB9, 0xFE} // "哈哈" in GBK code
	s := BytesToString(b)

	id1 := SimpleHashString2ID(s)
	id2 := uint32(3986945288)

	if id1 == id2 {
		t.Log("Test_SimpleHashString2ID_1 succeeded")
	} else {
		t.Error("Test_SimpleHashString2ID_1 failed")
	}
}

func Test_SimpleHashString2ID_2(t *testing.T) {
	s := "哈哈"

	id1 := SimpleHashString2ID(s)
	id2 := uint32(3369549050)

	if id1 == id2 {
		t.Log("Test_SimpleHashString2ID_2 succeeded")
	} else {
		t.Error("Test_SimpleHashString2ID_2 failed")
	}
}

func Test_SimpleHashString2ID_3(t *testing.T) {
	b1 := []byte{0xB9, 0xFE, 0xB9, 0xFE}       // "哈哈" in GBK code
	b2 := []byte{0xB9, 0xFE, 0xB9, 0xFE, 0x61} // "哈哈a" in GBK code

	s1 := BytesToString(b1)
	s2 := BytesToString(b2)

	id1 := SimpleHashString2ID(s1)
	id2 := SimpleHashString2ID(s2)

	if id1 != id2 {
		t.Log("Test_SimpleHashString2ID_3 succeeded")
	} else {
		t.Error("Test_SimpleHashString2ID_3 failed")
	}
}

func Test_SimpleHashString2ID_4(t *testing.T) {
	s1 := "哈哈"
	s2 := "哈哈a"

	id1 := SimpleHashString2ID(s1)
	id2 := SimpleHashString2ID(s2)

	if id1 != id2 {
		t.Log("Test_SimpleHashString2ID_4 succeeded")
	} else {
		t.Error("Test_SimpleHashString2ID_4 failed")
	}
}
