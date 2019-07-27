package midi

import (
	"errors"
	"fmt"
	"log"
)

// ChunkParser is a midi file parser. It can be used to parse the whole midi file.
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

// ParserData contains data read from midiparse.yaml
type ParserData struct {
	ControlChange []struct {
		Controller byte   `yaml:"control"`
		Message    string `yaml:"message"`
		//Values     uint   `yaml:"values"`
	} `yaml:"controlChange,flow"`

	ChannelVoiceMessage []struct {
		StatusByte byte   `yaml:"statusByte"`
		NValues    uint   `yaml:"nvalues"`
		Message    string `yaml:"message"`
	} `yaml:"channelVoiceMessage,flow"`
}

// NewChunkParser creates a new instance of ChunkParser.
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

// ParseTop parses the whole midi file.
func (cp *ChunkParser) ParseTop() []MidiElement {
	if cp.Chunk.Type == MThd {
		ret := make([]MidiElement, 1)
		parsed := cp.ParseMThd()
		ret = append(ret, parsed)
		return ret
	}
	panic("Not implemented!")
}

// ParseMThd parses the header part of a midi file.
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

// ParseMTrk parses the track part of a midi file.
// It stops when the current track is finished.
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

// ParseEventPair parses one deltatime-event pair.
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

// ParseEvent parses the event part of a deltatime-event pair.
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
		//fmt.Printf("New running status %X\n", head)
		event, err = cp.parseMidiEventWithRunningStatus(head)
	} else { // Midi events without running status
		//fmt.Printf("No new running status %X\n", head)
		event, err = cp.parseMidiEventWithoutRunningStatus(head)
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
	}
	text := string(cp.nextNByte(length))
	return length, text
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
		control := cp.nextNByte(2)
		message, _, err := cp.getControlChangeByController(control[0])
		if err != nil {
			event = &MidiEvent{
				runningStatus,
				"Unknown event",
				channel,
				control,
			}
			cp.registerRunningStatus(runningStatus)
		} else {
			event = &MidiEvent{
				runningStatus,
				message,
				channel,
				control,
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
			cp.registerRunningStatus(runningStatus)
		} else {
			data := cp.nextNByte(nvalues)
			event = &MidiEvent{
				runningStatus,
				message,
				channel,
				data,
			}
			cp.registerRunningStatus(runningStatus)
		}
	}
	return event, nil
}

func (cp *ChunkParser) parseMidiEventWithoutRunningStatus(firstData byte) (Event, error) {
	runningStatus := cp.runningStatus.statusByte
	channel := runningStatus & 0x0F
	var event Event
	if 0xB0 <= runningStatus && runningStatus <= 0xBF { // control change
		secondData := cp.nextByte()
		control := []byte{firstData, secondData}
		message, _, err := cp.getControlChangeByController(control[0])
		if err != nil {
			event = &MidiEvent{
				runningStatus,
				"Unknown event",
				channel,
				control,
			}
			cp.registerRunningStatus(runningStatus)
		} else {
			event = &MidiEvent{
				runningStatus,
				message,
				channel,
				control,
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
			restData := cp.nextNByte(nvalues - 1)
			data := []byte{firstData}
			data = append(data, restData...)
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

func (cp *ChunkParser) registerRunningStatus(statusByte byte) {
	cp.runningStatus.statusByte = statusByte
}

func (cp *ChunkParser) getControlChangeByController(controller byte) (string, uint, error) {
	for _, v := range cp.ParserData.ControlChange {
		if v.Controller == controller {
			return v.Message, 2, nil
		}
	}
	return "", 0, errors.New("Unknown event")
}

func (cp *ChunkParser) getChannelVoiceMessageByStatusByte(statusByte byte) (string, uint, error) {
	for _, v := range cp.ParserData.ChannelVoiceMessage {
		//fmt.Printf("%d", v.StatusByte)
		if v.StatusByte <= statusByte && statusByte < v.StatusByte+16 {
			return v.Message, v.NValues, nil
		}
	}
	return "", 0, errors.New("Unknown event")
}
