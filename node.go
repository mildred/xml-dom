package xmldom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type Node struct {
	nodeType      NodeType
	pos           int
	nodeName      string
	nodeValue     string
	ValueDirty    bool
	parentNode    *Node
	childNodes    NodeList
	ownerDocument *Node
	attributes    NamedNodeMap // not nil only for elements

	// For:
	// - start and end elements: "<", tagName, ">", "</tagName>"
	// - self closing elements: "<", tagName, "/>", ""
	// - attributes: " ", name, "=\"", value, "\""
	Raw []string
}

func (n *Node) String() string {
	return n.XML()
}

func (n *Node) XML() string {
	switch n.nodeType {
	case DocumentFragmentNode:
		fallthrough
	case DocumentNode:
		var res string
		for _, cn := range n.childNodes {
			res += cn.XML()
		}
		return res
	case ElementNode:
		var res string
		if len(n.Raw) >= 1 {
			res += n.Raw[0]
		} else {
			res += "<"
		}
		res += n.nodeName
		for i := 0; i < n.Attributes().Length(); i++ {
			a := n.Attributes().Item(i)
			res += a.XML()
		}

		var last string
		if len(n.Raw) >= 4 &&
			(n.Raw[3] == "") == (len(n.childNodes) == 0) &&
			(len(n.childNodes) == 0 || n.Raw[1] == n.nodeName) {
			last = n.Raw[3]
			res += n.Raw[2]
		} else if len(n.childNodes) == 0 {
			res += "/>"
			last = ""
		} else {
			res += ">"
			last = fmt.Sprintf("</%s>", n.nodeName)
		}

		for _, cn := range n.childNodes {
			res += cn.XML()
		}

		res += last
		return res
	case AttributeNode:
		var res string
		if len(n.Raw) >= 1 {
			res += n.Raw[0]
		} else {
			res += " "
		}
		res += n.nodeName
		if len(n.Raw) >= 5 && !n.ValueDirty {
			res += n.Raw[2]
			res += n.Raw[3]
			res += n.Raw[4]
		} else {
			var b bytes.Buffer
			res += "=\""
			err := xml.EscapeText(&b, []byte(n.nodeValue))
			if err != nil {
				panic(err)
			}
			res += string(b.Bytes())
			res += "\""
		}
		return res
	case TextNode:
		if len(n.Raw) > 0 && !n.ValueDirty {
			return strings.Join(n.Raw, "")
		}
		var b bytes.Buffer
		err := xml.EscapeText(&b, []byte(n.nodeValue))
		if err != nil {
			panic(err)
		}
		return string(b.Bytes())
	case CDATASectionNode:
		if len(n.Raw) > 0 && !n.ValueDirty {
			return strings.Join(n.Raw, "")
		}
		return "<![CDATA[" + n.nodeValue + "]]>"
	case ProcessingInstructionNode:
		if len(n.Raw) > 0 && !n.ValueDirty {
			return strings.Join(n.Raw, "")
		}
		return "<?" + n.nodeName + " " + n.nodeValue + "?>"
	case CommentNode:
		if len(n.Raw) > 0 && !n.ValueDirty {
			return strings.Join(n.Raw, "")
		}
		return "<!--" + n.nodeValue + "-->"
	case DocumentTypeNode:
		if len(n.Raw) > 0 && !n.ValueDirty {
			return strings.Join(n.Raw, "")
		}
		return "<!" + n.nodeValue + ">"
	case EntityReferenceNode:
		fallthrough
	case EntityNode:
		fallthrough
	case NotationNode:
		fallthrough
	default:
		if len(n.Raw) > 0 && !n.ValueDirty {
			return strings.Join(n.Raw, "")
		}
		panic("Unknown node type " + n.nodeName)
	}
}

func (n *Node) CloneNode(deep bool) *Node {
	var children NodeList
	if deep {
		for _, c := range n.childNodes {
			children = append(children, c.CloneNode(deep))
		}
	}
	var attributes NamedNodeMap
	if n.attributes != nil {
		attributes = n.attributes.Clone(deep)
	}
	return &Node{
		nodeType:      n.nodeType,
		pos:           n.pos,
		nodeName:      n.nodeName,
		nodeValue:     n.nodeValue,
		parentNode:    nil,
		childNodes:    children,
		ownerDocument: n.ownerDocument,
		attributes:    attributes,
	}
}

func (n *Node) NodeType() NodeType {
	return n.nodeType
}

func (n *Node) NodeName() string {
	return n.nodeName
}

func (n *Node) LocalNodeName() string {
	slice := strings.SplitN(n.nodeName, ":", 2)
	if len(slice) >= 2 {
		return slice[1]
	} else {
		return n.nodeName
	}
}

func (n *Node) NodeNamePrefix() string {
	slice := strings.SplitN(n.nodeName, ":", 2)
	if len(slice) >= 2 {
		return slice[0]
	} else {
		return ""
	}
}

func (n *Node) NodeValue() string {
	return n.nodeValue
}

func (n *Node) SetNodeName(s string) {
	n.nodeName = s
}

func (n *Node) SetNodeValue(s string) {
	n.nodeValue = s
	n.ValueDirty = true
}

func (n *Node) ParentNode() *Node {
	return n.parentNode
}

func (n *Node) ChildNodes() NodeList {
	return n.childNodes
}

func (n *Node) FirstChild() *Node {
	if len(n.childNodes) > 0 {
		return n.childNodes[0]
	} else {
		return nil
	}
}

func (n *Node) LastChild() *Node {
	if len(n.childNodes) > 0 {
		return n.childNodes[len(n.childNodes)-1]
	} else {
		return nil
	}
}

func (n *Node) PreviousSibling() *Node {
	if n.ParentNode() == nil || n.pos == 0 {
		return nil
	} else {
		return n.ParentNode().ChildNodes()[n.pos-1]
	}
}

func (n *Node) NextSibling() *Node {
	if n.ParentNode() == nil {
		return nil
	}
	siblings := n.ParentNode().ChildNodes()
	if n.pos+1 >= len(siblings) {
		return nil
	} else {
		return siblings[n.pos+1]
	}
}

func (n *Node) Attributes() NamedNodeMap {
	return n.attributes
}

func (n *Node) OwnerDocument() *Node {
	return n.ownerDocument
}

func (n *Node) InsertBefore(newChild, refChild *Node) (*Node, Error) {
	if refChild == nil {
		return n.AppendChild(newChild)
	}

	i := refChild.pos
	if refChild.ParentNode() != n || n.childNodes[i] != refChild {
		return nil, err(NotFoundError)
	}

	var newChildren NodeList
	if i > 0 {
		newChildren = n.childNodes[0:i]
	}
	var length int

	if newChild.NodeType() == DocumentFragmentNode {
		for j, c := range newChild.ChildNodes() {
			err := c.attach(n, i+j)
			if err != nil {
				return nil, err
			}
		}
		newChildren = append(newChildren, newChild.ChildNodes()...)
		length = len(newChild.ChildNodes())
	} else {
		err := newChild.attach(n, i)
		if err != nil {
			return nil, err
		}
		length = 1
		newChildren = append(newChildren, newChild)
	}

	for _, child := range n.childNodes[i:len(n.childNodes)] {
		child.pos += length
	}
	newChildren = append(newChildren, n.childNodes[i:len(n.childNodes)]...)
	n.childNodes = newChildren
	return newChild, nil
}

func (n *Node) ReplaceChild(newChild, oldChild *Node) (*Node, Error) {
	if newChild.NodeType() == DocumentFragmentNode {
		return nil, err(HierarchyRequestError)
	}
	i := oldChild.pos
	if oldChild.ParentNode() != n || n.childNodes[i] != oldChild {
		return nil, err(NotFoundError)
	}
	err := newChild.attach(n, i)
	if err != nil {
		return nil, err
	}
	n.childNodes[i] = newChild
	oldChild.parentNode = nil
	oldChild.pos = 0
	return oldChild, nil
}

func (n *Node) RemoveChild(oldChild *Node) (*Node, Error) {
	i := oldChild.pos
	if oldChild.parentNode != n || n.childNodes[i] != oldChild {
		return nil, err(NotFoundError)
	}
	var newChildren NodeList
	if i > 0 {
		newChildren = n.childNodes[0:i]
	}
	if i+1 < len(n.childNodes) {
		for _, cn := range n.childNodes[i+1 : len(n.childNodes)] {
			cn.pos -= 1
			newChildren = append(newChildren, cn)
		}
	}
	n.childNodes = newChildren
	oldChild.parentNode = nil
	oldChild.pos = 0
	return oldChild, nil
}

func (n *Node) AppendChild(newChild *Node) (*Node, Error) {
	if newChild.NodeType() == DocumentFragmentNode {
		for _, child := range newChild.ChildNodes() {
			_, err := n.AppendChild(child)
			if err != nil {
				return nil, err
			}
		}
		return newChild, nil
	}

	err := newChild.attach(n, len(n.childNodes))
	if err != nil {
		return nil, err
	}
	n.childNodes = append(n.childNodes, newChild)
	return newChild, nil
}

func (n *Node) attach(newParent *Node, pos int) Error {
	if newParent.OwnerDocument() != n.ownerDocument {
		return err(WrongDocumentError)
	} else if newParent.IsAncestor(n) {
		return err(HierarchyRequestError)
	}
	if n.parentNode != nil {
		_, err := n.parentNode.RemoveChild(n)
		if err != nil {
			panic(err)
		}
	}
	n.parentNode = newParent
	n.pos = pos
	return nil
}

func (n *Node) IsAncestor(ancestor *Node) bool {
	var nn *Node = n
	for nn.ParentNode() != nil {
		if nn.ParentNode() == ancestor {
			return true
		}
		nn = nn.ParentNode()
	}
	return false
}

func (n *Node) HasChildNodes() bool {
	return len(n.childNodes) > 0
}
