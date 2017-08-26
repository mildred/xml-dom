package xpath

import (
	"github.com/antchfx/xpath"
	"github.com/mildred/xml-dom"
	"github.com/mildred/xml-dom/node-navigator"
)

type Expr struct {
	*xpath.Expr
}

func convert(e *xpath.Expr) *Expr {
	if e == nil {
		return nil
	} else {
		return &Expr{e}
	}
}

func CompileNS(expr string, namespaces map[string]string) (*Expr, error) {
	// FIXME: handle namespaces
	e, err := xpath.Compile(expr)
	return convert(e), err
}

func MustCompileNS(expr string, namespaces map[string]string) *Expr {
	// FIXME: handle namespaces
	return convert(xpath.MustCompile(expr))
}

func Compile(expr string) (*Expr, error) {
	e, err := xpath.Compile(expr)
	return convert(e), err
}

func MustCompile(expr string) *Expr {
	return convert(xpath.MustCompile(expr))
}

func (e *Expr) Evaluate(n *xmldom.Node) interface{} {
	res := e.Expr.Evaluate(node_navigator.NewNodeNavigator(n))
	switch res.(type) {
	case *xpath.NodeIterator:
		i := res.(*xpath.NodeIterator)
		return &Iterator{i}
	default:
		return res
	}
}

// Iterator for nodes, initialized before the first element. Call MoveNext() to
// get started. When MoveNext() returns false, it means we are past the end and
// there is no current element.
type Iterator struct {
	*xpath.NodeIterator
}

func (i *Iterator) Current() *xmldom.Node {
	res := i.NodeIterator.Current()
	if res == nil {
		return nil
	} else if nn, ok := res.(*node_navigator.NodeNavigator); ok {
		return nn.Node
	} else {
		panic("Could not convert NodeNavigator")
	}
}

func (i *Iterator) MoveNext() bool {
	return i.NodeIterator.MoveNext()
}
