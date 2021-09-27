package parser

import (
	"errors"
	"go/ast"
	"go/parser"
	"reflect"
	"strconv"
)

// ParseExpr is a convenience function for obtaining the AST of an expression x.
func ParseExpr(expr string) (ast.Expr, error) {
	return parser.ParseExpr(expr)
}

// ParseAndExecuteExpr parse and execute expr
func ParseAndExecuteExpr(expr string, attrs map[string]interface{}) (interface{}, error) {
	var (
		exprAst ast.Expr
		iExpr   *ast.Ident
		lExpr   *ast.BasicLit
		bExpr   *ast.BinaryExpr
		ok      bool
		err     error
	)

	exprAst, err = parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}

	iExpr, ok = exprAst.(*ast.Ident)
	if ok {
		return attrs[iExpr.Name], nil
	}

	lExpr, ok = exprAst.(*ast.BasicLit)
	if ok {
		return parseBasicLit(lExpr)
	}

	bExpr, ok = exprAst.(*ast.BinaryExpr)
	if ok {
		return parseBinaryExpr(bExpr, attrs)
	}

	return nil, errors.New("expr " + expr + " not support")
}

func parseBinaryExpr(expr *ast.BinaryExpr, attrs map[string]interface{}) (interface{}, error) {
	var (
		lRes interface{}
		rRes interface{}
		opt  = expr.Op.String()
		err  error
	)

	if !expr.Op.IsOperator() {
		return nil, errors.New(opt + "is invalid operator")
	}

	switch ex := expr.X.(type) {
	case *ast.Ident:
		lRes = attrs[ex.Name]
	case *ast.BasicLit:
		lRes, err = parseBasicLit(ex)
		if err != nil {
			return nil, err
		}
	case *ast.BinaryExpr:
		lRes, err = parseBinaryExpr(ex, attrs)
		if err != nil {
			return nil, err
		}
	}

	switch ex := expr.Y.(type) {
	case *ast.Ident:
		rRes = attrs[ex.Name]
	case *ast.BasicLit:
		rRes, err = parseBasicLit(ex)
		if err != nil {
			return nil, err
		}
	case *ast.BinaryExpr:
		rRes, err = parseBinaryExpr(ex, attrs)
		if err != nil {
			return nil, err
		}
	}

	if !(lRes != nil && rRes != nil && reflect.TypeOf(lRes) == reflect.TypeOf(rRes)) {
		return nil, errors.New("expr type error")
	}

	switch opt {
	case "==":
		return lRes == rRes, nil
	case "!=":
		return lRes != rRes, nil
	default:
		return nil, errors.New("operator" + opt + " is not be support")
	}
	return nil, nil
}

func parseBasicLit(expr *ast.BasicLit) (interface{}, error) {
	var err error

	if expr.Kind.String() == "INT" {
		var res int64
		res, err = strconv.ParseInt(expr.Value, 10, 0)
		if err != nil {
			return nil, err
		}
		return int(res), nil
	}

	if expr.Kind.String() == "FLOAT" {
		var res float64
		res, err = strconv.ParseFloat(expr.Value, 0)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	if expr.Kind.String() == "STRING" {
		return expr.Value[1 : len(expr.Value)-1], nil
	}

	return expr.Value, nil
}
