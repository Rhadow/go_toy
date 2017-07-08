package parser

import (
	"regexp"
	"sort"
	"strings"

	"github.com/rhadow/go_toy/browser_engine/css"
)

type bySpecificity []css.Selector

func (s bySpecificity) Len() int {
	return len(s)
}

func (s bySpecificity) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s bySpecificity) Less(i, j int) bool {
	return s[i].GetSpecifity() > s[j].GetSpecifity()
}

// CSSParser - An CSS parser
type CSSParser struct {
	Position uint64
	Input    string
}

// Read the current character without consuming it.
func (c *CSSParser) nextChar() string {
	return string(c.Input[c.Position])
}

// Do the next characters start with the given string?
func (c *CSSParser) startsWith(str string) bool {
	return strings.HasPrefix(c.Input[c.Position:], str)
}

// Return true if all Input is consumed.
func (c *CSSParser) eof() bool {
	return c.Position >= uint64(len(c.Input))
}

// Return the current character, and advance self.pos to the next character.
func (c *CSSParser) consumeChar() string {
	currentChar := string(c.Input[c.Position])
	c.Position++
	return currentChar
}

// Consume characters until `test` returns false.
func (c *CSSParser) consumeWhile(test func(string) bool) string {
	result := ""
	for !c.eof() && test(c.nextChar()) {
		result += c.consumeChar()
	}
	return result
}

// Consume and discard zero or more whitespace characters.
func (c *CSSParser) consumeWhitespace() {
	c.consumeWhile(func(c string) bool {
		return c == " "
	})
}

// Parse one simple selector, e.g.: `type#id.class1.class2.class3`
func (c *CSSParser) parseSimpleSelector() css.SimpleSelector {
	simpleSelector := css.SimpleSelector{
		TagName: "",
		ID:      "",
		Classes: []string{},
	}
	for !c.eof() {
		nextChar := c.nextChar()
		switch {
		case nextChar == "#":
			c.consumeChar()
			simpleSelector.ID = c.parseIdentifier()
		case nextChar == ".":
			c.consumeChar()
			simpleSelector.Classes = append(simpleSelector.Classes, c.parseIdentifier())
		case validIdentifierChar(nextChar):
			simpleSelector.TagName = c.parseIdentifier()
		default:
			break
		}
	}
	return simpleSelector
}

func (c *CSSParser) parseIdentifier() string {
	return c.consumeWhile(validIdentifierChar)
}

func (c *CSSParser) parseRule() css.Rule {
	return css.Rule{
		Selectors: c.parseSelectors(),
		// TODO: Complete parseDeclarations
		Declarations: c.parseDeclarations(),
	}
}

func (c *CSSParser) parseSelectors() []css.Selector {
	selectorsBySpecificity := bySpecificity{}
	for {
		selectorsBySpecificity = append(selectorsBySpecificity, c.parseSimpleSelector())
		c.consumeWhitespace()
		nextChar := c.nextChar()
		switch {
		case nextChar == ",":
			c.consumeChar()
			c.consumeWhitespace()
		case nextChar == "{":
			break
		default:
			panic("Unexpected character in selector list")
		}
	}
	sort.Sort(selectorsBySpecificity)
	return []css.Selector(selectorsBySpecificity)
}

func validIdentifierChar(char string) bool {
	result, err := regexp.MatchString("[a-zA-Z0-9-_]", char)
	if err != nil {
		panic(err)
	}
	return result
}

// Parse an CSS document and return a css.StyleSheet
func (c *CSSParser) Parse() css.StyleSheet {
	return css.StyleSheet{
		Rules: c.parseRules(),
	}
}

func (c *CSSParser) parseRules() []css.Rule {
	rules := []css.Rule{}
	for !c.eof() {
		c.consumeWhitespace()
		rules = append(rules, c.parseRule())
	}
	return rules
}
