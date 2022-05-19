package ch5

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
)

func count(n *html.Node, nodeMap *map[string]int) {
	if n == nil {
		return
	}
	(*nodeMap)[n.Data]++
	count(n.FirstChild, nodeMap)
	count(n.NextSibling, nodeMap)
}

func main() {
	nodeMap := make(map[string]int)
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	count(doc, &nodeMap)
	for k, v := range nodeMap {
		if v > 1 {
			fmt.Printf("%s: %d", k, v)
		}
	}
}
