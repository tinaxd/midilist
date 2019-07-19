package midi

type ChunkParser struct {
	Chunk   *Chunk
	pointer int
}

func NewChunkParser(chunk *Chunk) *ChunkParser {
	return &ChunkParser{
		chunk,
		0,
	}
}

func (cp *ChunkParser) ParseTop() []MidiElement {
	if cp.Chunk.Type == MThd {
		ret := make([]MidiElement, 1)
		parsed := cp.ParseMThd()
		ret = append(ret, parsed)
		return ret
	} else {
		panic("Not implemented!")
	}
}

func (cp *ChunkParser) ParseMThd() MidiMeta {
	format := int(cp.Chunk.Data[0]<<1 + cp.Chunk.Data[1])
	tracks := int(cp.Chunk.Data[2]<<1 + cp.Chunk.Data[3])
	division := int(cp.Chunk.Data[4]<<1 + cp.Chunk.Data[5])
	return MidiMeta{
		format,
		tracks,
		division,
	}
}

func (cp *ChunkParser) parseVLRep() int {
	ans := 0
	for _, v := range cp.Chunk.Data {
		if v > 127 {
			ans += int(v)
			break
		} else {
			ans += int(v) - 128
		}
	}
	return ans
}
