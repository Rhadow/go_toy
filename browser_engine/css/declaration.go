package css

import "fmt"

// StyleDeclaration - A StyleDeclaration
type StyleDeclaration struct {
	Name  string
	Value StyleDeclarationValue
}

// StyleDeclarationValue - A style declaration value
type StyleDeclarationValue interface {
	getStyleDeclarationValue() string
}

// Length - A length value. e.g: 50px
type Length struct {
	Value float32
	Unit  string
}

func (l Length) getStyleDeclarationValue() string {
	return fmt.Sprintf("%f %s", l.Value, l.Unit)
}

// Keyword - A keyword value. e.g: inline-block
type Keyword string

func (k Keyword) getStyleDeclarationValue() string {
	return string(k)
}

// ColorValue - A keyword value. e.g: rgba(255,255,255,0)
type ColorValue struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func (k ColorValue) getStyleDeclarationValue() string {
	return fmt.Sprintf("r:%d, g:%d, b:%d, a:%d", k.R, k.G, k.B, k.A)
}
