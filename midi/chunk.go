package midi

import (
	"fmt"
	"log"
)

const (
	MThd = iota
	MTrk
)

type Chunk struct {
	Type   int
	Length uint
	Data   []byte
}

func GetChunkType(raw []byte) int {
	if len(raw) != 4 {
		log.Fatal("Chunk Type invalid length")
	}
	if raw[2] == 0x68 {
		return MThd
	} else {
		return MTrk
	}
}

func (chunk *Chunk) String() string {
	getStr := func(i int) string {
		if i == MThd {
			return "MThd"
		} else {
			return "MTrk"
		}
	}
	return fmt.Sprintf("Type: %s, Length: %d, Data: %v",
		getStr(chunk.Type), chunk.Length, chunk.Data,
	)
}
