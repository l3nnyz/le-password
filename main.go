package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	accent = lipgloss.Color("#7D56F4")

	countStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(accent).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C9C9C9")).
			Background(lipgloss.Color("#3C3C3C")).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accent).
			Padding(0, 2)

	bulletStyle   = lipgloss.NewStyle().Foreground(accent)
	passwordStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575"))
)

type model struct {
	passwords []string
}

func (m model) Init() tea.Cmd { return tea.Quit }

func (m model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m model) View() string {
	n := len(m.passwords)
	label := "passwords generated"
	if n == 1 {
		label = "password generated"
	}
	header := countStyle.Render(strconv.Itoa(n)) + labelStyle.Render(label)

	lines := make([]string, n)
	for i, p := range m.passwords {
		lines[i] = bulletStyle.Render("› ") + passwordStyle.Render(p)
	}

	return "\n" + header + "\n" + boxStyle.Render(strings.Join(lines, "\n")) + "\n"
}

func main() {
	length := flag.Int("l", 16, "password length")
	count := flag.Int("n", 1, "number of passwords to generate")
	flag.Parse()

	passwords, err := GenerateMany(*length, *count)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(model{passwords: passwords}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
