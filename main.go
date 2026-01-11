package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Zone represents a time zone configuration
type Zone struct {
	Name     string
	Location *time.Location
}

// Model holds the application state
type model struct {
	zones []Zone
	err   error
}

// TickMsg is sent every second to update clocks
type tickMsg time.Time

func initialModel() model {
	// Initialize with local time and UTC
	local := Zone{
		Name:     "Local",
		Location: time.Local,
	}
	
	utc, _ := time.LoadLocation("UTC")
	utcZone := Zone{
		Name:     "UTC",
		Location: utc,
	}
	
	return model{
		zones: []Zone{local, utcZone},
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" {
			return m, tea.Quit
		}
	case tickMsg:
		return m, tickCmd()
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	
	// Title
	b.WriteString("⏰ aeon - Time Zone Manager\n\n")
	
	// Clock view
	now := time.Now()
	for _, zone := range m.zones {
		t := now.In(zone.Location)
		timeStr := t.Format("15:04:05")
		dateStr := t.Format("Mon, Jan 02")
		offset := t.Format("-07:00")
		line := fmt.Sprintf("%-15s  %s  %s  (UTC%s)", 
			zone.Name, timeStr, dateStr, offset)
		b.WriteString(line + "\n")
	}
	
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}