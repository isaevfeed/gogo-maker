package logger

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Base styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")). // Bright white
			MarginBottom(1)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196")) // Red

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46")) // Green

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")) // Cyan

	// Command styles
	commandStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214")) // Orange

	flagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("183")) // Purple

	exampleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")) // Gray

	// Highlight styles
	highlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("51")) // Bright cyan

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")) // Gray
)

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Title(text string) {
	fmt.Println(titleStyle.Render(text))
}

func (l *Logger) Info(text string) {
	fmt.Println(infoStyle.Render("  " + text))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	fmt.Println(infoStyle.Render("  " + text))
}

func (l *Logger) Success(text string) {
	fmt.Println(successStyle.Render("✓ " + text))
}

func (l *Logger) Successf(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	fmt.Println(successStyle.Render("✓ " + text))
}

func (l *Logger) Error(err error) {
	fmt.Println(errorStyle.Render("✗ " + err.Error()))
}

func (l *Logger) Fatal(err error) {
	l.Error(err)
	panic(err)
}

func (l *Logger) Command(cmd, desc string) {
	fmt.Printf("  %s  %s\n",
		commandStyle.Render(cmd),
		dimStyle.Render(desc),
	)
}

func (l *Logger) Example(cmd string) {
	fmt.Printf("  %s\n", exampleStyle.Render(cmd))
}

func (l *Logger) Commands(title string, commands []CommandInfo) {
	l.Title(title)
	fmt.Println()

	for _, cmd := range commands {
		l.Command(cmd.Cmd, cmd.Desc)
	}

	fmt.Println()
	for _, cmd := range commands {
		if len(cmd.Examples) > 0 {
			l.Info("Examples:")
			for _, ex := range cmd.Examples {
				l.Example(ex)
			}
		}
	}
}

type CommandInfo struct {
	Cmd      string
	Desc     string
	Examples []string
}