package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rhadow/go_toy/browser_engine/dom"
)

// HTMLParser - An Html parser
type HTMLParser struct {
	Position uint64
	Input    string
}

// Read the current character without consuming it.
func (h *HTMLParser) nextChar() string {
	return string(h.Input[h.Position])
}

// Do the next characters start with the given string?
func (h *HTMLParser) startsWith(str string) bool {
	return strings.HasPrefix(h.Input[h.Position:], str)
}

// Return true if all Input is consumed.
func (h *HTMLParser) eof() bool {
	return h.Position >= uint64(len(h.Input))
}

// Return the current character, and advance self.pos to the next character.
func (h *HTMLParser) consumeChar() string {
	currentChar := string(h.Input[h.Position])
	h.Position++
	return currentChar
}

// Consume characters until `test` returns false.
func (h *HTMLParser) consumeWhile(test func(string) bool) string {
	result := ""
	for !h.eof() && test(h.nextChar()) {
		result += h.consumeChar()
	}
	return result
}

// Consume and discard zero or more whitespace characters.
func (h *HTMLParser) consumeWhitespace() {
	h.consumeWhile(func(c string) bool {
		return c == " "
	})
}

// Parse a tag or attribute name.
func (h *HTMLParser) parseTagName() string {
	return h.consumeWhile(func(c string) bool {
		result, err := regexp.MatchString("[a-zA-Z0-9]", c)
		if err != nil {
			panic(err)
		}
		return result
	})
}

// Parse a single node.
func (h *HTMLParser) parseNode() dom.Node {
	if h.nextChar() == "<" {
		return h.parseElement()
	}
	return h.parseText()
}

// Parse a text node.
func (h *HTMLParser) parseText() dom.Node {
	innerText := h.consumeWhile(func(c string) bool {
		return c != "<"
	})
	return dom.CreateTextNode(innerText)
}

// Parse a single element, including its open tag, contents, and closing tag.
func (h *HTMLParser) parseElement() dom.Node {
	// opening tag
	if h.consumeChar() != "<" {
		panic("parse element error: no opening brace in opening tag")
	}
	tagName := h.parseTagName()
	attrs := h.parseAttributes()
	if h.consumeChar() != ">" {
		panic("parse element error: no closing brace in openeing tag")
	}

	// contents
	children := h.parseNodes()

	// closing tag
	if h.consumeChar() != "<" {
		panic("parse element error: no opening brace in closing tag")
	}
	if h.consumeChar() != "/" {
		panic("parse element error: no slash in closing tag")
	}
	if tagName != h.parseTagName() {
		panic("parse element error: tag name is different in closing tag")
	}
	if h.consumeChar() != ">" {
		panic("parse element error: no closing brace in closing tag")
	}

	return dom.CreateElementNode(tagName, attrs, children)
}

// Parse a list of name="value" pairs, separated by whitespace
func (h *HTMLParser) parseAttributes() dom.AttrMap {
	attributes := dom.AttrMap{}
	for {
		h.consumeWhitespace()
		if h.nextChar() == ">" {
			break
		}
		name, value := h.parseAttr()
		attributes[name] = value
	}
	return attributes
}

// Parse a single name="value" pair.
func (h *HTMLParser) parseAttr() (string, string) {
	name := h.parseTagName()
	if h.nextChar() != "=" {
		panic("Parse attribute error: no = sign")
	}
	h.consumeChar()
	value := h.parseAttrValue()
	return name, value
}

// Parse a quoted value
func (h *HTMLParser) parseAttrValue() string {
	openQuote := h.consumeChar()
	fmt.Println(openQuote)
	if openQuote != "\"" && openQuote != "'" {
		panic("Parse attribute error: no quote sign")
	}
	value := h.consumeWhile(func(c string) bool {
		return c != openQuote
	})
	if h.consumeChar() != openQuote {
		panic("Parse attribute error: no closing quote sign")
	}
	return value
}

// Parse a sequence of sibling nodes
func (h *HTMLParser) parseNodes() []dom.Node {
	nodes := []dom.Node{}
	for {
		h.consumeWhitespace()
		if h.eof() || h.startsWith("</") {
			break
		}
		nodes = append(nodes, h.parseNode())
	}
	return nodes
}

// Parse an HTML document and return the root element
func (h *HTMLParser) Parse() dom.Node {
	nodes := h.parseNodes()
	if len(nodes) == 1 {
		return nodes[0]
	}
	return dom.CreateElementNode("html", dom.AttrMap{}, nodes)
}
