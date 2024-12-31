package main

import (
	"net/url"
	"strings"
  "fmt"

	"golang.org/x/net/html"
)

func ReadLinks(input string, baseURL *url.URL) ([]string, error) {
	tokenizer := html.NewTokenizer(strings.NewReader(input))

	links := make([]string, 0)
	for tokenizer.Err() == nil {
		token := tokenizer.Next()

		if token != html.StartTagToken {
			continue
		}

		t := tokenizer.Token()

		isAnchor := t.Data == "a"
		if !isAnchor {
			continue
		}

		for _, attr := range t.Attr {
			if attr.Key != "href" || attr.Val == "" {
				continue
			}


      //  NOTE: Even better, baseURL.ResolveReference(href)
      // u := attr.Val
      // if !strings.HasPrefix(u, "http") {
      //   u, err = url.JoinPath(baseURL.String(), u)
      //   if err != nil {
      //     return nil, fmt.Errorf("failed to join url %w", err)  
      //   }
      // }

      href, err := url.Parse(attr.Val)
      if err != nil {
        return nil, fmt.Errorf("failed to parse url %w", err)
      }

      u := baseURL.ResolveReference(href).String()
      links = append(links, u)
		}
	}

	return links, nil
}
