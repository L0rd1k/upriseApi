package bits

type Bits []byte //> pieces of peers

func (bitField Bits) hasBit(idx int) bool {
	byteIdx := idx / 8
	offset := byteIdx % 8
	if byteIdx < 0 || byteIdx >= len(bitField) {
		return false
	}
	return bitField[byteIdx]>>uint(7-offset)&1 != 0
}

func (bitField Bits) setBit(idx int) {
	byteIdx := idx / 8
	offset := idx % 8
	if byteIdx < 0 || byteIdx >= len(bitField) {
		return
	}
	bitField[byteIdx] |= 1 << uint(7-offset)
}
