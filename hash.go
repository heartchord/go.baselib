package goblazer

// SimpleHashString2ID :
func SimpleHashString2ID(s string) uint32 {
	if ok := CheckStringEmpty(s); !ok {
		return 0xFFFFFFFF
	}

	var id uint32
	var i, l int32
	for l = int32(len(s)); i < l; i++ {
		c1 := int8(s[i])
		c2 := (i + 1) * int32(c1)
		id = (id + uint32(c2)) % 0x8000000b * 0xffffffef
	}

	return (id ^ 0x12345678)
}
