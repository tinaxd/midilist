package midi

type EventPair struct {
	DeltaTime int32
	Event     Event
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

type MidiEvent struct {
	StatusByte byte
	SataBytes  MidiEventData
}

type SysexEvent struct {
	Type      byte
	IncludeF0 bool
	Data      SysexEventData
}

func (event *MidiEvent) Event()       {}
func (event *SysexEvent) SysexEvent() {}

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

type SysexSetTempo struct {
	Tempo int
}
