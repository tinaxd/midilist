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
	len := int(length)
	ret := lexer.data[lexer.pointer : lexer.pointer+len]
	lexer.pointer += len
	return ret
}

func (lexer *Lexer) NextToken() Chunk {
	chunktype := lexer.NextNByte(4)
	lengthR := lexer.NextNByte(4)
	length := binary.BigEndian.Uint32(lengthR)
	data := lexer.NextNByte(uint(length))
	return Chunk{
		GetChunkType(chunktype),
		uint(length),
		data,
	}
}

func (lexer *Lexer) TestFunc() {
	for i := 0; i < 100; i++ {
		fmt.Printf("%d", lexer.data[i])
	}
}
