package repl

import (
	"bufio"
	"fmt"
	"github.com/jpiechowka/micron-language-interpreter-go/lexer"
	"github.com/jpiechowka/micron-language-interpreter-go/token"
	"io"
)

const PROMPT = ">>"

func Start(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)

	for {
		fmt.Fprint(output, PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)

		for tok := lex.NextToken(); tok.TokenType != token.EOF; tok = lex.NextToken() {
			fmt.Fprintf(output, "%+v\n", tok)
		}
	}
}
