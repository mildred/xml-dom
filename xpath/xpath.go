package xpath

import (
	"fmt"
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

// Return *Iterator,bool,float64,string
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

// Evaluate that return always an iterator (empty in case the result is not a
// node)
func (e *Expr) EvaluateNode(n *xmldom.Node) *Iterator {
	return &Iterator{e.Expr.Evaluate(node_navigator.NewNodeNavigator(n)).(*xpath.NodeIterator)}
}

func (e *Expr) Exists(n *xmldom.Node) bool {
	return Exists(e.Evaluate(n))
}

func Exists(res interface{}) bool {
	if res == nil {
		return false
	} else if i, ok := res.(*Iterator); ok {
		return i.MoveNext()
	} else {
		return true
	}
}

// Iterator for nodes, initialized before the first element. Call MoveNext() to
// get started. When MoveNext() returns false, it means we are past the end and
// there is no current element.
type Iterator struct {
	*xpath.NodeIterator
}

func (i *Iterator) Node() *xmldom.Node {
	return i.Current()
}

func (i *Iterator) String() string {
	return fmt.Sprintf("%s", i.NodeIterator)
}

func (i *Iterator) Next() bool {
	return i.MoveNext()
}

func (i *Iterator) Nodes() []*xmldom.Node {
	var res []*xmldom.Node
	for i.MoveNext() {
		res = append(res, i.Current())
	}
	return res
}

func (i *Iterator) Current() *xmldom.Node {
	if i.NodeIterator == nil {
		return nil
	}
	res := i.NodeIterator.Current()
	if res == nil {
		return nil
	} else if nn, ok := res.(*node_navigator.NodeNavigator); ok {
		return nn.Current()
	} else {
		panic("Could not convert NodeNavigator")
	}
}

func (i *Iterator) MoveNext() bool {
	if i.NodeIterator == nil {
		return false
	}
	return i.NodeIterator.MoveNext()
}
