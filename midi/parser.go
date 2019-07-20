package midi

import (
	"errors"
	"log"
)

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

func (cp *ChunkParser) nextNByte(length uint) []byte {
	leng := int(length)
	if len(cp.Chunk.Data) <= cp.pointer+leng {
		return nil
	}
	ret := cp.Chunk.Data[cp.pointer : cp.pointer+leng]
	cp.pointer += leng
	return ret
}

func (cp *ChunkParser) nextByte() byte {
	next := cp.nextNByte(1)
	if next == nil {
		panic("TODO")
	}
	return next[0]
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

func (cp *ChunkParser) parseVLRep() (uint, error) {
	var ans uint = 0
	for {
		vr := cp.nextNByte(1)
		if vr == nil {
			return 0, errors.New("EOF")
		}
		v := vr[0]
		if v < 128 {
			ans += uint(v)
			break
		} else {
			ans += uint(v) - 128
		}
	}
	return ans, nil
}

func (cp *ChunkParser) ParseMTrk() []EventPair {
	ret := make([]EventPair, 32)
	for {
		evpair, err := cp.ParseEventPair()
		if err != nil {
			log.Printf("Error Mes: %s", err)
			break
		}
		log.Printf("Parsed: {%s}", evpair.String())
		ret = append(ret, evpair)
	}
	return ret
}

func (cp *ChunkParser) ParseEventPair() (EventPair, error) {
	deltaTime, err := cp.parseVLRep()
	if err != nil {
		return EventPair{}, errors.New("no more event pairs")
	}
	//log.Printf("DeltaTime: %d\n", deltaTime)
	event, err := cp.ParseEvent()
	if err != nil {
		return EventPair{}, errors.New("invalid event")
	}
	//log.Printf("%v\n", event)
	return EventPair{uint32(deltaTime), event}, nil
}

func (cp *ChunkParser) ParseEvent() (Event, error) {
	head := cp.nextNByte(1)[0]
	var event Event
	var err error
	if head == 0xFF { // When meta events
		metaType := cp.nextNByte(1)[0]
		switch metaType {
		case 0x51:
			event, err = cp.parseMetaEventSetTempo()
		case 0x03:
			event, err = cp.parseMetaEventSequenceTrackName()
		case 0x58:
			event, err = cp.parseMetaEventTimeSignature()
		default:
			event, err = nil, errors.New("invalid event")
		}
	} else {
		event, err = nil, errors.New("invalid event")
	}
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (cp *ChunkParser) parseMetaEventSetTempo() (Event, error) {
	cp.nextNByte(1)
	raw := cp.nextNByte(3)
	tempo := raw[0]<<2 + raw[1]<<1 + raw[2]
	event := MetaEvent{51, &SetTempo{int(tempo)}}
	return &event, nil
}

func (cp *ChunkParser) parseMetaEventSequenceTrackName() (Event, error) {
	length := uint(cp.nextNByte(1)[0])
	if length == 0 {
		event := MetaEvent{3, &TrackName{""}}
		return &event, nil
	} else {
		name := string(cp.nextNByte(length))
		event := MetaEvent{3, &TrackName{name}}
		return &event, nil
	}
}

func (cp *ChunkParser) parseMetaEventTimeSignature() (Event, error) {
	cp.nextNByte(1)
	numerator := cp.nextByte()
	denominator := cp.nextByte()
	clocks := cp.nextByte()
	notes := cp.nextByte()

	event := MetaEvent{58, &TimeSignature{numerator, denominator, clocks, notes}}
	return &event, nil
}
