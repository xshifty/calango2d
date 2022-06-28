package engine

type AlphaColor uint64

func NewAlphaColor(r, g, b, a uint8) AlphaColor {
	return AlphaColor(uint32(a) + (uint32(b) << 8) + (uint32(g) << 16) + (uint32(r) << 24))
}

func (c *AlphaColor) GetRed() uint8 {
	return uint8(uint64(*c)>>24) & 255
}

func (c *AlphaColor) GetGreen() uint8 {
	return uint8(uint64(*c)>>16) & 255
}

func (c *AlphaColor) GetBlue() uint8 {
	return uint8(uint64(*c)>>8) & 255
}

func (c *AlphaColor) GetAlpha() uint8 {
	return uint8(*c) & 255
}
