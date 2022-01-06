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
	// val x = 5; x + 45
	//
	// x0 = 5
	// x1 = x0
	// x2 = 45
	// x3 = x1 + x2
	path := "sources/main.scala"
	node := parse(path)
	fmt.Printf("%x", node)
}
