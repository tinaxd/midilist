package midi

import "fmt"

type EventPair struct {
	DeltaTime uint32
	Event     Event
}

func (ep *EventPair) MidiElement() {}

func (ep *EventPair) String() string {
	return fmt.Sprintf("Time: %d, Event: %s", ep.DeltaTime, ep.Event)
}

/*
Event --- MidiEvent + MidiEventData
	   |
	   |- ... // TODO
*/

// General types

type Event interface {
	Event()
}

type MidiEventData interface {
	MidiEventData()
}

type SysexEventData interface {
	SysexEventData()
}

type MetaEventData interface {
	MetaEventData()
}

type MidiEvent struct {
	StatusByte byte
	Message    string
	Channel    byte
	Values     []byte
}

type SysexEvent struct {
	Type      byte
	IncludeF0 bool
	Data      SysexEventData
}

type MetaEvent struct {
	Type byte
	Data MetaEventData
}

type ChannelModeMessage struct {
	Channel    byte
	StatusByte byte
	Message    string
	Value      byte
}

func (event *MidiEvent) Event()          {}
func (event *SysexEvent) Event()         {}
func (event *MetaEvent) Event()          {}
func (event *ChannelModeMessage) Event() {}

// MidiEventData

type NoteOff struct {
	Channel  byte
	Key      byte
	Velocity byte
}

type NoteOn struct {
	Channel  byte
	Key      byte
	Velocity byte
}

type PolyphonicKeyPressure struct {
	Channel  byte
	Key      byte
	Pressure byte
}

type ControllerChange struct {
	Channel    byte
	Controller byte
	Value      byte
}

type ProgramChange struct {
	Channel byte
	Program byte
}

type ChannelKeyPressure struct {
	Channel  byte
	Pressure byte
}

type PitchBend struct {
	Channel byte
	Lsb     byte
	Msb     byte
}

func (self *NoteOff) MidiEventData()                 {}
func (self *NoteOn) MidiEventData()                  {}
func (self *PolyphonicKeyPressure) MidiEventData()   {}
func (self *ControllerChange) MidiEventData()        {}
func (self *ProgramChange) ProgramChange()           {}
func (self *ChannelKeyPressure) ChannelKeyPressure() {}
func (self *PitchBend) PitchBend()                   {}

// Sysex events

// TODO

// -- Meta events --

type SetTempo struct {
	Tempo int
}

func (self *SetTempo) MetaEventData() {}
func (self *SetTempo) String() string {
	return fmt.Sprintf("Tempo: %d", self.Tempo)
}

type TrackName struct {
	Name string
}

func (self *TrackName) MetaEventData() {}
func (self *TrackName) String() string {
	return fmt.Sprintf("Track name: %s", self.Name)
}

type TimeSignature struct {
	Numerator   byte
	Denominator byte
	Clocks      byte
	Notes       byte
}

func (self *TimeSignature) MetaEventData() {}
func (self *TimeSignature) String() string {
	return fmt.Sprintf("TimeSignature: %d/%d %d %d",
		self.Numerator, self.Denominator, self.Clocks, self.Notes)
}

type Marker struct {
	Name string
}

func (self *Marker) MetaEventData() {}
func (self *Marker) String() string {
	return fmt.Sprintf("Marker: %s", self.Name)
}

type EndOfTrack struct{}

func (self *EndOfTrack) MetaEventData() {}
func (self *EndOfTrack) String() string {
	return "End of Track"
}

// Channnel mode messages
// not type safe (Experimental)

// deprecated
func (mes *ChannelModeMessage) String() string {
	return fmt.Sprintf("Channel Mode Message Track: %d Type: %X (%s) Value: %d",
		mes.Channel, mes.StatusByte, mes.Message, mes.Value)
}

func (me *MidiEvent) String() string {
	return fmt.Sprintf("<Channel %d> %X (%s) [%v]", me.Channel+1, me.StatusByte, me.Message, me.Values)
}
