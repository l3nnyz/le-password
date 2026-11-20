package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

// Maximum number of passwords per run: displaying more than a hundred
// in the console is unreadable, a file can hold more.
const (
	maxConsoleCount = 100
	maxExportCount  = 1000
)

type model struct {
	passwords  []string
	exportPath string
}

func (m model) Init() tea.Cmd { return tea.Quit }

func (m model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m model) View() string {
	n := len(m.passwords)
	noun := "passwords"
	if n == 1 {
		noun = "password"
	}

	if m.exportPath != "" {
		header := countStyle.Render(strconv.Itoa(n)) + labelStyle.Render(noun+" exported to")
		return "\n" + header + "\n" + boxStyle.Render(passwordStyle.Render(m.exportPath)) + "\n"
	}

	header := countStyle.Render(strconv.Itoa(n)) + labelStyle.Render(noun+" generated")

	lines := make([]string, n)
	for i, p := range m.passwords {
		lines[i] = bulletStyle.Render("› ") + passwordStyle.Render(p)
	}

	return "\n" + header + "\n" + boxStyle.Render(strings.Join(lines, "\n")) + "\n"
}

func main() {
	length := flag.Int("l", 16, "password length")
	count := flag.Int("n", 1, "number of passwords to generate")
	export := flag.String("export", "", "export the passwords to this file")
	flag.Parse()

	maxCount := maxConsoleCount
	if *export != "" {
		maxCount = maxExportCount
	}
	if *count > maxCount {
		fmt.Fprintf(os.Stderr, "error: count must be at most %d", maxCount)
		if *export == "" {
			fmt.Fprintf(os.Stderr, " (up to %d with -export)", maxExportCount)
		}
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	passwords, err := GenerateMany(*length, *count)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	m := model{passwords: passwords}
	if *export != "" {
		path, err := exportPasswords(*export, passwords)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		m.exportPath = path
	}

	if _, err := tea.NewProgram(m, tea.WithInput(nil)).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// exportPasswords writes the passwords to path as "1: password" lines
// and returns the absolute path of the file.
func exportPasswords(path string, passwords []string) (string, error) {
	var b strings.Builder
	for i, p := range passwords {
		fmt.Fprintf(&b, "%d: %s\n", i+1, p)
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		return "", err
	}
	return filepath.Abs(path)
}
