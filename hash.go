package goblazer

// SimpleHashString2ID uses a simple hash algorithm to generate a string id.
func SimpleHashString2ID(s string) uint32 {
	l := len(s)

	if l <= 0 {
		return 0xFFFFFFFF
	}

	var id uint32
	for i := 0; i < l; i++ {
		c1 := int8(s[i])
		c2 := (int32(i) + 1) * int32(c1)
		id = (id + uint32(c2)) % 0x8000000b * 0xffffffef
	}

	return (id ^ 0x12345678)
}
