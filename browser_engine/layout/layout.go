package layout

// BoxType - Node Type
type BoxType uint8

// BlockNode - A block node constant
const (
	BlockNode BoxType = iota + 1
	InlineNode
	AnonymousBlockMode
)

// Display - Display Type
type Display uint8

// BLOCK - Block for display
const (
	BLOCK Display = iota + 1
	INLINE
	NONE
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
