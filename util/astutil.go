package util

import (
	"go/ast"
	"strconv"
)

// importPath returns the unquoted import path of s,
// or "" if the path is not properly quoted.
func ImportPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err == nil {
		return t
	}
	return ""
}
