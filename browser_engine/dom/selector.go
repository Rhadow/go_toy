package dom

import (
	"fmt"
	"strings"
)

// Selector - A CSS selector
type Selector interface {
	getSelectorValue() string
}

// SimpleSelector - A simple css selector. E.g: #testId.testClass
type SimpleSelector struct {
	TagName string
	ID      string
	Classes []string
}

func (s SimpleSelector) getSelectorValue() string {
	result := ""
	if s.TagName != "" {
		result += fmt.Sprintf("Tag name: %s\n", s.TagName)
	}
	if s.ID != "" {
		result += fmt.Sprintf("ID: %s\n", s.ID)
	}
	result += "Class: "
	result += strings.Join(s.Classes, ",")
	return result
}
