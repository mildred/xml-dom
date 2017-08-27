package xmldom

type namedNodeMap struct {
	document *Node
	nodes    NodeList
	index    map[string]int
}

func NewEmptyNamedNodeMap(document *Node) *namedNodeMap {
	return &namedNodeMap{document, nil, map[string]int{}}
}

func (nm *namedNodeMap) Clone(deep bool) NamedNodeMap {
	attributes := &namedNodeMap{
		document: nm.document,
		index:    map[string]int{},
		nodes:    nil,
	}
	for _, n := range nm.nodes {
		attributes.nodes = append(attributes.nodes, n.CloneNode(deep))
	}
	for k, v := range nm.index {
		attributes.index[k] = v
	}
	return attributes
}

func (nm *namedNodeMap) GetNamedItem(name string) *Node {
	if i, ok := nm.index[name]; ok {
		return nm.nodes[i]
	} else {
		return nil
	}
}

func (nm *namedNodeMap) SetNamedItem(item *Node) (*Node, Error) {
	if item.OwnerDocument() != nm.document {
		return nil, err(WrongDocumentError)
	} else if item.ParentNode() != nil {
		return nil, err(InuseAttributeError)
	}
	name := item.NodeName()
	if i, ok := nm.index[name]; ok {
		old := nm.nodes[i]
		nm.nodes[i] = item
		return old, nil
	} else {
		nm.index[name] = len(nm.nodes)
		nm.nodes = append(nm.nodes, item)
		return nil, nil
	}
}

func (nm *namedNodeMap) RemoveNamedItem(name string) Error {
	if i, ok := nm.index[name]; ok {
		nm.nodes[i] = nil
		delete(nm.index, name)
		// reindex
		var newNodes NodeList
		if i > 0 {
			newNodes = nm.nodes[0:i]
		}
		if i+1 < len(nm.nodes) {
			newNodes = append(newNodes, nm.nodes[i+1:len(nm.nodes)]...)
		}
		nm.nodes = newNodes
		for n, j := range nm.index {
			if j > i {
				nm.index[n] = j - 1
			}
		}
		return nil
	} else {
		return err(NotFoundError)
	}
}

func (nm *namedNodeMap) Item(index int) *Node {
	if index < len(nm.nodes) && index >= 0 {
		return nm.nodes[index]
	} else {
		return nil
	}
}

func (nm *namedNodeMap) Length() int {
	return len(nm.nodes)
}
