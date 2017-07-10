package css

import "fmt"

// StyleDeclaration - A StyleDeclaration
type StyleDeclaration struct {
	Name  string
	Value StyleDeclarationValue
}

// StyleDeclarationValue - A style declaration value
type StyleDeclarationValue interface {
	GetStyleDeclarationValue() string
}

// Length - A length value. e.g: 50px
type Length struct {
	Value float32
	Unit  string
}

// GetStyleDeclarationValue - Get style declaration value for Length
func (l Length) GetStyleDeclarationValue() string {
	return fmt.Sprintf("%f %s", l.Value, l.Unit)
}

// Keyword - A keyword value. e.g: inline-block
type Keyword string

// GetStyleDeclarationValue - Get style declaration value for Keyword
func (k Keyword) GetStyleDeclarationValue() string {
	return string(k)
}

// ColorValue - A keyword value. e.g: rgba(255,255,255,0)
type ColorValue struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// GetStyleDeclarationValue - Get style declaration value for ColorValue
func (k ColorValue) GetStyleDeclarationValue() string {
	return fmt.Sprintf("r:%d, g:%d, b:%d, a:%d", k.R, k.G, k.B, k.A)
}
