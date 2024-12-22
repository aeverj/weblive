package utils

import (
	"bytes"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io"
	"strings"
)

func ReadBody(r io.Reader) string {
	content, _ := io.ReadAll(r)
	if len(content) == 0 {
		return ""
	}
	start := bytes.Index(content, []byte("title"))
	if start == -1 {
		start = 0
	}
	e, _, _ := charset.DetermineEncoding(content, "")
	if e == nil && content == nil {
		return ""
	} else {
		ht, _, _ := transform.Bytes(e.NewDecoder(), content)
		return string(ht)
	}
}

func GetTitle(body string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(body))
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return ""
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "title" {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					title := tokenizer.Token()
					return strings.TrimSpace(title.Data)
				}
			}
		}
	}
}
