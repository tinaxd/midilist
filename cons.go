package main

import (
	"fmt"

	"github.com/tinaxd/midilist/midi"
)

func main() {
	lexer := midi.NewLexer("./test.mid")
	lexer.LoadData()
	token, _ := lexer.NextToken()
	fmt.Printf("%s\n", token.String())
	parser := midi.NewChunkParser(token)
	fmt.Printf("%s\n", parser.ParseMThd().String())
	for {
		tok, err := lexer.NextToken()
		if err != nil {
			break
		}
		fmt.Printf("%s\n", tok.String())
		p := midi.NewChunkParser(tok)
		p.ParseMTrk()
		break
	}
	//lexer.TestFunc()
}
