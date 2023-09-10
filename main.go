package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var fname = flag.String(
	"read from file",
	"",
	"If filename is not given, you'll get a prompt",
)

func main() {
	fmt.Println("starting...")

	foo := "./resources/sample-code/sample-state.txt"
	fname = &foo

	fContent, err := openFile(*fname)
	if err != nil {
		log.Fatal("read file: %v", err)
	}
	fmt.Println("\nFile content:\n" + string(fContent))
	lex := Lexer{
		Source:  string(fContent),
		Tokens:  []Token{},
		line:    0,
		start:   0,
		current: 0,
	}
	fmt.Println("scanning...")
	if err = lex.Scan(); err != nil {
		log.Fatal("lexer scan: %v", err)
	}
	fmt.Println("scanning DONE")

	fmt.Println("parsing...")
	stmts := parse(lex.Tokens)
	interp := NewInterpreter()
	fmt.Println("interpreting...")
	if err := interp.interpret(stmts); err != nil {
		fmt.Printf("err - interpreter: %v", err)
	}
	fmt.Println("\n\n==[DONE]==")
}

func openFile(fname string) ([]byte, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	return content, err
}

func parse(t []Token) []Stmt {
	p := NewParser(t)
	s, err := p.Parse()
	if err != nil {
		fmt.Println("\n[ERROR] parse: %v", err)
	}
	return s
}
