package main

import (
	"strings"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type pane int

func sshHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	_, _, active := sess.Pty()
	if !active {
		wish.Fatalln(sess, "no terminal detected")
		return nil, nil
	}
	m := &model{
		cursor:  0,
		cCursor: 0,
	}
	// The wish bubbletea middleware handles window size messages for us
	// but we can also get the initial size from the PTY.
	if pty, _, ok := sess.Pty(); ok {
		m.width = pty.Window.Width
		m.height = pty.Window.Height
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}


const (
	paneHome pane = iota
	paneAbout
	paneContact
)

var homeMenu = []string{
	"About",
	"Contacts",
	"Tetris",
	"Quit",
}

var contactMenu = []string{
	"GitHub",
	"LinkedIn",
	"Twitter",
	"« Back to main menu",
}

// Catppuccin Mocha–inspired palette (readable on dark terminals).
var (
	colRose    = lipgloss.Color("#f5c2e7")
	colSky     = lipgloss.Color("#89dceb")
	colBlue    = lipgloss.Color("#89b4fa")
	colText    = lipgloss.Color("#cdd6f4")
	colSubtext = lipgloss.Color("#a6adc8")
	colGreen   = lipgloss.Color("#a6e3a1")
	colPeach   = lipgloss.Color("#fab387")
	colMauve   = lipgloss.Color("#cba6f7")
)

type model struct {
	pane    pane
	cursor  int
	cCursor int
	width   int
	height  int
	status  string
}

func runTUI() error {
	p := tea.NewProgram(&model{cursor: 0, cCursor: 0}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.pane {
		case paneHome:
			return m.updateHome(msg)
		case paneAbout:
			if key := msg.String(); key == "b" || key == "esc" {
				m.pane = paneHome
				m.status = ""
				return m, nil
			}
		case paneContact:
			return m.updateContact(msg)
		}
	}
	return m, nil
}

func (m *model) updateHome(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(homeMenu)-1 {
			m.cursor++
		}
	case "enter":
		switch m.cursor {
		case 0:
			m.pane = paneAbout
		case 1:
			m.pane = paneContact
			m.cCursor = 0
			m.status = ""
		case 2:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) updateContact(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cCursor > 0 {
			m.cCursor--
		}
	case "down", "j":
		if m.cCursor < len(contactMenu)-1 {
			m.cCursor++
		}
	case "enter":
		switch m.cCursor {
		case 0:
			_ = openURL(GitHubURL)
			m.status = "Opened GitHub in your browser."
		case 1:
			_ = openURL(LinkedInURL)
			m.status = "Opened LinkedIn in your browser."
		case 2:
			_ = openURL(TwitterURL)
			m.status = "Opened Twitter / X in your browser."
		case 3:
			m.pane = paneHome
			m.status = ""
		}
	case "b", "esc":
		m.pane = paneHome
		m.status = ""
	}
	return m, nil
}

func (m *model) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(colRose).
		Render("ayushkashyap.me")
	sub := lipgloss.NewStyle().Foreground(colSubtext).Render("terminal portfolio")

	switch m.pane {
	case paneAbout:
		return m.viewAbout(title, sub)
	case paneContact:
		return m.viewContact(title, sub)
	case paneAbout:
		return m.viewTetris(title,sub)
	default:
		return m.viewHome(title, sub)
	}
}

func (m *model) heroBlock() string {
	left := lipgloss.NewStyle().Foreground(colSky).Render(ASCIILeftPanel)
	hello := lipgloss.NewStyle().Bold(true).Foreground(colMauve).Render(strings.TrimSpace(ASCIIHelloLetters))
	world := lipgloss.NewStyle().Foreground(colGreen).Render(strings.TrimSpace(ASCIIWorldLine))
	cat := lipgloss.NewStyle().Foreground(colPeach).Render(strings.TrimSpace(CatASCII))
	worldCat := lipgloss.JoinVertical(lipgloss.Left, hello, world, "", cat)
	gap := lipgloss.NewStyle().Render("   ")

	if m.width >= 74 {
		return lipgloss.JoinHorizontal(lipgloss.Top, left, gap, worldCat)
	}
	return lipgloss.JoinVertical(lipgloss.Left, left, "", worldCat)
}

func (m *model) viewTetris(title,sub string) string {
	
}

func (m *model) viewHome(title, sub string) string {
	heroFramed := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colSky).
		Padding(0, 1).
		Render(m.heroBlock())

	helloLines := []string{
		lipgloss.NewStyle().Bold(true).Foreground(colMauve).Render("✦ " + welcomeTitle),
		lipgloss.NewStyle().Foreground(colText).Width(max(20, m.width-8)).Render(welcomeBody),
	}
	helloBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colBlue).
		Padding(0, 1).
		Render(strings.Join(helloLines, "\n"))

	head := lipgloss.JoinVertical(lipgloss.Left, title, sub, "")
	blk := lipgloss.JoinVertical(lipgloss.Left, heroFramed, "", helloBox, "")

	menuTitle := lipgloss.NewStyle().Foreground(colPeach).Bold(true).Render("Menu")
	normal := lipgloss.NewStyle().Foreground(colText)
	sel := lipgloss.NewStyle().Foreground(colGreen).Bold(true)

	var rows []string
	for i, item := range homeMenu {
		prefix := "  "
		st := normal
		if i == m.cursor {
			prefix = "› "
			st = sel
		}
		rows = append(rows, st.Render(prefix+item))
	}
	menuBlock := strings.Join(rows, "\n")

	footer := lipgloss.NewStyle().Foreground(colSubtext).Render("↑/↓ or j/k · Enter · q quit")

	var parts []string
	parts = append(parts, head, blk, menuTitle, menuBlock)
	if m.status != "" {
		parts = append(parts, "", lipgloss.NewStyle().Foreground(colPeach).Render(m.status))
	}
	parts = append(parts, "", footer)
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m *model) viewAbout(title, sub string) string {
	box := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colMauve).
		Padding(0, 1).
		Width(max(20, m.width-4)).
		Foreground(colText).
		Render(aboutText)
	help := lipgloss.NewStyle().Foreground(colSubtext).Render("b or esc — back")
	return lipgloss.JoinVertical(lipgloss.Left,
		title, sub, "",
		lipgloss.NewStyle().Bold(true).Foreground(colPeach).Render("About"),
		"", box, "", help,
	)
}

func (m *model) viewContact(title, sub string) string {
	intro := lipgloss.NewStyle().Foreground(colText).Render("Choose a platform — Enter opens in your browser.")

	normal := lipgloss.NewStyle().Foreground(colText)
	sel := lipgloss.NewStyle().Foreground(colGreen).Bold(true)

	var rows []string
	for i, item := range contactMenu {
		prefix := "  "
		st := normal
		if i == m.cCursor {
			prefix = "› "
			st = sel
		}
		line := st.Render(prefix + item)
		switch i {
		case 0:
			line = lipgloss.JoinHorizontal(lipgloss.Left, st.Render(prefix+item),
				lipgloss.NewStyle().Foreground(colSubtext).Render("   "+GitHubURL))
		case 1:
			line = lipgloss.JoinHorizontal(lipgloss.Left, st.Render(prefix+item),
				lipgloss.NewStyle().Foreground(colSubtext).Render("   "+LinkedInURL))
		case 2:
			line = lipgloss.JoinHorizontal(lipgloss.Left, st.Render(prefix+item),
				lipgloss.NewStyle().Foreground(colSubtext).Render("   "+TwitterURL))
		}
		rows = append(rows, line)
	}

	menuBlock := strings.Join(rows, "\n")
	help := lipgloss.NewStyle().Foreground(colSubtext).Render("b or esc — back · q quit")

	var extra string
	if m.status != "" {
		extra = "\n" + lipgloss.NewStyle().Foreground(colSky).Render(m.status)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		title, sub, "",
		lipgloss.NewStyle().Bold(true).Foreground(colPeach).Render("Contacts"),
		"", intro, "", menuBlock, extra, "", help,
	)
}
