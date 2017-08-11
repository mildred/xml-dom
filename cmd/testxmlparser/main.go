package main

import (
	"bytes"
	"github.com/mildred/xml-dom"
	"log"
)

const xml = `
<?xml version="1.0"?>
<hello >
	<world attr="value"  >
	&lt;
	<![CDATA[ cdata ]]>
	<self-closing tag=val />
</hello >
`

func main() {
	log.Printf("Input document:\n%s\n-----", xml)
	dom, err := xmldom.ParseXML(bytes.NewReader([]byte(xml)))
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Output document:\n%s\n-----", dom.XML())
	}
}
