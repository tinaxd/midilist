package main

import (
	"fmt"

	"github.com/tinaxd/midilist/midi"
)

func main() {
	lexer := midi.NewLexer("./test.mid")
	lexer.LoadData()
	token := lexer.NextToken()
	fmt.Printf("%s\n", token.String())
	token = lexer.NextToken()
	fmt.Printf("%s\n", token.String())
	//lexer.TestFunc()
}
