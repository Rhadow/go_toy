package main

import (
	"fmt"

	"github.com/rhadow/go_toy/browser_engine/layout"
	"github.com/rhadow/go_toy/browser_engine/parser"
	"github.com/rhadow/go_toy/browser_engine/render"
)

func main() {
	sourceHTML := `<div class="wrapper">
		<div id="test1">Hi1</div>
		<div id="test2">Hi2</div>
		<div class="test3">Hi3</div>
		<div class="test4">Hi4</div>
		<div class="test3">Hi5</div>
		<div class="test3">Hi6</div>
	</div>`
	htmlParser := parser.HTMLParser{
		Input:    sourceHTML,
		Position: 0,
	}
	rootNode := htmlParser.Parse()

	sourceCSS := `
	.wrapper {
		display: block;
	}
	#test1 {
		display: block;
	}
	#test2 {
		display: inline;
	}
	.test3 {
		display: inline;
	}
	.test4 {
		display: block;
	}`
	cssParser := parser.CSSParser{
		Input:    sourceCSS,
		Position: 0,
	}
	styleSheet := cssParser.Parse()

	renderTree := render.BuildRenderTree(rootNode, styleSheet)
	fmt.Println(renderTree)
	fmt.Println()

	layoutTree := layout.BuildLayoutTree(renderTree)
	fmt.Println(layoutTree)
}
