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
}

func (i *Interpreter) Error(msg string) {
	panic(msg)
}

func (i *Interpreter) GetNextToken() *Token {
	text := i.Text
	var token *Token
	if i.Pos > len(text) - 1 {
		return &Token{
			Type:EOF, Value: nil}
	}
	currentChar := rune(text[i.Pos])

	if unicode.IsDigit(currentChar) {
		val, err := strconv.Atoi(string(currentChar))
		if err != nil {
			i.Error("GetNextToken. Not Digit")
		}
		token = &Token{
			Type: INTEGER,
			Value:val,
		}
		i.Pos += 1
		return token
	}

	if currentChar == '+' {
		token = &Token{Type:PLUS, Value:currentChar}
		i.Pos += 1
		return token
	}
	i.Error("GetNextToken. Not match. Error parsing input")
	return nil
}

func (i *Interpreter) Eat(tType TokenType) {
	if i.CurrentToken.Type == tType {
		i.CurrentToken = i.GetNextToken()
	} else {
		i.Error("Eat: Error parsing input")
	}
}

func (i *Interpreter) Expr() int {
	i.CurrentToken = i.GetNextToken()

	left := i.CurrentToken
	i.Eat(INTEGER)
	_ = i.CurrentToken
	i.Eat(PLUS)

	right := i.CurrentToken
	i.Eat(INTEGER)

	lVal := left.Value.(int);
	rVal := right.Value.(int);

	return lVal + rVal
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	//	text  := "2+3"
	fmt.Println(text)

	inter := Interpreter{
		Text: text,
	}
	res := inter.Expr()
	fmt.Printf("Result: %v\n", res)

}