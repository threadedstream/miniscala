package main

type Parser struct {
	in *TokenScanner
}

func (p *Parser) stmt() Stmt {
	switch p.in.peekToken.(type) {
	case *TokenVal:
		return p.valDecl()
	case *TokenVar:
		return p.varDecl()
	}
	return nil
}

func (p *Parser) varDecl() *VarDeclStmt {
	return nil
}

func (p *Parser) valDecl() *ValDeclStmt {
	return nil
}

func (p *Parser) consume(token Token) {
}

//
//func parse(path string) Exp {
//	stream, err := os.Open(path)
//	if err != nil {
//		panic("no file with such path was found")
//	}
//
//	var parser = &Parser{
//		in: newTokenScanner(stream),
//	}
//	res := parser.expr()
//	if parser.in.hasNext() {
//		panic("expected EOF")
//	}
//	return res
//}

func isComparisonOp(op Operator) bool {
	switch op {
	default:
		return false
	case GreaterThan, GreaterThanOrEqual, LessThan, LessThanOrEqual, Equal, NotEqual:
		return true
	}
}
