package pkg

func (this *Smooth) add(b1 byte) {
	this.data = append([]byte{b1}, this.data...)
}
func (this *Smooth) pop() uint8 {
	if 0 < len(this.data) {
		b1 := this.data[0]
		this.data = this.data[1:]
		this.nPos += 1
		return b1
	}
	b, _ := this._p.readUInt8()
	return b
}

func (this *Smooth) size() int {
	return this._p.maxDataBlockSize - this.nPos
}
func (this *Smooth) peek() uint8 {
	if 0 < len(this.data) {
		return this.data[0]
	}
	b, _ := this._p.readUInt8()
	this.add(b)

	return b
}
