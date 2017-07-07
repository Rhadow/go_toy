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
}
