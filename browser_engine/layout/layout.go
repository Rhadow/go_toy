package layout

import "github.com/rhadow/go_toy/browser_engine/render"

// BoxType - Node Type
type BoxType uint8

// BlockNode - A block node constant
const (
	BlockNode BoxType = iota + 1
	InlineNode
	AnonymousBlockNode
)

// Dimensions - A dimension
type Dimensions struct {
	Content Rect
	Padding EdgeSizes
	Border  EdgeSizes
	Margin  EdgeSizes
}

// Rect - A rectangle
type Rect struct {
	X      float32
	Y      float32
	width  float32
	height float32
}

// EdgeSizes - An edge size
type EdgeSizes struct {
	left   float32
	right  float32
	top    float32
	bottom float32
}

// Box - A layout box
type Box struct {
	Dimensions Dimensions
	BoxType    BoxType
	Children   []Box
}

// BuildLayoutTree - Build a layout tree
func BuildLayoutTree(renderedNode render.RenderedNode) Box {
	var root Box
	switch renderedNode.GetDisplay() {
	case render.BLOCK:
		root = createBlockNode(renderedNode)
	case render.INLINE:
		root = createInlineNode(renderedNode)
	case render.NONE:
		panic("Root node has display none")
	}
	var lastType render.Display
	var anonymousLayoutBox *Box
	for _, child := range renderedNode.Children {
		switch child.GetDisplay() {
		case render.BLOCK:
			if anonymousLayoutBox != nil {
				root.Children = append(root.Children, *anonymousLayoutBox)
				anonymousLayoutBox = nil
			}
			root.Children = append(root.Children, BuildLayoutTree(child))
			lastType = render.BLOCK
		case render.INLINE:
			if lastType != render.INLINE {
				anonymousLayoutBox = createAnonymousBlockNode()
			}
			anonymousLayoutBox.Children = append(anonymousLayoutBox.Children, BuildLayoutTree(child))
			lastType = render.INLINE
		}
	}
	if anonymousLayoutBox != nil {
		root.Children = append(root.Children, *anonymousLayoutBox)
		anonymousLayoutBox = nil
	}
	return root
}

func createBlockNode(r render.RenderedNode) Box {
	var initialDimensions Dimensions
	return Box{
		BoxType:    BlockNode,
		Dimensions: initialDimensions,
		Children:   []Box{},
	}
}

func createInlineNode(r render.RenderedNode) Box {
	var initialDimensions Dimensions
	return Box{
		BoxType:    InlineNode,
		Dimensions: initialDimensions,
		Children:   []Box{},
	}
}

func createAnonymousBlockNode() *Box {
	var initialDimensions Dimensions
	return &Box{
		BoxType:    AnonymousBlockNode,
		Dimensions: initialDimensions,
		Children:   []Box{},
	}
}
