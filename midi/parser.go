package midi

import (
	"errors"
	"fmt"
	"log"
)

const (
	CHANNELVOICEMESSAGE = iota
	CHANNELMODEMESSAGE
)

type ChunkParser struct {
	Chunk         *Chunk
	pointer       int
	ParserData    ParserData
	runningStatus runningStatus
	absoluteTime  uint32
}

type runningStatus struct {
	statusByte byte
	//message     string
	//messageType int
	//channel     byte
	//nValues     uint
}

type ParserData struct {
	ChannelModeMessage []struct {
		Controller byte   `yaml:"control"`
		Message    string `yaml:"message"`
		//Values     uint   `yaml:"values"`
	} `yaml:"channelModeMessage,flow"`

	ChannelVoiceMessage []struct {
		StatusByte byte   `yaml:"statusByte"`
		NValues    uint   `yaml:"nvalues"`
		Message    string `yaml:"message"`
	} `yaml:"channelVoiceMessagemflow"`
}

func NewChunkParser(chunk *Chunk) *ChunkParser {
	return &ChunkParser{
		Chunk:   chunk,
		pointer: 0,
	}
}

func (cp *ChunkParser) nextNByte(length uint) []byte {
	leng := int(length)
	if len(cp.Chunk.Data) < cp.pointer+leng {
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
	var ans uint
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

func (cp *ChunkParser) ParseMTrk(deltaTime bool) []EventPair {
	ret := make([]EventPair, 32)
	for {
		evpair, err := cp.ParseEventPair()
		if err != nil {
			log.Printf("Error Mes: %s", err)
			break
		}
		if !deltaTime {
			cp.absoluteTime += evpair.DeltaTime
			evpair.DeltaTime = cp.absoluteTime
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
		return EventPair{}, err
	}
	//log.Printf("%v\n", event)
	return EventPair{uint32(deltaTime), event}, nil
}

func (cp *ChunkParser) ParseEvent() (Event, error) {
	head := cp.nextNByte(1)[0]
	var event Event
	var err error
	if head == 0xFF { // When meta events
		metaType := cp.nextByte()
		switch metaType {
		case 0x51:
			event, err = cp.parseMetaEventSetTempo()
		case 0x03:
			event, err = cp.parseMetaEventSequenceTrackName()
		case 0x58:
			event, err = cp.parseMetaEventTimeSignature()
		case 0x06:
			event, err = cp.parseMetaEventMarker()
		case 0x2f:
			event, err = cp.parseMetaEventEndOfTrack()
		default:
			msg := fmt.Sprintf("unknown meta event: FF %X", metaType)
			event, err = nil, errors.New(msg)
		}
	} else if head >= 0x80 { // Midi events with running status
		event, err = cp.parseMidiEventWithRunningStatus(head)
	} else { // Midi events without running status
		event, err = cp.parseMidiEventWithoutRunningStatus()
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

func (cp *ChunkParser) parseMetaEventLengthTextHelper() (uint, string) {
	length := uint(cp.nextByte())
	if length == 0 {
		return length, ""
	} else {
		text := string(cp.nextNByte(length))
		return length, text
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

func (cp *ChunkParser) parseMetaEventSequenceTrackName() (Event, error) {
	_, text := cp.parseMetaEventLengthTextHelper()
	event := MetaEvent{3, &TrackName{text}}
	return &event, nil
}

func (cp *ChunkParser) parseMetaEventMarker() (Event, error) {
	_, text := cp.parseMetaEventLengthTextHelper()
	event := MetaEvent{6, &Marker{text}}
	return &event, nil
}

func (cp *ChunkParser) parseMetaEventEndOfTrack() (Event, error) {
	check := cp.nextByte()
	if check != byte(0) {
		return &MetaEvent{0x2f, &EndOfTrack{}}, errors.New("invalid End of Track event")
	}
	return &MetaEvent{0x2f, &EndOfTrack{}}, nil
}

func (cp *ChunkParser) parseMidiEventWithRunningStatus(runningStatus byte) (Event, error) {
	channel := runningStatus & 0x0F
	var event Event
	if 0xB0 <= runningStatus && runningStatus <= 0xBF { // channel mode messages
		control := cp.nextByte()
		message, _, err := cp.getChannelModeMessageByController(control)
		if err != nil {
			event = &MidiEvent{
				runningStatus,
				"Unknown event",
				channel,
				[]byte{},
			}
			cp.registerRunningStatus(runningStatus)
		} else {
			event = &MidiEvent{
				runningStatus,
				message,
				channel,
				[]byte{control},
			}
			cp.registerRunningStatus(runningStatus)
		}
	} else { // channel voice message
		message, nvalues, err := cp.getChannelVoiceMessageByStatusByte(runningStatus)
		channel := runningStatus & 0x0F
		if err != nil {
			event = &MidiEvent{
				runningStatus,
				"Unknown event",
				channel,
				[]byte{},
			}
		} else {
			data := cp.nextNByte(nvalues)
			event = &MidiEvent{
				runningStatus,
				message,
				channel,
				data,
			}
		}
	}
	return event, nil
}

func (cp *ChunkParser) parseMidiEventWithoutRunningStatus() (Event, error) {
	return cp.parseMidiEventWithRunningStatus(cp.runningStatus.statusByte)
}

func (cp *ChunkParser) registerRunningStatus(control byte) {
	cp.runningStatus.statusByte = control
}

func (cp *ChunkParser) getChannelModeMessageByController(controller byte) (string, uint, error) {
	for _, v := range cp.ParserData.ChannelModeMessage {
		if v.Controller == controller {
			return v.Message, 2, nil
		}
	}
	return "", 0, errors.New("Unknown event")
}

func (cp *ChunkParser) getChannelVoiceMessageByStatusByte(statusByte byte) (string, uint, error) {
	for _, v := range cp.ParserData.ChannelVoiceMessage {
		if v.StatusByte == statusByte {
			return v.Message, v.NValues, nil
		}
	}
	return "", 0, errors.New("Unknown event")
}
