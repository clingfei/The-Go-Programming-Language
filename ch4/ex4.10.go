package ch4

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	LTONEM = "less than one month"
	LTONEY = "less than one year"
	MTONEY = "more than one year"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	//QueryEscape对参数进行转码以安全使用参数
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}

	resp.Body.Close()
	return &result, nil
}

func main() {
	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
	categoryMap := make(map[string][]Issue)
	for _, item := range result.Items {
		y, m, _ := item.CreatedAt.Date()
		curY, curM, _ := time.Now().Date()
		if curY == y && m == curM {
			categoryMap[LTONEM] = append(categoryMap[LTONEM], *item)
		} else if curY == y || curY-y == 1 && m >= curM {
			categoryMap[LTONEY] = append(categoryMap[LTONEY], *item)
		} else {
			categoryMap[MTONEY] = append(categoryMap[MTONEY], *item)
		}
	}
	fmt.Println()
	for c, category := range categoryMap {
		fmt.Println(c)
		for _, issue := range category {
			fmt.Printf("#%-5d %9.9s %.55s\n",
				issue.Number, issue.User.Login, issue.Title)
		}
	}
}
