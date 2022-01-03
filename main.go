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
)

const (
	asmForPlus = `
	movq %d, %%rax\n
	addq %d, %%rax
	`

	asmForMinus = `
	movq %d, %%rax\n
	subq %d, %%rax
	`

	asmForTimes = `
	movq %d, %%rax\n
	imulq %d, %%rax
	`

	// the same applies to mod
	asmForDiv = `
	movq %d, %%rax
	movq %d, %%rcx
	cltd 
	idivq %%rcx
	`
)

type Exp interface {
}

type Plus struct {
	Exp
	x int
	y int
}

type Minus struct {
	Exp
	x int
	y int
}

type Times struct {
	Exp
	x int
	y int
}

type Div struct {
	Exp
	x int
	y int
}

type Mod struct {
	Exp
	x int
	y int
}

// kind of interpreter
func eval(e Exp) int {
	switch e.(type) {
	case Plus:
		plusExpr := e.(Plus)
		return plusExpr.x + plusExpr.y
	case Minus:
		minusExpr := e.(Minus)
		return minusExpr.x - minusExpr.y
	case Times:
		timesExpr := e.(Times)
		return timesExpr.x * timesExpr.y
	case Div:
		divExpr := e.(Div)
		return divExpr.x / divExpr.y
	case Mod:
		modExpr := e.(Mod)
		return modExpr.x % modExpr.y
	default:
		panic(fmt.Errorf("encountered unknown node %v", e))
	}
}

// kind of compiler

func trans(e Exp) string {
	switch e.(type) {
	case Plus:
		plusExpr := e.(Plus)
		return fmt.Sprintf(asmForPlus, plusExpr.x, plusExpr.y)
	case Minus:
		minusExpr := e.(Minus)
		return fmt.Sprintf(asmForMinus, minusExpr.x, minusExpr.y)
	case Times:
		timesExpr := e.(Times)
		return fmt.Sprintf(asmForTimes, timesExpr.x, timesExpr.y)
	case Div:
		divExpr := e.(Div)
		return fmt.Sprintf(asmForDiv, divExpr.x, divExpr.y)
	case Mod:
		modExpr := e.(Mod)
		return fmt.Sprintf(asmForDiv, modExpr.x, modExpr.y)
	default:
		panic(fmt.Errorf("encountered unknown node %v", e))
	}
}

func prepareRuntimeLib() {
	driverCode := `
		#include <stdio.h>

		int miniscala_main();

		int main(int argc, const char* argv[]){
			printf("%d\n", miniscala_main());
			return 0;
		}
	`

	handle, err := os.OpenFile("sources/driver.c", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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
		.globl _miniscala_main
		_miniscala_main:
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
	log.Printf("wrote %d bytes", bts)
	handle.Close()

	compileCodeCmd := exec.Command("gcc", "runtime.o", "driver.s", "-o", "driver")
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
	run(trans(Times{x: 20, y: 30}))
}
