package render

import (
	"errors"
	"sort"

	"github.com/rhadow/go_toy/browser_engine/css"
	"github.com/rhadow/go_toy/browser_engine/dom"
)

// Display - Display Type
type Display uint8

// BLOCK - Block for display
const (
	BLOCK Display = iota + 1
	INLINE
	NONE
)

// MatchedRule - A map with specificity and rule
type MatchedRule struct {
	Specificity uint32
	Rule        css.Rule
}

// MatchedRuleBySpecificity - For sortng matched rule by specificity
type MatchedRuleBySpecificity []MatchedRule

func (m MatchedRuleBySpecificity) Len() int {
	return len(m)
}

func (m MatchedRuleBySpecificity) Less(i, j int) bool {
	return m[i].Specificity < m[j].Specificity
}

func (m MatchedRuleBySpecificity) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// PropertyMap - Map from CSS property names to values
type PropertyMap map[string]css.StyleDeclarationValue

// RenderedNode - A node with associated style data.
type RenderedNode struct {
	Node       dom.Node
	Properties PropertyMap
	Children   []RenderedNode
}

func (r RenderedNode) getPropertyValueByName(name string) (css.StyleDeclarationValue, error) {
	var result css.StyleDeclarationValue
	value, ok := r.Properties[name]
	if !ok {
		return result, errors.New("No name field in given rendered node")
	}
	return value, nil
}

// GetDisplay - get display of render node
func (r RenderedNode) GetDisplay() Display {
	styleDeclaration, err := r.getPropertyValueByName("display")
	if err != nil {
		return INLINE
	}
	value := styleDeclaration.GetStyleDeclarationValue()
	switch value {
	case "block":
		return BLOCK
	case "none":
		return NONE
	default:
		return INLINE
	}
}

// MatchSelector - Check selector matches selected element
func MatchSelector(element dom.ElementNode, selector css.Selector) bool {
	if s, ok := selector.(css.SimpleSelector); ok {
		return MatchSimpleSelector(element, css.SimpleSelector(s))
	}
	return false
}

// MatchSimpleSelector - Check a simple selector matches selected element
func MatchSimpleSelector(element dom.ElementNode, selector css.SimpleSelector) bool {
	if selector.TagName != "" {
		if element.TagName != selector.TagName {
			return false
		}
	}
	if selector.ID != "" {
		if id, err := element.GetID(); err != nil || id != selector.ID {
			return false
		}
	}

	if len(selector.Classes) > 0 {
		classes, err := element.GetClasses()
		if err != nil && len(selector.Classes) != 0 {
			return false
		}
		for _, selectorClass := range selector.Classes {
			match := false
			for _, elementClass := range classes {
				if elementClass == selectorClass {
					match = true
				}
			}
			if !match {
				return false
			}
		}
	}
	return true
}

func matchRule(element dom.ElementNode, rule css.Rule) (MatchedRule, error) {
	for _, selector := range rule.Selectors {
		if MatchSelector(element, selector) {
			return MatchedRule{
				Specificity: selector.GetSpecifity(),
				Rule:        rule,
			}, nil
		}
	}
	return MatchedRule{}, errors.New("No matched rule")
}

func matchRules(element dom.ElementNode, styleSheet css.StyleSheet) MatchedRuleBySpecificity {
	result := MatchedRuleBySpecificity{}
	for _, rule := range styleSheet.Rules {
		if matchedRule, err := matchRule(element, rule); err == nil {
			result = append(result, matchedRule)
		}
	}
	return result
}

func getPropertyMapForElement(element dom.ElementNode, styleSheet css.StyleSheet) PropertyMap {
	result := PropertyMap{}
	matchedRules := matchRules(element, styleSheet)
	sort.Sort(matchedRules)
	for _, matchedRule := range matchedRules {
		for _, declaration := range matchedRule.Rule.Declarations {
			result[declaration.Name] = declaration.Value
		}
	}
	return result
}

// BuildRenderTree - Builds render tree for given dom and stylesheet
func BuildRenderTree(element dom.Node, styleSheet css.StyleSheet) RenderedNode {
	childrenRenderTree := []RenderedNode{}
	for _, childElement := range element.GetChildren() {
		childrenRenderTree = append(childrenRenderTree, BuildRenderTree(childElement, styleSheet))
	}
	propertyMap := PropertyMap{}
	if elementNode, ok := element.(dom.ElementNode); ok {
		propertyMap = getPropertyMapForElement(elementNode, styleSheet)
	}
	return RenderedNode{
		Node:       element,
		Properties: propertyMap,
		Children:   childrenRenderTree,
	}
}
