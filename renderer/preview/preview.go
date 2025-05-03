package preview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
	style    glamour.TermRendererOption
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc || msg.String() == "q" {
			return m, tea.Quit
		}
		if msg.Type == tea.KeyHome || msg.String() == "g" {
			m.viewport.GotoTop()
		}
		if msg.Type == tea.KeyEnd || msg.String() == "G" {
			m.viewport.GotoBottom()
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 1
		}

		content, err := m.getGlamouredContent()
		if err != nil {
			return m, tea.Quit
		}

		m.viewport.SetContent(content)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.viewport.Height == 0 {
		return ""
	}

	var b strings.Builder
	fmt.Fprint(&b, m.viewport.View()+"\n")

	scrollPercent := fmt.Sprintf(" %3.f%% ", m.viewport.ScrollPercent()*100.0)
	padding := strings.Repeat(" ", m.viewport.Width-len(scrollPercent))
	fmt.Fprintf(&b, "%s%s",
		padding,
		scrollPercent,
	)

	return b.String()
}

func (m model) getGlamouredContent() (string, error) {
	r, err := glamour.NewTermRenderer(m.style, glamour.WithWordWrap(m.viewport.Width))
	if err != nil {
		return "", err
	}

	out, err := r.Render(m.content)
	if err != nil {
		return "", err
	}

	return out, nil
}

func RenderPreview(content string, style glamour.TermRendererOption) error {
	if _, err := tea.NewProgram(model{content: content, style: style}, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		panic(err)
	}

	return nil
}
