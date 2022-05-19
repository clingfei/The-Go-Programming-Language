package ch5

import (
	"fmt"
	"golang.org/x/net/html"
)

func printText(n *html.Node) {
	if n == nil {
		return
	}
	if n.Data != "script" && n.Data != "style" {
		fmt.Println(n.DataAtom.String())
	}
	printText(n.FirstChild)
	printText(n.NextSibling)
}
