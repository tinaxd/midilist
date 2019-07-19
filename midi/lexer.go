package midi

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const BUFSIZE = 1024

type Lexer struct {
	path    string
	data    []byte
	pointer int
}

func NewLexer(path string) *Lexer {
	return &Lexer{
		path,
		make([]byte, BUFSIZE),
		0,
	}
}

func (lexer *Lexer) LoadData() {
	data, err := ioutil.ReadFile(lexer.path)
	if err != nil {
		log.Fatalf("Error loading file")
	}
	lexer.data = data
}

func (lexer *Lexer) LoadData2() {
	file, err := os.Open(lexer.path)
	if err != nil {
		log.Fatal("Error opening file")
	}
	defer file.Close()

	buf := make([]byte, BUFSIZE)
	lexer.data = make([]byte, BUFSIZE)
	log.Print("Loading")
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			log.Fatal("Error loading data")
			break
		}
		log.Print(".")
		lexer.data = append(lexer.data, buf...)
	}
	log.Println("Loaded a midi file")
}

func (lexer *Lexer) NextNByte(length uint) []byte {
	leng := int(length)
	if len(lexer.data) <= lexer.pointer+leng {
		return nil
	}
	ret := lexer.data[lexer.pointer : lexer.pointer+leng]
	lexer.pointer += leng
	return ret
}

type MidiError struct {
	Message string
}

func (le *MidiError) Error() string {
	return fmt.Sprintf("MidiError: %s", le.Message)
}

func (lexer *Lexer) NextToken() (*Chunk, error) {
	chunktype := lexer.NextNByte(4)
	if chunktype == nil {
		return &Chunk{}, &MidiError{"EOF"}
	}
	lengthR := lexer.NextNByte(4)
	if lengthR == nil {
		return &Chunk{}, &MidiError{"EOF"}
	}
	length := binary.BigEndian.Uint32(lengthR)
	data := lexer.NextNByte(uint(length))
	if data == nil {
		return &Chunk{}, &MidiError{"EOF"}
	}
	return &Chunk{
		GetChunkType(chunktype),
		uint(length),
		data,
	}, nil
}

func (lexer *Lexer) TestFunc() {
	for i := 0; i < 100; i++ {
		fmt.Printf("%d", lexer.data[i])
	}
}
