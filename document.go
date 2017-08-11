package xmldom

func (d *Node) DocumentElement() *Node {
	if d.nodeType != DocumentNode {
		panic("only on a document")
	}
	for _, n := range d.childNodes {
		if n.NodeType() == ElementNode {
			return n
		}
	}
	return nil
}

func NewDocument() *Node {
	n := &Node{
		nodeType:      DocumentNode,
		pos:           -1,
		nodeName:      "#document",
		nodeValue:     "",
		ValueDirty:    false,
		parentNode:    nil,
		childNodes:    NodeList{},
		ownerDocument: nil,
		attributes:    nil,
	}
	n.ownerDocument = n
	return n
}

func (d *Node) CreateElement(tagName string) (*Node, Error) {
	return &Node{
		nodeType:      ElementNode,
		pos:           -1,
		nodeName:      tagName,
		nodeValue:     "",
		ValueDirty:    false,
		parentNode:    nil,
		childNodes:    NodeList{},
		ownerDocument: d.ownerDocument,
		attributes:    NewEmptyNamedNodeMap(d.ownerDocument),
	}, nil
}

func (d *Node) CreateDocumentFragment() *Node {
	return &Node{
		nodeType:      DocumentFragmentNode,
		pos:           -1,
		nodeName:      "#document-fragment",
		nodeValue:     "",
		ValueDirty:    false,
		parentNode:    nil,
		childNodes:    NodeList{},
		ownerDocument: d.ownerDocument,
		attributes:    nil,
	}
}

func (d *Node) CreateTextNode(data string) *Node {
	return &Node{
		nodeType:      TextNode,
		pos:           -1,
		nodeName:      "#text",
		nodeValue:     data,
		ValueDirty:    false,
		parentNode:    nil,
		childNodes:    nil,
		ownerDocument: d.ownerDocument,
		attributes:    nil,
	}
}

func (d *Node) CreateComment(data string) *Node {
	return &Node{
		nodeType:      CommentNode,
		pos:           -1,
		nodeName:      "#comment",
		nodeValue:     data,
		ValueDirty:    false,
		parentNode:    nil,
		childNodes:    nil,
		ownerDocument: d.ownerDocument,
		attributes:    nil,
	}
}

func (d *Node) CreateDocumentType(data string) *Node {
	return &Node{
		nodeType:      DocumentTypeNode,
		pos:           -1,
		nodeName:      "#document-type",
		nodeValue:     data,
		ValueDirty:    false,
		parentNode:    nil,
		childNodes:    nil,
		ownerDocument: d.ownerDocument,
		attributes:    nil,
	}
}

func (d *Node) CreateCDATASection(data string) (*Node, Error) {
	return &Node{
		nodeType:      CDATASectionNode,
		pos:           -1,
		nodeName:      "#cdata-section",
		nodeValue:     data,
		ValueDirty:    false,
		parentNode:    nil,
		childNodes:    nil,
		ownerDocument: d.ownerDocument,
		attributes:    nil,
	}, nil
}

func (d *Node) CreateProcessingInstruction(target, data string) (*Node, Error) {
	return &Node{
		nodeType:      ProcessingInstructionNode,
		pos:           -1,
		nodeName:      target,
		nodeValue:     data,
		ValueDirty:    false,
		parentNode:    nil,
		childNodes:    nil,
		ownerDocument: d.ownerDocument,
		attributes:    nil,
	}, nil
}

func (d *Node) CreateAttribute(name string) (*Node, Error) {
	return &Node{
		nodeType:      AttributeNode,
		pos:           -1,
		nodeName:      name,
		nodeValue:     "",
		ValueDirty:    false,
		parentNode:    nil,
		childNodes:    nil,
		ownerDocument: d.ownerDocument,
		attributes:    nil,
	}, nil
}

func (d *Node) CreateEntityReference(name string) (*Node, Error) {
	return &Node{
		nodeType:      EntityReferenceNode,
		pos:           -1,
		nodeName:      name,
		nodeValue:     "",
		ValueDirty:    false,
		parentNode:    nil,
		childNodes:    nil,
		ownerDocument: d.ownerDocument,
		attributes:    nil,
	}, nil
}
