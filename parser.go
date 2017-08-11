package xmldom

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"log"
)

func xmlName(n xml.Name) string {
	if n.Space == "" {
		return n.Local
	} else {
		return n.Space + ":" + n.Local
	}
}

func (parent *Node) NewChildFromToken(tok xml.Token, data string) *Node {
	var n *Node
	var err error
	doc := parent.ownerDocument
	switch tok := tok.(type) {
	case xml.StartElement:
		se := parseStartElement(data)
		n, err = doc.CreateElement(xmlName(tok.Name))
		if err != nil {
			panic(err)
		}
		for i, attr := range tok.Attr {
			sea := se.Attributes[i]
			a, err := doc.CreateAttribute(xmlName(attr.Name))
			if err != nil {
				panic(err)
			}
			a.SetNodeValue(attr.Value)
			a.Raw = []string{sea.Before, sea.Name, sea.Between, sea.Value, sea.After}
			a.ValueDirty = false
			log.Printf("attr: %#v", a)
			n.SetAttributeNode(a)
		}
		n.Raw = []string{se.Before, se.TagName, se.After}
		n.ValueDirty = false
	case xml.EndElement:
		panic("end element")
	case xml.CharData:
		n = doc.CreateTextNode(string(tok))
		n.Raw = append(n.Raw, data)
		n.ValueDirty = false
	case xml.Comment:
		n = doc.CreateComment(string(tok))
		n.Raw = append(n.Raw, data)
		n.ValueDirty = false
	case xml.ProcInst: // Processing Instruction
		n, err = doc.CreateProcessingInstruction(tok.Target, string(tok.Inst))
		n.Raw = append(n.Raw, data)
		n.ValueDirty = false
	case xml.Directive:
		n = doc.CreateDocumentType(string(tok))
		n.Raw = append(n.Raw, data)
		n.ValueDirty = false
	}
	if err != nil {
		panic(err)
	}
	_, err = parent.AppendChild(n)
	if err != nil {
		panic(err)
	}
	return n
}

type xmlReader struct {
	r      io.ByteReader
	line   int64
	col    int64
	offset int64
	acc    []byte
	last   []byte
}

func (r *xmlReader) Read(p []byte) (n int, err error) {
	panic("Read")
}

func (r *xmlReader) ReadByte() (byte, error) {
	for _, b := range r.last {
		r.offset++
		if b == '\n' {
			r.line++
			r.col = 0
		} else {
			r.col++
		}
		r.acc = append(r.acc, b)
	}
	r.last = nil
	b, err := r.r.ReadByte()
	if err != nil {
		return b, err
	}
	r.last = append(r.last, b)
	return b, nil
}

func ParseXML(rr io.Reader) (*Node, error) {
	var r *xmlReader
	if rb, ok := rr.(io.ByteReader); ok {
		r = &xmlReader{rb, 0, 0, 0, nil, nil}
	} else {
		r = &xmlReader{bufio.NewReader(rr), 0, 0, 0, nil, nil}
	}

	decoder := xml.NewDecoder(r)
	decoder.Strict = false
	doc := NewDocument()
	node := doc
	for {
		//log.Println()
		//log.Printf("Before token o=%d l=%d c=%d", r.offset, r.line+1, r.col+1)
		//log.Printf("xml offset=%d", decoder.InputOffset())
		r.acc = nil
		tok, err := decoder.RawToken()
		switch {
		case err == io.EOF:
			//log.Printf("End of File")
			goto quit
		case err != nil:
			return nil, err
		}

		inclLast := false
		switch tok.(type) {
		case xml.ProcInst:
			inclLast = true
		case xml.EndElement:
			inclLast = true
		case xml.StartElement:
			inclLast = true
		case xml.CharData:
			if bytes.HasPrefix(r.acc, []byte("<![CDATA[")) && bytes.HasSuffix(r.acc, []byte("]]")) && string(r.last) == ">" {
				inclLast = true
			}
		}
		if inclLast {
			r.acc = append(r.acc, r.last...)
			r.last = nil
		}

		//log.Printf("After token o=%d l=%d c=%d", r.offset, r.line+1, r.col+1)
		//log.Printf("acc=%#v", string(r.acc))
		//log.Printf("xml offset=%d", decoder.InputOffset())
		//log.Printf("tok=%#v", tok)

		switch tok := tok.(type) {
		default:
			log.Printf("node: %#v", string(r.acc))
			node.NewChildFromToken(tok, string(r.acc))
		case xml.StartElement:
			log.Printf("start: %#v", string(r.acc))
			node = node.NewChildFromToken(tok, string(r.acc))
		case xml.EndElement:
			log.Printf("end: %#v", string(r.acc))
			for node.parentNode != nil && xmlName(tok.Name) != node.nodeName {
				p := node.parentNode
				for _, n := range node.childNodes {
					_, err := p.AppendChild(n)
					if err != nil {
						panic(err)
					}
				}
				node.childNodes = NodeList{}
				node = p
			}
			node.Raw = append(node.Raw, string(r.acc))
			log.Printf("end2: %#v, %#v", string(node.nodeName), node.parentNode)
			node = node.parentNode
			log.Printf("end2: %#v", string(node.nodeName))
		}
	}
quit:
	log.Printf("doc: %s", doc.XML())
	return doc, nil
}
