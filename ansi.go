package main

import (
	"fmt"
	"strings"
)

// ANSI for curl / terminals that interpret escape codes (Windows Terminal, iTerm, etc.).
const ansiReset = "\033[0m"

func ansi256(fg int, s string) string {
	return fmt.Sprintf("\033[38;5;%dm%s%s", fg, s, ansiReset)
}

func ansiBold256(fg int, s string) string {
	return fmt.Sprintf("\033[1;38;5;%dm%s%s", fg, s, ansiReset)
}

func ansiHR(width int, label string) string {
	if width < 8 {
		width = 52
	}
	inner := width - len(label) - 2
	if inner < 4 {
		inner = 4
	}
	left := inner / 2
	right := inner - left
	return ansi256(109, strings.Repeat("─", left)) +
		ansiBold256(213, " "+label+" ") +
		ansi256(109, strings.Repeat("─", right))
}

// ansiHero two-column layout: left panel (teal), HELLO (magenta), world (green), cat (peach).
func ansiHero() string {
	rs := splitLines(heroRightBlock())
	ls := splitLines(ASCIILeftPanel)
	maxW := 0
	for _, s := range ls {
		if w := runeLen(s); w > maxW {
			maxW = w
		}
	}
	g := strings.Repeat(" ", 4)
	n := max(len(ls), len(rs))
	var b strings.Builder
	for i := 0; i < n; i++ {
		var l string
		if i < len(ls) {
			l = ls[i]
		}
		pad := maxW - runeLen(l)
		if pad < 0 {
			pad = 0
		}
		leftPad := l + strings.Repeat(" ", pad) + g
		b.WriteString(ansi256(109, leftPad))
		if i < len(rs) {
			r := rs[i]
			switch {
			case i < 5:
				b.WriteString(ansiBold256(213, r))
			case i == 5:
				b.WriteString(ansi256(86, r))
			default:
				b.WriteString(ansi256(215, r))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func ansiCurlPage() string {
	var b strings.Builder

	b.WriteString(ansiHero())
	b.WriteString("\n")

	b.WriteString(ansiBold256(213, "  ✦  "+welcomeTitle))
	b.WriteString("\n")
	for _, line := range strings.Split(welcomeBody, "\n") {
		b.WriteString(ansi256(252, "     "+line))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	b.WriteString(ansiHR(52, "About"))
	b.WriteString("\n")
	b.WriteString(ansi256(252, aboutText))
	b.WriteString("\n\n")

	b.WriteString(ansiHR(52, "Contact"))
	b.WriteString("\n")
	b.WriteString(ansi256(246, "     Inside this section — links you can copy:\n\n"))
	b.WriteString(fmt.Sprintf("     %s  %s\n", ansiBold256(84, "● github"), ansi256(246, GitHubURL)))
	b.WriteString(fmt.Sprintf("     %s  %s\n", ansiBold256(33, "● linkedin"), ansi256(246, LinkedInURL)))
	b.WriteString(fmt.Sprintf("     %s  %s\n", ansiBold256(117, "● twitter"), ansi256(246, TwitterURL)))
	b.WriteString("\n")
	b.WriteString(ansi256(240, strings.Repeat("·", 56)))
	b.WriteString("\n")
	b.WriteString(ansi256(245, "     ") + ansi256(240, "Colors: use a modern terminal. "))
	b.WriteString(ansi256(245, "Full ↑/↓ UI: "))
	b.WriteString(ansiBold256(86, "go run ."))
	b.WriteString(ansi256(245, " (omit -serve)."))
	b.WriteString("\n")

	return b.String()
}
