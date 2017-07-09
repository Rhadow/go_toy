package main

import (
	"fmt"

	"github.com/rhadow/go_toy/browser_engine/parser"
	"github.com/rhadow/go_toy/browser_engine/render"
)

func main() {
	sourceHTML := `<html>
	<body>
        <h1 class="test" id="test2">Title</h1>
				<div class="test">Test</div>
    </body>
</html>`
	htmlParser := parser.HTMLParser{
		Input:    sourceHTML,
		Position: 0,
	}
	rootNode := htmlParser.Parse()

	sourceCSS := `
	h1 {
		display: inline-block;
		margin-top: 50px;
		color: #01cafe;
	}
	#test2 {
		color: #0000CC;
	}
	.test {
		color: #cc0000;
	}`
	cssParser := parser.CSSParser{
		Input:    sourceCSS,
		Position: 0,
	}
	styleSheet := cssParser.Parse()

	renderTree := render.BuildRenderTree(rootNode, styleSheet)
	fmt.Println(renderTree)
}
