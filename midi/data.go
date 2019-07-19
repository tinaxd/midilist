package midi

import "fmt"

type MidiElement interface {
	MidiElement()
}

type MidiMeta struct {
	Format   int
	Tracks   int
	Division int
}

func (m MidiMeta) MidiElement() {}

func (m MidiMeta) String() string {
	return fmt.Sprintf("Format: %d, Tracks: %d, Division: %d",
		m.Format, m.Tracks, m.Division)
}
