package main

import (
	"fmt"
	"log"
	"os"

	"io/ioutil"

	"github.com/tinaxd/midilist/midi"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	/*
		p := midi.ParserData{
			[]struct {
				StatusByte byte   `yaml:"statusByte"`
				Message    string `yaml:"message"`
			}{{64, "RPN (LSB)"}},
		}
		e, err := yaml.Marshal(p)
		fmt.Println(string(e))
		return */

	lexer := midi.NewLexer("./test.mid")
	lexer.LoadData()
	token, _ := lexer.NextToken()
	fmt.Printf("%s\n", token.String())
	parser := midi.NewChunkParser(token)
	fmt.Printf("%s\n", parser.ParseMThd().String())

	var parserData midi.ParserData
	f, err := os.Open("./midiparse.yaml")
	if err != nil {
		log.Panicln("cannot read midiparse.yaml")
	}
	b, err := ioutil.ReadAll(f)
	yaml.Unmarshal(b, &parserData)

	i := 1
	for {
		tok, err := lexer.NextToken()
		if err != nil {
			break
		}
		//fmt.Printf("%s\n", tok.String())
		log.Printf("[Track %d]", i)
		p := midi.NewChunkParser(tok)
		p.ParserData = parserData
		p.ParseMTrk()
		i++
	}
	//lexer.TestFunc()
}
