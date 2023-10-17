package domain

import "golang.org/x/net/html"

type Content struct {
	HTML *html.Node
	Text string
}
