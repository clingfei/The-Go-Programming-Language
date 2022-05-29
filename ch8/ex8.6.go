package ch8

import (
	"flag"
	"fmt"
	"gopl.io/ch5/links"
	"log"
	"os"
)

var depth = flag.Int("depth", -1, "depth")

func crawl(url string) []string {
	fmt.Println(url)
	// 使用有容量限制的buffer channel来限制并发数，同一时刻只能有20个连接
	tokens <- struct{}{}
	list, err := links.Extract(url)
	<-tokens
	if err != nil {
		log.Print(err)
	}
	return list
}

var tokens = make(chan struct{}, 20)

func main() {
	worklist := make(chan []string)

	go func() { worklist <- os.Args[1:] }()

	seen := make(map[string]bool)
	if *depth == -1 {
		for list := range worklist {
			for _, link := range list {
				if !seen[link] {
					seen[link] = true
					go func(link string) {
						worklist <- crawl(link)
					}(link)
				}
			}
		}
	} else {
		for i := 0; i < *depth; i++ {
			lists := worklist
			for list := range lists {
				for _, link := range list {
					if !seen[link] {
						seen[link] = true
						go func(link string) {
							worklist <- crawl(link)
						}(link)
					}
				}
			}
		}
	}

}
