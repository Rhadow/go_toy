package parser

import (
	"encoding/hex"
	"regexp"
	"sort"
	"strconv"
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
		return c == " " || c == "\n" || c == "\t"
	})
}

// Parse one simple selector, e.g.: `type#id.class1.class2.class3`
func (c *CSSParser) parseSimpleSelector() css.SimpleSelector {
	simpleSelector := css.SimpleSelector{
		TagName: "",
		ID:      "",
		Classes: []string{},
	}
	breakOut := false
	for !c.eof() {
		if breakOut {
			break
		}
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
			breakOut = true
		}
	}
	return simpleSelector
}

func (c *CSSParser) parseIdentifier() string {
	return c.consumeWhile(validIdentifierChar)
}

func (c *CSSParser) parseRule() css.Rule {
	return css.Rule{
		Selectors:    c.parseSelectors(),
		Declarations: c.parseDeclarations(),
	}
}

func (c *CSSParser) parseSelectors() []css.Selector {
	selectorsBySpecificity := bySpecificity{}
	breakOut := false
	for {
		if breakOut {
			break
		}
		selectorsBySpecificity = append(selectorsBySpecificity, c.parseSimpleSelector())
		c.consumeWhitespace()
		nextChar := c.nextChar()
		switch {
		case nextChar == ",":
			c.consumeChar()
			c.consumeWhitespace()
		case nextChar == "{":
			breakOut = true
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

func (c *CSSParser) parseDeclarations() []css.StyleDeclaration {
	if c.consumeChar() != "{" {
		panic("Invalid declaration: no opening brace")
	}
	declarations := []css.StyleDeclaration{}

	for {
		c.consumeWhitespace()
		if c.nextChar() == "}" {
			c.consumeChar()
			break
		}
		declarations = append(declarations, c.parseDeclaration())
	}
	return declarations
}

func (c *CSSParser) parseDeclaration() css.StyleDeclaration {
	propertyName := c.parseIdentifier()
	c.consumeWhitespace()
	if c.consumeChar() != ":" {
		panic("Invalid declaration: no colon")
	}
	c.consumeWhitespace()
	value := c.parseValue()
	c.consumeWhitespace()
	if c.consumeChar() != ";" {
		panic("Invalid declaration: no semicolon")
	}
	return css.StyleDeclaration{
		Name:  propertyName,
		Value: value,
	}
}

func (c *CSSParser) parseValue() css.StyleDeclarationValue {
	nextChar := c.nextChar()
	switch nextChar {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return c.parseLength()
	case "#":
		return c.parseColor()
	default:
		return css.Keyword(c.parseIdentifier())
	}
}

func (c *CSSParser) parseLength() css.Length {
	return css.Length{
		Value: c.parseFloat(),
		Unit:  c.parseUnit(),
	}
}

func (c *CSSParser) parseFloat() float32 {
	floatString := c.consumeWhile(func(char string) bool {
		result, err := regexp.MatchString("[0-9]", char)
		if err != nil {
			panic("Parse float error")
		}
		return result
	})
	result, err := strconv.ParseFloat(floatString, 32)
	if err != nil {
		panic("Parse float error")
	}
	return float32(result)
}

func (c *CSSParser) parseUnit() string {
	c.consumeWhitespace()
	unit := strings.ToLower(c.parseIdentifier())
	switch {
	case unit == "px":
		return unit
	default:
		panic("Invalid unit")
	}
}

func (c *CSSParser) parseColor() css.ColorValue {
	if c.consumeChar() != "#" {
		panic("Invalid color value")
	}
	return css.ColorValue{
		R: c.parseHexPair(),
		G: c.parseHexPair(),
		B: c.parseHexPair(),
		A: 255,
	}
}

func (c *CSSParser) parseHexPair() uint8 {
	hexString := c.Input[int(c.Position):int(c.Position+2)]
	c.Position += 2
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		panic("Parse hex pair error")
	}
	return uint8(decoded[0])
}
