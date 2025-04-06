package parser

import (
	"io"
	"mado/renderer/jira"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func ToHtml(source []byte, w io.Writer) {
	if err := goldmark.Convert(source, w); err != nil {
		panic(err)
	}
}

func ToJira(source []byte, w io.Writer, language string) {
	md := goldmark.New(goldmark.WithRenderer(renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(jira.NewRenderer(jira.Config{
		Language: language,
	}), 1000)))))
	if err := md.Convert(source, w); err != nil {
		panic(err)
	}
}
