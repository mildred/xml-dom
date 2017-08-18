package node_navigator

import (
	"github.com/antchfx/xpath"
	"github.com/mildred/xml-dom"
	"io/ioutil"
	"log"
)

var _ xpath.NodeNavigator = &NodeNavigator{}

var l *log.Logger = log.New(ioutil.Discard, "", log.LstdFlags)

type NodeNavigator struct {
	Node *xmldom.Node
	Attr int
}

func NewNodeNavigator(node *xmldom.Node) *NodeNavigator {
	return &NodeNavigator{node, 0}
}

// NodeType returns the XPathNodeType of the current node.
func (nn *NodeNavigator) NodeType() xpath.NodeType {
	l.Printf("NodeType(%v) = %#v", nn.node(), nn.node().NodeType())
	switch nn.node().NodeType() {
	case xmldom.DocumentNode:
		return xpath.RootNode
	case xmldom.ElementNode:
		return xpath.ElementNode
	case xmldom.AttributeNode:
		return xpath.AttributeNode
	case xmldom.TextNode:
		return xpath.TextNode
	case xmldom.CDATASectionNode:
		return xpath.TextNode
	case xmldom.CommentNode:
		return xpath.CommentNode
	case xmldom.DocumentFragmentNode:
		fallthrough
	case xmldom.ProcessingInstructionNode:
		fallthrough
	case xmldom.DocumentTypeNode:
		fallthrough
	case xmldom.EntityReferenceNode:
		fallthrough
	case xmldom.EntityNode:
		fallthrough
	case xmldom.NotationNode:
		fallthrough
	default:
		return xpath.CommentNode
	}
}

// LocalName gets the Name of the current node.
func (nn *NodeNavigator) node() *xmldom.Node {
	if nn.Attr == 0 {
		return nn.Node
	} else {
		return nn.Node.Attributes().Item(nn.Attr - 1)
	}
}

// LocalName gets the Name of the current node.
func (nn *NodeNavigator) LocalName() string {
	l.Printf("LocalName(%v) = %#v", nn.node(), nn.node().LocalNodeName())
	return nn.node().LocalNodeName()
}

// Prefix returns namespace prefix associated with the current node.
func (nn *NodeNavigator) Prefix() string {
	l.Printf("Prefix(%v)", nn.node())
	return nn.node().NodeNamePrefix()
}

// Value gets the value of current node.
func (nn *NodeNavigator) Value() string {
	l.Printf("Value(%v)", nn.node())
	return nn.node().NodeValue()
}

// Copy does a deep copy of the NodeNavigator and all its components.
func (nn *NodeNavigator) Copy() xpath.NodeNavigator {
	var nn2 NodeNavigator = *nn
	return &nn2
}

// MoveToRoot moves the NodeNavigator to the root node of the current node.
func (nn *NodeNavigator) MoveToRoot() {
	l.Printf("MoveToRoot(%v)", nn.node())
	defer l.Printf("MoveToRoot(%v) *END*", nn.node())
	for nn.MoveToParent() {
	}
}

// MoveToParent moves the NodeNavigator to the parent node of the current node.
func (nn *NodeNavigator) MoveToParent() bool {
	l.Printf("MoveToParent(%v)", nn.node())
	if nn.Attr != 0 {
		nn.Attr = 0
		return true
	}
	if parent := nn.Node.ParentNode(); parent != nil {
		nn.Node = parent
		return true
	} else {
		return false
	}
}

// MoveToNextAttribute moves the NodeNavigator to the next attribute on current node.
func (nn *NodeNavigator) MoveToNextAttribute() bool {
	l.Printf("MoveToNextAttribute(%v)", nn.node())
	if nn.Attr >= nn.Node.Attributes().Length() {
		return false
	} else {
		nn.Attr++
		return true
	}
}

// MoveToChild moves the NodeNavigator to the first child node of the current node.
func (nn *NodeNavigator) MoveToChild() bool {
	l.Printf("MoveToChild(%v)", nn.node())
	if nn.Attr != 0 {
		return false
	}
	if child := nn.Node.FirstChild(); child != nil {
		nn.Node = child
		return true
	} else {
		return false
	}
}

// MoveToFirst moves the NodeNavigator to the first sibling node of the current node.
func (nn *NodeNavigator) MoveToFirst() bool {
	l.Printf("MoveToFirst(%v)", nn.node())
	defer l.Printf("MoveToFirst(%v) *END*", nn.node())
	if nn.Attr != 0 {
		return false
	}
	if !nn.MoveToPrevious() {
		return false
	}
	for nn.MoveToPrevious() {
	}
	return true
}

// MoveToNext moves the NodeNavigator to the next sibling node of the current node.
func (nn *NodeNavigator) MoveToNext() bool {
	if nn.Attr != 0 {
		l.Printf("MoveToNext(%v) ERROR", nn.node())
		return false
	}
	if sibling := nn.Node.NextSibling(); sibling != nil {
		l.Printf("MoveToNext(%v) = %v", nn.node(), sibling)
		nn.Node = sibling
		return true
	} else {
		l.Printf("MoveToNext(%v) = false", nn.node())
		return false
	}
}

// MoveToPrevious moves the NodeNavigator to the previous sibling node of the current node.
func (nn *NodeNavigator) MoveToPrevious() bool {
	l.Printf("MoveToPrevious(%v)", nn.node())
	if nn.Attr != 0 {
		return false
	}
	if sibling := nn.Node.PreviousSibling(); sibling != nil {
		nn.Node = sibling
		return true
	} else {
		return false
	}
}

// MoveTo moves the NodeNavigator to the same position as the specified NodeNavigator.
func (nn *NodeNavigator) MoveTo(nn2 xpath.NodeNavigator) bool {
	l.Printf("MoveTo(%v, %v)", nn.node(), nn2)
	if n := nn2.(*NodeNavigator); n != nil {
		nn.Node = n.Node
		nn.Attr = n.Attr
		return true
	} else {
		return false
	}
}

func (nn *NodeNavigator) String() string {
	return nn.node().String()
}
