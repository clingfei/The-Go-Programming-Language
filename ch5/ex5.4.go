package ch5

import "golang.org/x/net/html"

func visit1(links []string, images []string, scripts []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			} else if a.Key == "images" {
				images = append(images, a.Val)
			} else if a.Key == "scripts" {
				scripts = append(scripts, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
