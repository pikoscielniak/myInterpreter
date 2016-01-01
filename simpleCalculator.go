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
	LPAREN = "("
	RPAREN = ")"
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

func (i *Lexer) getNextToken() Token {

	for i.CurrentChar != nil {
		if unicode.IsSpace(*i.CurrentChar) {
			i.skipWhitespace()
			continue
		}
		if unicode.IsDigit(*i.CurrentChar) {
			return Token{
				Type: INTEGER,
				Value:i.integer(),
			}
		}
		if *i.CurrentChar == '+' {
			i.advance()
			return Token{Type:PLUS, Value:"+"}
		}
		if *i.CurrentChar == '-' {
			i.advance()
			return Token{Type:MINUS, Value:"-"}
		}
		if *i.CurrentChar == '*' {
			i.advance()
			return Token{Type:MUL, Value:"*"}
		}
		if *i.CurrentChar == '/' {
			i.advance()
			return Token{Type:DIV, Value:"/"}
		}
		if *i.CurrentChar == '(' {
			i.advance()
			return Token{Type:LPAREN, Value:"("}
		}
		if *i.CurrentChar == ')' {
			i.advance()
			return Token{Type:RPAREN, Value:")"}
		}
		i.Error("getNextToken Not match")
	}
	return Token{
		Type:EOF, Value: nil}
}

type AST struct {

}

type BinOp struct {
	Left  interface{}
	Token Token
	Op    Token
	Right interface{}
}

func NewBinOp(left interface{}, op Token, right interface{}) *BinOp {
	return &BinOp{
		Left:left,
		Token: op,
		Op: op,
		Right:right,
	}
}

type Num struct {
	Token Token
	Value int
}

func NewNum(t Token) *Num {
	return &Num{
		Token: t,
		Value:t.Value.(int),
	}
}

type Parser struct {
	Lexer        *Lexer
	CurrentToken Token
}

func NewParser(l *Lexer) *Parser {
	return &Parser{
		Lexer:l,
		CurrentToken:l.getNextToken(),
	}
}

func (i *Parser) Error(msg string) {
	panic("PARSER " + msg)
}

func (i *Parser) eat(tType TokenType) {
	if (i.CurrentToken.Type == tType) {
		i.CurrentToken = i.Lexer.getNextToken()
	} else {
		i.Error("Eat wrong type")
	}
}

func (i *Parser) factor() interface{} {
	t := i.CurrentToken
	if t.Type == INTEGER {
		i.eat(INTEGER)
		return NewNum(t)
	} else if t.Type == LPAREN {
		i.eat(LPAREN)
		node := i.expr()
		i.eat(RPAREN)
		return node
	}
	i.Error("factor not supported type")
	return nil
}

func (i *Parser) term() interface{} {
	node := i.factor()

	for i.CurrentToken.Type == MUL ||
	i.CurrentToken.Type == DIV {
		token := i.CurrentToken
		if token.Type == MUL {
			i.eat(MUL)
		} else if token.Type == DIV {
			i.eat(DIV)
		}

		node = &BinOp{Left:node, Op:token, Right: i.factor()}
	}
	return node
}

func (i *Parser) expr() interface{} {
	node := i.term()

	for i.CurrentToken.Type == PLUS ||
	i.CurrentToken.Type == MINUS {
		token := i.CurrentToken
		if token.Type == PLUS {
			i.eat(PLUS)
		} else if token.Type == MINUS {
			i.eat(MINUS)
		}
		node = &BinOp{Left:node, Op:token, Right:i.term()}
	}
	return node
}

func (i *Parser) parse() interface{} {
	return i.expr()
}

type Interpreter struct {
	Parser *Parser
}

func NewInterpreter(p *Parser) *Interpreter {
	return &Interpreter{
		Parser:p,
	}
}

func (n *Interpreter) visitBinOp(node *BinOp) int {
	if node.Op.Type == PLUS {
		return n.visit(node.Left) + n.visit(node.Right)
	} else if node.Op.Type == MINUS {
		return n.visit(node.Left) - n.visit(node.Right)
	} else if node.Op.Type == MUL {
		return n.visit(node.Left) * n.visit(node.Right)
	} else if node.Op.Type == DIV {
		return n.visit(node.Left) / n.visit(node.Right)
	}
	n.Error("visitBinOp unknown operator type")
	return 0;
}


func (i *Interpreter) Error(msg string) {
	panic("INTERPRETER " + msg)
}

func (n *Interpreter) visitNum(node *Num) int {
	return node.Value
}

func (n *Interpreter) visit(node interface{}) int {
	switch t := node.(type){
	case *BinOp:
		return n.visitBinOp(t)
	case *Num:
		return n.visitNum(t)
	}
	return 0
}

func (n *Interpreter) interpret() int {
	tree := n.Parser.parse()
	return n.visit(tree)
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	//			text  := "2+3"
	//	fmt.Println(text)

	lex := NewLexer(text)
	parser := NewParser(lex)
	inter := NewInterpreter(parser)
	res := inter.interpret()
	fmt.Printf("Result: %v\n", res)

//	mulTok := Token{Type:MUL, Value:"*"}
//	plusTok := Token{Type:PLUS, Value:"+"}
//	num2 := NewNum(Token{INTEGER, 2})
//	num7 := NewNum(Token{INTEGER, 7})
//	mulNode := NewBinOp(num2, mulTok, num7)
//	num3 := NewNum(Token{INTEGER, 3})
//	addNode := NewBinOp(mulNode, plusTok, num3)
//
//	intr := NewInterpreter(nil)
//	res := intr.visit(addNode)
//	fmt.Printf("Result: %v\n", res)
}