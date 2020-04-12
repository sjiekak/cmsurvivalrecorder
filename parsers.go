package main

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func crawl(node *html.Node) (string, error) {
	if node.Data == "span" {
		for _, attr := range node.Attr {
			if attr.Key == "class" && attr.Val == "income" {
				for child := node.FirstChild; child != nil; child = child.NextSibling {
					if child.Type == html.TextNode {
						return child.Data, nil
					}
				}
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			if v, err := crawl(child); err == nil {
				return v, nil
			}
		}
	}
	return "", fmt.Errorf("not found")
}

func parseAmount(s string) (float64, error) {
	for _, substr := range strings.Split(s, "€") {
		if len(substr) > 0 {
			return strconv.ParseFloat(substr, 64)
		}
	}
	return 0.0, fmt.Errorf("empty string")
}
