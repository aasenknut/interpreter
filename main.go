package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const prompt = "\n>> "

var fname = flag.String(
	"read from file",
	"",
	"If filename is not given, you'll get a prompt",
)

func main() {
	fmt.Println("starting...")

	foo := "./sample-small.txt"
	fname = &foo

	var fContent []byte

	var err error
	if *fname == "" {
		err = openCLI(os.Stdin, os.Stdout)
		if err != nil {
			os.Exit(65)
		}
	} else {
		fContent, err = openFile(*fname)
		if err != nil {
			log.Fatal("read file: %v", err)
		}
	}
	fmt.Println(string(fContent))
	fmt.Println("scanning...")
	lex := Lexer{
		Source:  string(fContent),
		Tokens:  []Token{},
		line:    0,
		start:   0,
		current: 0,
	}
	if err = lex.Scan(); err != nil {
		log.Fatal("lexer scan: %v", err)
	}
	for _, t := range lex.Tokens {
		fmt.Println(t.Lexeme)
	}

	//printAst(exampleExpr())
	expr := parse(lex.Tokens)
	printAst(expr)
	interp := Interpreter{}
	interp.interpret(expr)
	fmt.Println("\n\n==[DONE]==")
}

func openCLI(reader io.Reader, writer io.Writer) error {
	in := bufio.NewScanner(reader)
	for {
		fmt.Fprintf(writer, prompt)
		fmt.Printf("You wrote: %s", in.Text())
		data := in.Scan()
		l := Lexer{}
		l.Source = in.Text()
		l.Scan()
		if !data {
			return fmt.Errorf("erronoeus data")
		}
		fmt.Printf("Tokens: %v", &l.Tokens)
		//line := in.Text()
	}
}

func openFile(fname string) ([]byte, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	return content, nil
}

func printAst(e Expr) {
	fmt.Println("\n[INFO] printing ast...")
	ap := AstPrinter{}
	str, err := ap.Print(e)
	if err != nil {
		fmt.Printf("\n\n[ERROR] %v\n", err)
	}
	fmt.Println(str)
}

func exampleExpr() Expr {
	return &BinaryExpr{
		Left: &UnaryExpr{
			Operator: Token{Type: Minus, Lexeme: "-"},
			Right:    &LiteralExpr{Value: 123},
		},
		Operator: Token{Type: Star, Lexeme: "*"},
		Right:    &GroupingExpr{Expr: &LiteralExpr{45.65}},
	}
}

func parse(t []Token) Expr {
	p := NewParser(t)
	e, err := p.Parse()
	if err != nil {
		fmt.Println("\n[ERROR] parse: %v", err)
	}
	return e
}
