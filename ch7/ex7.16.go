package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"text/scanner"
)

type Var string

type literal float64

type binary struct {
	op   rune
	x, y Expr
}

type unary struct {
	op rune
	x  Expr
}

type call struct {
	fn   string
	args []Expr
}

type Expr interface {
	Eval(env Env) float64
}

func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

func (v Var) Eval(e Env) float64 {
	return e[v]
}

func (u unary) Eval(e Env) float64 {
	switch u.op {
	case '+':
		return u.x.Eval(e)
	case '-':
		return -u.x.Eval(e)
	}
	panic(fmt.Sprintf("unsupported operator type: %c", u.op))
}

func (b binary) Eval(e Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(e) + b.y.Eval(e)
	case '-':
		return b.x.Eval(e) - b.y.Eval(e)
	case '*':
		return b.x.Eval(e) * b.y.Eval(e)
	case '/':
		return b.x.Eval(e) / b.y.Eval(e)
	}
	panic(fmt.Sprintf("unsupported operator type: %c", b.op))
}

func (c call) Eval(e Env) float64 {
	switch c.fn {
	case "sin":
		return math.Sin(c.args[0].Eval(e))
	case "cos":
		return math.Cos(c.args[0].Eval(e))
	case "tan":
		return math.Tan(c.args[0].Eval(e))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(e))
	case "pow":
		return math.Pow(c.args[0].Eval(e), c.args[1].Eval(e))
	case "pi":
		return math.Pi
	}
	panic(fmt.Sprintf("unsupported function type: %s", c.fn))
}

type lexer struct {
	scan  scanner.Scanner
	token rune
}

type lexPanic string

func (lex *lexer) next() {
	// Scan reads the next token or Unicode character from source and returns it.
	// It only recognizes tokens t for which the respective Mode bit (1<<-t) is set.
	lex.token = lex.scan.Scan()
}

func (lex *lexer) text() string {
	// TokenText returns the string corresponding to the most recently scanned token.
	// Valid after calling Scan and in calls of Scanner.Error.
	return lex.scan.TokenText()
}

func (lex *lexer) describe() string {
	switch lex.token {
	case scanner.EOF:
		return "end of file"
	case scanner.Ident:
		return fmt.Sprintf("identifier %s", lex.text())
	case scanner.Int, scanner.Float:
		return fmt.Sprintf("number %s", lex.text())
	}
	return fmt.Sprintf("%q", rune(lex.token))
}

func precedence(op rune) int {
	switch op {
	case '*', '/':
		return 2
	case '+', '-':
		return 1
	}
	return 0
}

func Parse(input string) (_ Expr, err error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
			//
		case lexPanic:
			err = fmt.Errorf("%s", err)
		default:
			panic(x)
		}
	}()
	lex := new(lexer)
	lex.scan.Init(strings.NewReader(input))
	lex.scan.Mode = scanner.ScanInts | scanner.ScanFloats | scanner.ScanIdents
	lex.next()
	e := parseExpr(lex)
	if lex.token != scanner.EOF {
		return nil, fmt.Errorf("unexpected %s", lex.describe())
	}
	return e, nil
}

func parseExpr(lex *lexer) Expr {
	return parseBinary(lex, 1)
}

func parseBinary(lex *lexer, prec1 int) Expr {
	lhs := parseUnary(lex)
	for prec := precedence(lex.token); prec >= prec1; prec-- {
		for precedence(lex.token) == prec {
			op := lex.token
			lex.next()
			rhs := parseBinary(lex, prec+1)
			lhs = binary{op, lhs, rhs}
		}
	}
	return lhs
}

// unary = '+' expr | '-' expr | primary
func parseUnary(lex *lexer) Expr {
	if lex.token == '+' || lex.token == '-' {
		op := lex.token
		lex.next()
		lhs := parseUnary(lex)
		return unary{op, lhs}
	}
	return parsePrimary(lex)
}

// primary = id | id '(' expr ',' ... ',' expr ')' | num | '(' expr ')'
func parsePrimary(lex *lexer) Expr {
	switch lex.token {
	case scanner.Ident:
		id := lex.text()
		lex.next()
		if lex.token != '(' {
			return Var(id)
		}
		lex.next()
		var args []Expr
		if lex.token != ')' {
			for {
				expr := parseExpr(lex)
				args = append(args, expr)
				if lex.token != ',' {
					break
				}
				lex.next()
			}
			if lex.token != ')' {
				msg := fmt.Sprintf("unexpected token, want ), get %c", lex.token)
				panic(lexPanic(msg))
			}
		}
		lex.next()
		return call{id, args}
	case scanner.ScanInts, scanner.ScanFloats:
		res, err := strconv.ParseFloat(lex.text(), 64)
		if err != nil {
			panic(lexPanic(err.Error()))
		}
		lex.next()
		return literal(res)
	case '(':
		lex.next()
		expr := parseExpr(lex)
		if lex.token != ')' {
			panic(lexPanic(fmt.Sprintf("unexpected token, want ), get %c", lex.token)))
		}
		lex.next()
		return expr
	}
	msg := fmt.Sprintf("unexpected %s", lex.describe())
	panic(lexPanic(msg))
}

type Env map[Var]float64

func main() {
	env := Env{}
	http.Handle("/setEnv", http.HandlerFunc(env.setEnv))
	http.Handle("/calc", http.HandlerFunc(env.calc))
	err := http.ListenAndServe("127.0.0.1:8080", nil)

	log.Fatal(err)
}

func (e *Env) calc(w http.ResponseWriter, r *http.Request) {
	expr := r.URL.Query().Get("expr")
	ex, err := Parse(expr)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong expression: %s", expr)
		return
	}
	res := ex.Eval(*e)
	fmt.Fprintf(w, "expr = %f", res)
}

func (e *Env) setEnv(w http.ResponseWriter, r *http.Request) {
	/*	s, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("%s", s)*/
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	_ = decoder.Decode(&params)
	for k, v := range params {
		fmt.Printf("Var %s = %s \n", k, v)
	}
	for k, v := range params {
		expr, err := Parse(v)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Wrong arguments: %s: %s", k, v)
			return
		}
		(*e)[Var(k)] = expr.Eval(*e)
	}
	w.WriteHeader(http.StatusOK)
}
