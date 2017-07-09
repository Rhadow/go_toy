package dom

import (
	"errors"
	"fmt"
	"strings"
)

// ElementNodeType - constant for element node
const ElementNodeType = "DOM/ELEMENT_NODE"

// TextNodeType - constant for element node
const TextNodeType = "DOM/TEXT_NODE"

// Node - Node interface
type Node interface {
	GetChildren() []Node
	GetNodeType() string
}

// TextNode - a dom text node
type TextNode struct {
	text string
}

// GetChildren - get children for text node
func (t TextNode) GetChildren() []Node {
	return []Node{}
}

// GetNodeType - get node type for text node
func (t TextNode) GetNodeType() string {
	return TextNodeType
}
func (t TextNode) String() string {
	return fmt.Sprintf("Text Node: %s", t.text)
}

// ElementNode - an html element consists of tagName and attributes
type ElementNode struct {
	TagName    string
	Attributes AttrMap
	Children   []Node
}

// GetChildren - get children for element node
func (e ElementNode) GetChildren() []Node {
	return e.Children
}

// GetNodeType - get node type for element node
func (e ElementNode) GetNodeType() string {
	return ElementNodeType
}

// GetID - Get the id of this element
func (e ElementNode) GetID() (string, error) {
	if id, ok := e.Attributes["id"]; ok {
		return id, nil
	}
	return "", errors.New("No id in attribute")
}

// GetClasses - Get the classes of this element
func (e ElementNode) GetClasses() ([]string, error) {
	if classes, ok := e.Attributes["class"]; ok {
		return strings.Split(classes, " "), nil
	}
	return []string{}, errors.New("No class in attribute")
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
		Children:   children,
		TagName:    tagName,
		Attributes: attrs,
	}
}
