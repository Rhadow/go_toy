package css

// Selector - A CSS selector
type Selector interface {
	GetSpecifity() uint32
}

// SimpleSelector - A simple css selector. E.g: #testId.testClass
type SimpleSelector struct {
	ID      string
	Classes []string
	TagName string
}

// GetSpecifity - Get the specificity of this selector
func (s SimpleSelector) GetSpecifity() uint32 {
	result := 0
	if s.ID != "" {
		result += 100
	}
	result += (len(s.Classes) * 10)
	if s.TagName != "" {
		result++
	}
	return uint32(result)
}
