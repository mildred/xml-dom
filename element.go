package xmldom

func (n *Node) TagName() string {
	if n.nodeType != ElementNode {
		panic("only on an element")
	}
	return n.nodeName
}

func (n *Node) GetAttribute(name string) string {
	if n.nodeType != ElementNode {
		panic("only on an element")
	}
	attr := n.GetAttributeNode(name)
	if attr == nil {
		return ""
	} else {
		return attr.NodeValue()
	}
}

func (n *Node) SetAttribute(name string, value string) Error {
	var err Error
	if n.nodeType != ElementNode {
		panic("only on an element")
	}
	attr := n.GetAttributeNode(name)
	if attr == nil {
		attr, err = n.OwnerDocument().CreateAttribute(name)
		if err != nil {
			return err
		}
		attr.nodeValue = value
		_, err = n.SetAttributeNode(attr)
		return err
	} else {
		attr.SetNodeValue(value)
		return nil
	}
}

func (n *Node) RemoveAttribute(name string) Error {
	if n.nodeType != ElementNode {
		panic("only on an element")
	}
	return n.attributes.RemoveNamedItem(name)
}

func (n *Node) GetAttributeNode(name string) *Node {
	if n.nodeType != ElementNode {
		panic("only on an element")
	}
	return n.attributes.GetNamedItem(name)
}

func (n *Node) SetAttributeNode(newAttr *Node) (*Node, Error) {
	if n.nodeType != ElementNode {
		panic("only on an element")
	}
	a := n.attributes.GetNamedItem(newAttr.NodeName())
	if a != nil && newAttr != a {
		n.attributes.RemoveNamedItem(newAttr.NodeName())
	} else if a == newAttr {
		return a, nil
	}
	return n.attributes.SetNamedItem(newAttr)
}

func (n *Node) RemoveAttributeNode(oldAttr *Node) (*Node, Error) {
	if n.nodeType != ElementNode {
		panic("only on an element")
	}
	a := n.attributes.GetNamedItem(oldAttr.NodeName())
	if a != oldAttr {
		return nil, err(NotFoundError)
	}
	n.attributes.RemoveNamedItem(oldAttr.NodeName())
	return oldAttr, nil
}
