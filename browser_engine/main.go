package main

import (
	"fmt"

	"github.com/rhadow/go_toy/browser_engine/parser"
)

func main() {
	sourceHTML := `<html>
	<body>
        <h1>Title</h1>
        <div id="main" class="test">
            <p>Hello <em>world</em>!</p>
        </div>
    </body>
</html>`

	htmlParser := parser.HTMLParser{
		Input:    sourceHTML,
		Position: 0,
	}

	rootNode := htmlParser.Parse()
	fmt.Println(rootNode)

	sourceCSS := `
	h1 {
		display: inline-block;
		margin-top: 50px;
		background-color: #01cafe;
	}
	.test {
		border-radius: 5px;
	}
	#test {
		position: absolute;
		color: #CC0000;
	}
	h1.test {color: #000000;}
	div#test {
		padding-right: 20px;
	}
	div.test.test2#test {
		color: #00CC00;
	}`

	cssParser := parser.CSSParser{
		Input:    sourceCSS,
		Position: 0,
	}

	styleSheet := cssParser.Parse()
	fmt.Println(styleSheet)
}
