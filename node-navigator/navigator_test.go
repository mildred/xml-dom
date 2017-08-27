package node_navigator

import (
	"bytes"
	"github.com/antchfx/xpath"
	"github.com/antchfx/xpath/test"
	"github.com/mildred/xml-dom"
	"testing"
)

func TestXQuery(t *testing.T) {
	create := func(doc string) (xpath.NodeNavigator, error) {
		dom, err := xmldom.ParseXML(bytes.NewBuffer([]byte(doc)))
		if err != nil {
			return nil, err
		}
		return NewNodeNavigator(dom), nil
	}

	xquerytest.TestAll(t, create, xquerytest.EnableAll)
}
