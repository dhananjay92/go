// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
This file contains the code to check for basic type comparison.
*/

package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
)

func init() {
	register("cmpbasic",
		"check for basic types comparisons",
		cmdBasic,
		binaryExpr)
}

func cmdBasic(f *File, node ast.Node) {
	e := node.(*ast.BinaryExpr)

	rval, ok := ival(e.Y)
	if !ok {
		return
	}

	ltype := f.pkg.types[e.X].Type
	unsigned := isUnsigned(ltype)
	if unsigned && e.Op == token.LEQ && rval < 0 {
		f.Badf(e.Pos(), "%v (unsigned) <= %d is always false", e.X, rval)
	} else if unsigned && e.Op == token.LSS && rval == 0 {
		f.Badf(e.Pos(), "%v (unsigned) < 0 is always false", e.X)
	}
	// TODO(dhananjay92): Handle 0 >= (unsigned)
}

// isUnsigned reports whether type is basic unsigned.
func isUnsigned(typ types.Type) bool {
	t, ok := typ.(*types.Basic)
	if !ok {
		return false
	}

	return types.Uint <= t.Kind() && t.Kind() <= types.Uintptr
}

// ival returns int value from given basic expr.
func ival(e ast.Expr) (int, bool) {
	switch n := e.(type) {
	case *ast.BasicLit:
		if n.Kind != token.INT {
			return 0, false
		}
		rval, err := strconv.Atoi(n.Value)
		if err != nil {
			return 0, false
		}
		return rval, true
	case *ast.UnaryExpr:
		if n.Op == token.SUB {
			if v, ok := ival(n.X); ok {
				return -1 * v, false
			}
		}
	}

	return 0, false
}
