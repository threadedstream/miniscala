package main

import "C"
import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ThreadedStream/miniscala/syntax"
	"github.com/ThreadedStream/miniscala/typecheck"
	"github.com/ThreadedStream/miniscala/vm"
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

// kind of vm

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
	path := "sources/composition.miniscala"

	program, hadErrors := syntax.Parse(path)
	if hadErrors {
		return
	}

	hadErrors = typecheck.Typecheck(program)
	if hadErrors {
		return
	}

	vmHandle := vm.NewVM(program)

	vmHandle.Run()

	//execVm := vm.NewVM(program)
	//
	//execVm.Run()

	//program := &syntax.Program{
	//	StmtList: []syntax.Stmt{
	//		&syntax.DefDeclStmt{
	//			Name: &syntax.Name{
	//				Value: "add",
	//			},
	//			ParamList: []*syntax.Field{
	//				{
	//					Name: &syntax.Name{
	//						Value: "x",
	//					},
	//					Type: &syntax.Name{
	//						Value: "Int",
	//					},
	//				},
	//				{
	//					Name: &syntax.Name{
	//						Value: "y",
	//					},
	//					Type: &syntax.Name{
	//						Value: "Int",
	//					},
	//				},
	//			},
	//			ReturnType: &syntax.Name{
	//				Value: "Int",
	//			},
	//			Body: &syntax.BlockStmt{
	//				Stmts: []syntax.Stmt{
	//					&syntax.VarDeclStmt{
	//						Name: syntax.Name{
	//							Value: "z",
	//						},
	//						Rhs: &syntax.BasicLit{
	//							Value: "0",
	//							Kind:  syntax.FloatLit,
	//						},
	//					},
	//					&syntax.Assignment{
	//						Lhs: &syntax.Name{
	//							Value: "z",
	//						},
	//						Rhs: &syntax.Operation{
	//							Op: syntax.Plus,
	//							Lhs: &syntax.Name{
	//								Value: "x",
	//							},
	//							Rhs: &syntax.Name{
	//								Value: "y",
	//							},
	//						},
	//					},
	//					&syntax.Return{
	//						Value: &syntax.Name{
	//							Value: "z",
	//						},
	//					},
	//				},
	//			},
	//		},
	//		&syntax.Call{
	//			CalleeName: &syntax.Name{
	//				Value: "print",
	//			},
	//			ArgList: []syntax.Expr{
	//				&syntax.Call{
	//					CalleeName: &syntax.Name{
	//						Value: "add",
	//					},
	//					ArgList: []syntax.Expr{
	//						&syntax.BasicLit{
	//							Value: "10",
	//							Kind:  syntax.FloatLit,
	//						},
	//						&syntax.BasicLit{
	//							Value: "30",
	//							Kind:  syntax.FloatLit,
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}

	//interpreter.Execute(program)
	//interpreter.DumpEnvState()

	//factorialCode := []vm.Instruction{
	//	&vm.InstrLoadRef{
	//		RefName: "x",
	//	},
	//	&vm.InstrLoadImm{
	//		Value: backing.Value{
	//			Value:     1.0,
	//			ValueType: backing.Float,
	//		},
	//	},
	//	&vm.InstrEqual{},
	//	&vm.InstrJmpIfFalse{
	//		Offset: 2,
	//	},
	//	&vm.InstrLoadImm{
	//		Value: backing.Value{
	//			Value:     1.0,
	//			ValueType: backing.Float,
	//		},
	//	},
	//	&vm.InstrReturn{},
	//	&vm.InstrLoadRef{
	//		RefName: "x",
	//	},
	//	&vm.InstrLoadRef{
	//		RefName: "x",
	//	},
	//	&vm.InstrLoadImm{
	//		Value: backing.Value{
	//			Value:     1.0,
	//			ValueType: backing.Float,
	//		},
	//	},
	//	&vm.InstrSub{},
	//	&vm.InstrCall{
	//		FuncName: "fac",
	//		ArgNames: []string{"x"},
	//	},
	//	&vm.InstrMul{},
	//	&vm.InstrReturn{},
	//}
	//
	//fibCode := []vm.Instruction{
	//	&vm.InstrLoadRef{
	//		RefName: "n",
	//	},
	//	&vm.InstrLoadImm{
	//		Value: backing.Value{
	//			Value:     1.0,
	//			ValueType: backing.Float,
	//		},
	//	},
	//	&vm.InstrLessThanOrEqual{},
	//	&vm.InstrJmpIfFalse{
	//		Offset: 2,
	//	},
	//	&vm.InstrLoadImm{
	//		Value: backing.Value{
	//			Value:     1.0,
	//			ValueType: backing.Float,
	//		},
	//	},
	//	&vm.InstrReturn{},
	//	&vm.InstrLoadRef{
	//		RefName: "n",
	//	},
	//	&vm.InstrLoadImm{
	//		Value: backing.Value{
	//			Value:     1.0,
	//			ValueType: backing.Float,
	//		},
	//	},
	//	&vm.InstrSub{},
	//	&vm.InstrCall{
	//		FuncName: "fib",
	//		ArgNames: []string{"n"},
	//	},
	//	&vm.InstrLoadRef{
	//		RefName: "n",
	//	},
	//	&vm.InstrLoadImm{
	//		Value: backing.Value{
	//			Value:     2.0,
	//			ValueType: backing.Float,
	//		},
	//	},
	//	&vm.InstrSub{},
	//	&vm.InstrCall{
	//		FuncName: "fib",
	//		ArgNames: []string{"n"},
	//	},
	//	&vm.InstrAdd{},
	//	&vm.InstrReturn{},
	//}
	//
	//mainCode := []vm.Instruction{
	//	&vm.InstrLoadImm{
	//		Value: backing.Value{
	//			Value:     5.0,
	//			ValueType: backing.Float,
	//		},
	//	},
	//	&vm.InstrCall{
	//		FuncName: "fib",
	//		ArgNames: []string{"n"},
	//	},
	//	&vm.InstrCall{
	//		FuncName: "fac",
	//		ArgNames: []string{"x"},
	//	},
	//	&vm.InstrCall{
	//		FuncName: "print",
	//	},
	//	&vm.InstrReturn{},
	//}
	//
	//executionEnv := vm.InitializeVm(mainCode, factorialCode, fibCode)
	//executionEnv.Run()

}
