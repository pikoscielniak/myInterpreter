package main
import (
	"fmt"
	"unicode"
	"strconv"
	"bufio"
	"os"
	"strings"
)
type TokenType string

const (
	INTEGER TokenType = "INTEGER"
	PLUS TokenType = "PLUS"
	MINUS TokenType = "MINUS"
	MUL TokenType = "MUL"
	DIV TokenType = "DIV"
	EOF TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value interface{}
}

func (t Token)String() string {
	return fmt.Sprintf("Token({%v}, {%v})", t.Type, t.Value)
}

type Lexer struct {
	Text        string
	Pos         int
	CurrentChar *rune
}

func NewLexer(text string) *Lexer {
	first := rune(text[0])
	i := &Lexer{
		Text: text,
		Pos: 0,
		CurrentChar: &first,
	}
	return i
}

func (i *Lexer) Error(msg string) {
	panic("LEXER " + msg)
}

func (i *Lexer) advance() {
	i.Pos += 1
	if i.Pos > len(i.Text) - 1 {
		i.CurrentChar = nil
	} else {
		r := rune(i.Text[i.Pos])
		i.CurrentChar = &r
	}
}

func (i *Lexer) skipWhitespace() {
	for i.CurrentChar != nil && unicode.IsSpace(*i.CurrentChar) {
		i.advance()
	}
}

func (i *Lexer) integer() int {
	result := ""
	for i.CurrentChar != nil && unicode.IsDigit(*i.CurrentChar) {
		result += string(*i.CurrentChar)
		i.advance()
	}
	res, err := strconv.Atoi(result)
	if err != nil {
		i.Error(fmt.Sprintf("Integer %v", err))
	}
	return res
}

func (i *Lexer) getNextToken() *Token {

	for i.CurrentChar != nil {
		if unicode.IsSpace(*i.CurrentChar) {
			i.skipWhitespace()
			continue
		}
		if unicode.IsDigit(*i.CurrentChar) {
			return &Token{
				Type: INTEGER,
				Value:i.integer(),
			}
		}
		if *i.CurrentChar == '+' {
			i.advance()
			return &Token{Type:PLUS, Value:"+"}
		}
		if *i.CurrentChar == '-' {
			i.advance()
			return &Token{Type:MINUS, Value:"-"}
		}
		if *i.CurrentChar == '*' {
			i.advance()
			return &Token{Type:MUL, Value:"*"}
		}
		if *i.CurrentChar == '/' {
			i.advance()
			return &Token{Type:DIV, Value:"/"}
		}
		i.Error("getNextToken Not match")
	}
	return &Token{
		Type:EOF, Value: nil}
}

type Interpreter struct {
	lexer        *Lexer
	CurrentToken *Token
}

func NewInterpreter(l *Lexer) *Interpreter {
	i := &Interpreter{
		lexer: l,
	}
	i.CurrentToken = i.lexer.getNextToken()
	return i
}

func (i *Interpreter) Error(msg string) {
	panic("INTERPRETER " + msg)
}

func (i *Interpreter) Eat(tType TokenType) {
	if i.CurrentToken.Type == tType {
		i.CurrentToken = i.lexer.getNextToken()
	} else {
		i.Error("Eat: Error parsing input")
	}
}

func (i *Interpreter) factor() int {
	t := i.CurrentToken
	i.Eat(INTEGER)
	return t.Value.(int)
}


func (i *Interpreter) term() int {
	result := i.factor()

	for i.CurrentToken.Type == MUL ||
	i.CurrentToken.Type == DIV {
		token := i.CurrentToken
		if token.Type == MUL {
			i.Eat(MUL)
			result = result * i.factor()
		} else if token.Type == DIV {
			i.Eat(DIV)
			result = result / i.factor()
		}
	}
	return result
}

func (i *Interpreter) Expr() int {
	result := i.term()

	for i.CurrentToken.Type == PLUS ||
	i.CurrentToken.Type == MINUS {
		token := i.CurrentToken
		if token.Type == PLUS {
			i.Eat(PLUS)
			result = result + i.term()
		} else if token.Type == MINUS {
			i.Eat(MINUS)
			result = result - i.term()
		}
	}
	return result
}

//
//func (i *Interpreter) term() int {
//	t := i.CurrentToken
//	i.Eat(INTEGER)
//	return t.Value.(int)
//}

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	//	text  := "2+3"
	fmt.Println(text)

	lex := NewLexer(text)
	inter := NewInterpreter(lex)
	res := inter.Expr()
	fmt.Printf("Result: %v\n", res)

}