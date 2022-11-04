package pkg

func (this *Smooth) pop() uint8 {
	b, _ := this._p.readUInt8()
	return b
}

func (this *Smooth) peek() uint8 {
	b, _ := this._p.readUInt8()
	return b
}
