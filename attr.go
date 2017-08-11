package xmldom

type DOMAttr struct {
	name      string
	specified bool
	value     string
}

func (a *DOMAttr) NodeType() NodeType {
	return AttributeNode
}

func (a *DOMAttr) Name() string {
	return a.name
}

func (a *DOMAttr) Specified() bool {
	return a.specified
}

func (a *DOMAttr) Value() string {
	return a.value
}

func (a *DOMAttr) SetValue(s string) {
	a.value = s
	a.specified = true
}
