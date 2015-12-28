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
	EOF TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value interface{}
}

func (t Token)String() string {
	return fmt.Sprintf("Token({%v}, {%v})", t.Type, t.Value)
}

type Interpreter struct {
	Text         string
	Pos          int
	CurrentToken *Token
	CurrentChar  *rune
}

func NewInterpreter(text string) *Interpreter {
	first := rune(text[0])
	i := &Interpreter{
		Text: text,
		Pos: 0,
		CurrentToken: nil,
		CurrentChar: &first,
	}
	return i
}

func (i *Interpreter) Error(msg string) {
	panic(msg)
}

func (i *Interpreter) advance() {
	i.Pos += 1
	if i.Pos > len(i.Text) - 1 {
		i.CurrentChar = nil
	} else {
		r := rune(i.Text[i.Pos])
		i.CurrentChar = &r
	}
}

func (i *Interpreter) skipWhitespace() {
	for i.CurrentChar != nil && unicode.IsSpace(*i.CurrentChar) {
		i.advance()
	}
}

func (i *Interpreter) integer() int {
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

func (i *Interpreter) getNextToken() *Token {

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
		i.Error("getNextToken Not match")
	}
	return &Token{
		Type:EOF, Value: nil}
}

func (i *Interpreter) Eat(tType TokenType) {
	if i.CurrentToken.Type == tType {
		i.CurrentToken = i.getNextToken()
	} else {
		i.Error("Eat: Error parsing input")
	}
}

func (i *Interpreter) Expr() int {
	i.CurrentToken = i.getNextToken()

	left := i.CurrentToken
	i.Eat(INTEGER)
	op := i.CurrentToken
	if op.Type == PLUS {
		i.Eat(PLUS)
	} else {
		i.Eat(MINUS)
	}

	right := i.CurrentToken
	i.Eat(INTEGER)

	lVal := left.Value.(int);
	rVal := right.Value.(int);

	var result int
	if op.Type == PLUS {
		result = lVal + rVal
	} else {
		result = lVal - rVal
	}
	return result
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	//	text  := "2+3"
	fmt.Println(text)

	inter := NewInterpreter(text)
	res := inter.Expr()
	fmt.Printf("Result: %v\n", res)

}