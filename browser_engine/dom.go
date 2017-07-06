package dom

// Node - Node interface
type Node interface{}

// TextNode - a basic dom node
type TextNode struct {
	children string
}

// ElementNode - an html element consists of tagName and attributes
type ElementNode struct {
	children   []Node
	tagName    string
	attributes AttrMap
}

// AttrMap - an attribute map
type AttrMap map[string]string

// Text - Create a text dom node
func Text(data string) Node {
	return TextNode{
		children: data,
	}
}

// Elem - Create an Element dom node
func Elem(tagName string, attr AttrMap, children []Node) ElementNode {
	return ElementNode{
		children:   children,
		tagName:    tagName,
		attributes: attr,
	}
}
