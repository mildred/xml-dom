package xmldom

import (
	"strings"
)

type startElementAttribute struct {
	Before  string
	Name    string
	Between string
	Value   string
	After   string
}

type startElement struct {
	Before     string
	TagName    string
	Attributes []startElementAttribute
	After      string
}

type whileUntil uint

const (
	parseWhileCode whileUntil = iota
	parseUntilCode
)

func parseWhileUntil(code string, separators string, while_until whileUntil) (string, string) {
	var prefix string
	for i, r := range code {
		has_rune := strings.ContainsRune(separators, r)
		if !has_rune && while_until == parseWhileCode {
			return prefix, code[i:]
		}
		if has_rune && while_until == parseUntilCode {
			return prefix, code[i:]
		}
		prefix += string(r)
	}
	return prefix, ""
}

func parseWhile(code, separators string) (string, string) {
	return parseWhileUntil(code, separators, parseWhileCode)
}

func parseUntil(code, separators string) (string, string) {
	return parseWhileUntil(code, separators, parseUntilCode)
}

const xmlWhitespace string = " \x0d\x09\x0a"

func parseStartElement(code string) startElement {
	var elem startElement
	var sp string
	elem.Before, code = parseWhile(code, "<"+xmlWhitespace)
	elem.TagName, code = parseUntil(code, "/>"+xmlWhitespace)
	sp, code = parseWhile(code, xmlWhitespace)
	for !strings.HasPrefix(code, "/") && !strings.HasPrefix(code, ">") {
		var attr startElementAttribute
		attr.Before = sp
		attr.Name, code = parseUntil(code, "="+xmlWhitespace)
		sp, code = parseWhile(code, xmlWhitespace)
		if strings.HasPrefix(code, "=") {
			attr.Between, code = parseWhile(code, "="+xmlWhitespace)
			attr.Between = sp + attr.Between
			if (strings.HasPrefix(code, "'") || strings.HasPrefix(code, "\"")) && len(code) >= 2 {
				var quote string = code[0:1]
				code = code[1:]
				attr.Between += quote
				attr.Value, code = parseUntil(code, quote)
				if len(code) >= 1 {
					attr.After = code[0:1]
					code = code[1:]
				} else {
					attr.After = quote
				}
			} else {
				attr.Value, code = parseUntil(code, xmlWhitespace)
				attr.After = ""
			}
		}
		elem.Attributes = append(elem.Attributes, attr)
	}
	elem.After = sp + code
	return elem
}
