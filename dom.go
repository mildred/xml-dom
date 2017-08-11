package xmldom

type ErrorCode uint

const (
	_ ErrorCode = iota
	IndexSizeError
	DOMStringSizeError
	HierarchyRequestError
	WrongDocumentError
	InvalidCharacterError
	NoDataAllowedError
	NoModificationAllowedError
	NotFoundError
	NotSupportedError
	InuseAttributeError
)

type Error interface {
	error
	Code() ErrorCode
}

type NodeType uint

const (
	ElementNode NodeType = iota + 1
	AttributeNode
	TextNode
	CDATASectionNode
	EntityReferenceNode
	EntityNode
	ProcessingInstructionNode
	CommentNode
	DocumentNode
	DocumentTypeNode
	DocumentFragmentNode
	NotationNode
)

type Implementation interface {
	HasFeature(feature, version string) bool
}

type NodeInterface interface {
	Pos() uint // extension
	NodeName() string
	NodeValue() string
	SetNodeName(string)
	SetNodeValue(string)
	NodeType() NodeType
	ParentNode() *Node
	ChildNodes() NodeList
	FirstChild() *Node
	LastChild() *Node
	PreviousSibling() *Node
	NextSibling() *Node
	Attributes() NamedNodeMap
	OwnerDocument() Document

	InsertBefore(newChild, refChild *Node) (*Node, Error)
	ReplaceChild(newChild, oldChild *Node) (*Node, Error)
	RemoveChild(oldChild *Node) (*Node, Error)
	AppendChild(newChild *Node) (*Node, Error)
	HasChildNodes() bool
	CloneNode(deep bool) *Node
}

type NodeList []*Node

type NamedNodeMap interface {
	GetNamedItem(name string) *Node
	SetNamedItem(node *Node) (*Node, Error)
	RemoveNamedItem(name string) Error
	Item(index int) *Node
	Length() int
	Clone(deep bool) NamedNodeMap
}

type DocumentFragment interface {
	NodeInterface
}

type Document interface {
	NodeInterface
	//Doctype() DocumentType
	//Implementation() Implementation
	DocumentElement() Element
	//GetElementsByTagName(tagName string) NodeList

	CreateElement(tagName string) (Element, Error)
	CreateDocumentFragment() DocumentFragment
	CreateTextNode(data string) Text
	CreateComment(data string) Comment
	CreateCDATASection(data string) (CDATASection, Error)
	CreateProcessingInstruction(target, data string) (ProcessingInstruction, Error)
	CreateAttribute(name string) (Attr, Error)
	CreateEntityReference(name string) (Entity, Error)
}

type CharacterData interface {
	NodeInterface
	Data() string
	SetData(s string)
	Length() uint

	//SubstringData(offset, count uint) (string, Error)
	//AppendData(arg string) Error
	//InsertData(offset uint, arg string) Error
	//DeteleData(offset, count uint) Error
	//ReplaceData(offset, count long, arg string) Error
}

type Attr interface {
	NodeInterface
	Name() string
	Specified() bool
	Value() string
	SetValue(s string)
}

type Element interface {
	NodeInterface
	TagName() string
	GetAttribute(name string) string
	SetAttribute(name string, value string) Error
	RemoveAttribute(name string) Error
	GetAttributeNode(name string) Attr
	SetAttributeNode(newAttr Attr) (Attr, Error)
	RemoveAttributeNode(oldAttr Attr) (Attr, Error)
	//GetElementsByTagName(name string) NodeList
	Normalize()
}

type Text interface {
	CharacterData
	//SplitText(offset uint) (Text, Error)
}

type Comment interface {
	CharacterData
}

type CDATASection interface {
	Text
}

type DocumentType interface {
	NodeInterface
	Name() string
	Entities() NamedNodeMap
	Notations() NamedNodeMap
}

type Notation interface {
	NodeInterface
	PublicId() string
	SystemId() string
}

type Entity interface {
	NodeInterface
	PublicId() string
	SystemId() string
	NotationName() string
}

type EntityReference interface {
	NodeInterface
}

type ProcessingInstruction interface {
	NodeInterface
	Target() string
	Data() string
}
