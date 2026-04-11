package main

import (
	"fmt"
	"strings"
)

// ANSI for curl / terminals that interpret escape codes (Windows Terminal, iTerm, etc.).
const ansiReset = "\033[0m"
const ansiST = "\033\\" // string terminator for OSC sequences

func ansi256(fg int, s string) string {
	return fmt.Sprintf("\033[38;5;%dm%s%s", fg, s, ansiReset)
}

func ansiBold256(fg int, s string) string {
	return fmt.Sprintf("\033[1;38;5;%dm%s%s", fg, s, ansiReset)
}

// osc8 wraps text as a clickable hyperlink in supporting terminals (OSC 8).
func osc8(url, text string) string {
	return "\033]8;;" + url + ansiST + text + "\033]8;;" + ansiST
}

func ansiHRColor(width int, label string, labelFG int) string {
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
		ansiBold256(labelFG, " "+label+" ") +
		ansi256(109, strings.Repeat("─", right))
}

// ansiResumeHero — large name art with a short tagline.
func ansiResumeHero() string {
	lines := append(splitLines(asciiNameAyush), splitLines(asciiNameKashyap)...)
	// Cool → warm as you read down
	palette := []int{45, 39, 38, 37, 43, 214, 220}
	var b strings.Builder
	for i, line := range lines {
		c := palette[i%len(palette)]
		b.WriteString(ansiBold256(c, line))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(ansi256(246, "  "))
	b.WriteString(ansiBold256(86, "Software engineer"))
	b.WriteString(ansi256(246, "  ·  "))
	b.WriteString(ansi256(252, Location))
	b.WriteString("\n")
	return b.String()
}

func writePrefixedLines(b *strings.Builder, prefix string, fg int, text string) {
	for _, line := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		b.WriteString(ansi256(fg, prefix+line))
		b.WriteString("\n")
	}
}

func ansiCurlPage() string {
	var b strings.Builder

	b.WriteString(ansiResumeHero())
	b.WriteString("\n")

	b.WriteString(ansiBold256(213, "  ✦  "+welcomeTitle))
	b.WriteString("\n")
	for _, line := range strings.Split(welcomeBody, "\n") {
		b.WriteString(ansi256(252, "     "+line))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	b.WriteString(ansiHRColor(56, "Introduction", 117))
	b.WriteString("\n")
	writePrefixedLines(&b, "     ", 252, introBlurb)
	b.WriteString("\n")

	b.WriteString(ansiHRColor(56, "Education", 75))
	b.WriteString("\n")
	writePrefixedLines(&b, "     ", 252, educationBlock)
	b.WriteString("\n")

	b.WriteString(ansiHRColor(56, "Technical skills", 141))
	b.WriteString("\n")
	writePrefixedLines(&b, "     ", 252, skillsBlock)
	b.WriteString("\n")

	b.WriteString(ansiHRColor(56, "Experience", 114))
	b.WriteString("\n")
	writePrefixedLines(&b, "     ", 252, experienceBlock)
	b.WriteString("\n")

	b.WriteString(ansiHRColor(56, "Projects", 220))
	b.WriteString("\n")
	writePrefixedLines(&b, "     ", 252, projectsBlock)
	b.WriteString("\n")
	b.WriteString(ansi256(246, "     "))
	b.WriteString(ansiBold256(220, "Links  "))
	b.WriteString(ansiBold256(117, osc8(GoChessDemoURL, "GoChess (demo / repo)")))
	b.WriteString(ansi256(240, "   "))
	b.WriteString(ansiBold256(117, osc8(PyGitURL, "PyGit")))
	b.WriteString(ansi256(240, "   "))
	b.WriteString(ansiBold256(117, osc8(PacmanRLURL, "Pacman RL")))
	b.WriteString("\n\n")

	b.WriteString(ansiHRColor(56, "Contact", 86))
	b.WriteString("\n")
	b.WriteString(ansi256(252, "     "))
	b.WriteString(ansiBold256(117, osc8(GmailURL, "Email")))
	b.WriteString(ansi256(252, "     "))
	b.WriteString(ansiBold256(84, osc8(GitHubURL, "GitHub")))
	b.WriteString(ansi256(240, "  ·  "))
	b.WriteString(ansiBold256(33, osc8(LinkedInURL, "LinkedIn")))
	b.WriteString(ansi256(240, "  ·  "))
	b.WriteString(ansiBold256(117, osc8(TwitterURL, "X / Twitter")))
	b.WriteString("\n")
	b.WriteString(ansi256(252, "     "))
	b.WriteString(ansi256(240, "Tip: clickable links need a terminal that supports OSC 8 hyperlinks."))
	b.WriteString("\n\n")

	b.WriteString(ansi256(240, strings.Repeat("·", 56)))
	b.WriteString("\n")
	b.WriteString(ansi256(245, "     "))
	b.WriteString(ansi256(240, "Interactive TUI (arrow keys): "))
	b.WriteString(ansiBold256(86, "go run ."))
	b.WriteString(ansi256(245, "  ·  "))
	b.WriteString(ansi256(240, "Same ANSI page: "))
	b.WriteString(ansiBold256(213, "curl -sL https://"+SiteDomain+"/"))
	b.WriteString(ansi256(245, " or "))
	b.WriteString(ansiBold256(213, "curl -sL https://"+SiteDomain+"/terminal"))
	b.WriteString("\n")
	b.WriteString(ansi256(245, "     "))
	b.WriteString(ansi256(240, "SSH shell coming soon: "))
	b.WriteString(ansiBold256(86, "ssh terminal.AyushKashyap.me"))
	b.WriteString(ansi256(240, " (placeholder — wire when ready)."))
	b.WriteString("\n")

	return b.String()
}
