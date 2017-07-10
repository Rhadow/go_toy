package layout

// BlockNode - A block node constant
const BlockNode = "BLOCK_NODE"

// InlineNode - A inline node constant
const InlineNode = "INLINE_NODE"

// AnonymousBlockNode - An anonymoous node constant
const AnonymousBlockNode = "ANONYMOUS_BLOCK_NODE"

// BLOCK - Block for display
const BLOCK = "DISPLAY_BLOCK"

// INLINE - Inline for display
const INLINE = "DISPLAY_INLINE"

// NONE - None for display
const NONE = "DISPLAY_NONE"

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
	BoxType    string
	Children   []Box
}
