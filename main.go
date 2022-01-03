package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
)

const (
	asmForPlus = `
	movq $%d, %%rax
	addq $%d, %%rax
	`

	asmForMinus = `
	movq $%d, %%rax
	subq $%d, %%rax
	`

	asmForTimes = `
	movq $%d, %%rax
	imulq $%d, %%rax
	`

	// the same applies to mod
	asmForDiv = `
	movq $%d, %%rax
	movq $%d, %%rcx
	cltd 
	idivq %%rcx
	`
)

// scanner

func expected(s string) {
	panic(fmt.Errorf("expected %s", s))
}

type Reader struct {
	scanner *scanner.Scanner
}

func (r *Reader) hasNext() bool {
	return r.scanner.Peek() != scanner.EOF
}

func (r *Reader) hasNextP(predicate func(rune) bool) bool {
	return predicate(r.scanner.Peek())
}

func isOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '%'
}

func (r *Reader) getNum() int {
	if r.hasNextP(unicode.IsDigit) {
		n := 0
		for r.hasNextP(unicode.IsDigit) {
			n = 10*n + (int)(r.scanner.Next()-'0')
		}
		return n
	} else {
		expected("number")
	}
	// shouldn't reach this point
	return 0
}

func (r *Reader) parseTerm() Exp {
	return Plus{}
}

func (r *Reader) parseExpression() Exp {
	var left = r.parseTerm()
	if r.hasNextP(isOperator) {
		var op = r.scanner.Next()
		switch op {
		case '+':
			return Plus{x: left, y: r.parseTerm()}
		default:
			panic(fmt.Errorf("unknown operator %c", op))
		}

	} else {
		expected("operator")
		return Plus{}
	}
}

func parse(code string) Exp {
	var (
		reader *Reader = &Reader{}
		res            = reader.parseExpression()
	)

	if reader.hasNext() {
		expected("EOF")
	}
	return res
}

// kind of compiler

func prepareRuntimeLib() {
	driverCode := `
		#include <stdio.h>

		int miniscala_main();

		int main(int argc, const char* argv[]){
			printf("%d\n", miniscala_main());
			return 0;
		}
	`

	handle, err := os.OpenFile("sources/runtime.c", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	bts, err := io.WriteString(handle, driverCode)
	if err != nil {
		panic(err)
	}
	log.Printf("wrote %d bytes", bts)
	handle.Close()

	compileCodeCmd := exec.Command("gcc", "-c", "sources/runtime.c", "-o", "runtime.o")
	if err = compileCodeCmd.Run(); err != nil {
		panic(err)
	}
}

func run(code string) int {
	if _, err := os.Stat("runtime.o"); err != nil {
		prepareRuntimeLib()
	}

	driverCode := `
	.text
	.globl miniscala_main
	miniscala_main:
	%s
	ret
	`

	driverCode = fmt.Sprintf(driverCode, code)

	handle, err := os.OpenFile("sources/driver.s", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		panic(err)
	}
	bts, err := io.WriteString(handle, driverCode)
	if err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes", bts)
	handle.Close()

	compileCodeCmd := exec.Command("gcc", "runtime.o", "sources/driver.s", "-o", "driver")
	if err = compileCodeCmd.Run(); err != nil {
		panic(err)
	}

	// running result
	runCompiledCodeCmd := exec.Command("./driver")
	var buf bytes.Buffer
	runCompiledCodeCmd.Stdout = bufio.NewWriter(&buf)
	if err = runCompiledCodeCmd.Run(); err != nil {
		panic(err)
	}

	res, _ := strconv.Atoi(strings.Trim(buf.String(), "\t\n\r"))

	return res
}

func main() {
	trans(Plus{x: Lit{x: 1}, y: Plus{x: Lit{x: 2}, y: Lit{x: 3}}}, 0)
}
