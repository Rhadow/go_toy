package dom

import "fmt"

// ElementNodeType - constant for element node
const ElementNodeType = "DOM/ELEMENT_NODE"

// TextNodeType - constant for element node
const TextNodeType = "DOM/TEXT_NODE"

// Node - Node interface
type Node interface {
	getChildren() []Node
	getNodeType() string
}

// TextNode - a dom text node
type TextNode struct {
	text string
}

func (t TextNode) getChildren() []Node {
	return []Node{}
}
func (t TextNode) getNodeType() string {
	return TextNodeType
}
func (t TextNode) String() string {
	return fmt.Sprintf("Text Node: %s", t.text)
}

// ElementNode - an html element consists of tagName and attributes
type ElementNode struct {
	tagName    string
	attributes AttrMap
	children   []Node
}

func (e ElementNode) getChildren() []Node {
	return e.children
}
func (e ElementNode) getNodeType() string {
	return ElementNodeType
}

// AttrMap - an attribute map
type AttrMap map[string]string

// CreateTextNode - Create a text dom node
func CreateTextNode(data string) Node {
	return TextNode{
		text: data,
	}
}

// CreateElementNode - Create an Element dom node
func CreateElementNode(tagName string, attrs AttrMap, children []Node) Node {
	return ElementNode{
		children:   children,
		tagName:    tagName,
		attributes: attrs,
	}
}
