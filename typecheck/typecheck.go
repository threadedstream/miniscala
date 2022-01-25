package typecheck

import (
	"fmt"
	"github.com/ThreadedStream/miniscala/backing"
	"github.com/ThreadedStream/miniscala/syntax"
	"os"
)

func typecheckError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	os.Exit(1)
}

// TODO(threadedstream): breathe some life into (yet unlearned) typechecking machinery

func inferValueType(expr syntax.Expr, typeEnv backing.TypeEnv) backing.ValueType {
	switch expr.(type) {
	default:
		typecheckError("unknown node in inferValueType()")
	case *syntax.BasicLit:
		basicLit := expr.(*syntax.BasicLit)
		return backing.LitKindToValueType(basicLit.Kind)
	case *syntax.Name:
		name := expr.(*syntax.Name)
		typeInfo, ok := backing.LookupType(name.Value, typeEnv, false)
		if !ok {
			// probably, this is just a type name
			valueType := backing.MiniscalaTypeToValueType(name.Value)
			if valueType == backing.Undefined {
				typecheckError("name %s is neither a type name nor var, nor val", name.Value)
			}
			return valueType
		}
		return typeInfo.ValueType
	case *syntax.Field:
		field := expr.(*syntax.Field)
		return inferValueType(field.Type, typeEnv)
	case *syntax.Operation:
		operation := expr.(*syntax.Operation)
		lhsType := inferValueType(operation.Lhs, typeEnv)
		rhsType := inferValueType(operation.Rhs, typeEnv)
		resultingType, compatible := typesCompatible(lhsType, rhsType)
		if !compatible {
			typecheckError("types %s and %s are not compatible",
				backing.ValueTypeToStr(lhsType),
				backing.ValueTypeToStr(rhsType))
		}
		return resultingType
	}
}

// typesCompatible checks against compatibility of passed types and returns
// a resulting type upon success
func typesCompatible(t1, t2 backing.ValueType) (backing.ValueType, bool) {
	switch {
	default:
		return backing.Undefined, false
	case t1 == backing.Float && t2 == backing.Int:
		return backing.Float, true
	case t1 == backing.Float && t2 == backing.Float:
		return backing.Float, true
	case t1 == backing.Int && t2 == backing.Int:
		return backing.Int, true
	case t1 == backing.Bool && t2 == backing.Bool:
		return backing.Bool, true
	case t1 == backing.String && t2 == backing.Bool:
		return backing.String, true
	}
}

func typecheckProgram(program *syntax.Program) {
	for _, stmt := range program.StmtList {
		typecheckStmt(stmt, nil)
	}
}

func typecheckStmt(stmt syntax.Stmt, typeEnv backing.TypeEnv) {
	switch stmt.(type) {
	case *syntax.BlockStmt:
		typecheckBlockStmt(stmt)
	}
}

func typecheckBlockStmt(stmt syntax.Stmt) {
	blockStmt := stmt.(*syntax.BlockStmt)
	typeEnv := make(backing.TypeEnv)
	for _, decStmt := range blockStmt.Stmts {
		typecheckStmt(decStmt, typeEnv)
	}
}

func typecheckVarDeclStmt(stmt syntax.Stmt, typeEnv backing.TypeEnv) {
	varDeclStmt := stmt.(*syntax.VarDeclStmt)
	inferredType := inferValueType(varDeclStmt.Rhs, typeEnv)
	backing.StoreType(varDeclStmt.Name.Value, inferredType, false, typeEnv)
}

func typecheckValDeclStmt(stmt syntax.Stmt, typeEnv backing.TypeEnv) {
	valDeclStmt := stmt.(*syntax.ValDeclStmt)
	inferredType := inferValueType(valDeclStmt.Rhs, typeEnv)
	backing.StoreType(valDeclStmt.Name.Value, inferredType, true, typeEnv)
}

func typecheckDefDeclStmt(stmt syntax.Stmt) {
	defDeclStmt := stmt.(*syntax.DefDeclStmt)
	for _, param := range defDeclStmt.ParamList {
		typecheckField(param, nil)
	}
}

func typecheckReturnStmt(stmt syntax.Stmt) {

}

func typecheckField(field *syntax.Field, typeEnv backing.TypeEnv) {
	valueType := inferValueType(field, typeEnv)
	backing.StoreType(field.Name.Value, valueType, false, typeEnv)
}
