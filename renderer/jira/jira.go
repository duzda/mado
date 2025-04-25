package jira

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Config struct {
	Language string
}

func (r Config) SetJiraOption(c *Config) { *c = r }

type Renderer struct {
	Config Config
}

type Option interface {
	SetJiraOption(*Config)
}

func NewRenderer(opts ...Option) renderer.NodeRenderer {
	r := &Renderer{
		Config: Config{},
	}
	for _, opt := range opts {
		opt.SetJiraOption(&r.Config)
	}
	return r
}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindString, r.renderString)
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)
}

func (r *Renderer) writeLines(w util.BufWriter, source []byte, node ast.Node) {
	l := node.Lines().Len()
	for i := 0; i < l; i++ {
		line := node.Lines().At(i)
		_, _ = w.Write(line.Value(source))
	}
}

func (r *Renderer) putSpaceIfMissing(w util.BufWriter, source []byte, node ast.Node) {
	s := node.PreviousSibling()
	if s != nil && s.Kind() == ast.KindText {
		seg := s.(*ast.Text).Segment
		previous := seg.Value(source)
		if len(previous) > 0 && previous[len(previous)-1] != ' ' {
			_ = w.WriteByte(' ')
		}
	}
}

func (r *Renderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.AutoLink)
	if !entering {
		return ast.WalkContinue, nil
	}
	_ = w.WriteByte('[')
	url := n.URL(source)
	label := n.Label(source)
	if n.AutoLinkType == ast.AutoLinkEmail && !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
		_, _ = w.WriteString("mailto:")
	}
	if len(label) > 0 {
		_, _ = w.Write(label)
		_ = w.WriteByte('|')
	}
	_, _ = w.Write(url)
	_ = w.WriteByte(']')
	return ast.WalkContinue, nil
}

func (r *Renderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("\n{quote}\n")
	} else {
		_, _ = w.WriteString("{quote}\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("\n{code:" + r.Config.Language + "}\n")
		r.writeLines(w, source, node)
	} else {
		_, _ = w.WriteString("{code}\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.putSpaceIfMissing(w, source, node)
		// _, _ = w.WriteString("\n{code:" + r.Config.Language + "}\n")
		_, _ = w.WriteString("{{")
	} else {
		// _, _ = w.WriteString("\n{code}\n")
		_, _ = w.WriteString("}}")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkStop, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	tag := "_"
	if n.Level == 2 {
		tag = "*"
	}
	if entering {
		r.putSpaceIfMissing(w, source, node)
	}
	_, _ = w.WriteString(tag)
	return ast.WalkContinue, nil
}

func (r *Renderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	language := n.Language(source)
	if language == nil {
		language = []byte(r.Config.Language)
	}
	if entering {
		_, _ = w.WriteString("\n{code:" + string(language) + "}\n")
		r.writeLines(w, source, node)
	} else {
		_, _ = w.WriteString("{code}\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.writeLines(w, source, node)
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		if n.PreviousSibling() != nil {
			_ = w.WriteByte('\n')
		}
		_ = w.WriteByte('h')
		_ = w.WriteByte("0123456"[n.Level])
		_, _ = w.WriteString(". ")
	} else {
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Image)
	if !entering {
		_ = w.WriteByte('!')
		return ast.WalkContinue, nil
	}
	s := n.PreviousSibling()
	if s != nil {
		ss := s.PreviousSibling()
		if ss != nil && ss.Kind() != ast.KindLink {
			_ = w.WriteByte('\n')
		}
	}
	_ = w.WriteByte('!')
	if n.Destination != nil {
		_, _ = w.Write(n.Destination)
		_ = w.WriteByte('|')
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if !entering {
		if n.Destination != nil {
			_ = w.WriteByte('|')
			_, _ = w.Write(n.Destination)
		}
		_ = w.WriteByte(']')
		return ast.WalkContinue, nil
	}
	if entering {
		r.putSpaceIfMissing(w, source, node)
	}
	_ = w.WriteByte('[')
	return ast.WalkContinue, nil
}

func (r *Renderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	p := node.Parent().(*ast.List)
	tag := "*"
	if p.IsOrdered() {
		tag = "#"
	}
	if entering {
		_, _ = w.WriteString(tag + " ")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Paragraph)
	if entering && n.PreviousSibling() == nil {
		return ast.WalkContinue, nil
	}
	if !n.IsRaw() {
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	r.writeLines(w, source, node)
	return ast.WalkContinue, nil
}

func (r *Renderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	r.writeLines(w, source, node)
	return ast.WalkContinue, nil
}

func (r *Renderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Text)
	if !entering {
		return ast.WalkContinue, nil
	}
	s := n.PreviousSibling()
	if s != nil {
		ss := s.PreviousSibling()
		if ss != nil && ss.Kind() == ast.KindImage {
			_ = w.WriteByte('\n')
		} else if s.Kind() == ast.KindText {
			seg := s.(*ast.Text).Segment
			previous := seg.Value(source)
			p := n.Parent()
			if p != nil {
				pp := p.Parent()
				if pp != nil && pp.Kind() == ast.KindBlockquote {
					goto render_segment
				}
			}
			if previous != nil && previous[len(previous)-1] != ' ' {
				_ = w.WriteByte(' ')
			}
		}
	}
render_segment:
	segment := n.Segment
	_, _ = w.Write(segment.Value(source))
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTextBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	_ = w.WriteByte('\n')
	return ast.WalkContinue, nil
}
