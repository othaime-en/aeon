package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View types
type viewType int

const (
	clockView viewType = iota
	convertView
	meetingView
)

// Zone represents a time zone configuration
type Zone struct {
	Name     string
	Location *time.Location
}

// Model holds the application state
type model struct {
	zones       []Zone
	currentView viewType
	err         error
}

// TickMsg is sent every second to update clocks
type tickMsg time.Time

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginBottom(1)

	tabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			Background(lipgloss.Color("235")).
			Padding(0, 2)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Padding(0, 2)

	clockStyle = lipgloss.NewStyle().
			Padding(0, 1).
			MarginBottom(1)
)

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
		zones:       []Zone{local, utcZone},
		currentView: clockView,
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
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "tab", "right":
			m.currentView = (m.currentView + 1) % 3
		case "shift+tab", "left":
			m.currentView = (m.currentView + 2) % 3
		case "1":
			m.currentView = clockView
		case "2":
			m.currentView = convertView
		case "3":
			m.currentView = meetingView
		}
	case tickMsg:
		return m, tickCmd()
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("⏰ aeon - Time Zone Manager"))
	b.WriteString("\n\n")

	// Tabs
	tabs := []string{}
	for i, name := range []string{"Clock", "Convert", "Meeting"} {
		if viewType(i) == m.currentView {
			tabs = append(tabs, tabActiveStyle.Render(name))
		} else {
			tabs = append(tabs, tabInactiveStyle.Render(name))
		}
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
	b.WriteString("\n\n")

	// View content
	switch m.currentView {
	case clockView:
		b.WriteString(m.renderClockView())
	case convertView:
		b.WriteString("Convert view (TODO)\n")
	case meetingView:
		b.WriteString("Meeting view (TODO)\n")
	}

	return b.String()
}

func (m model) renderClockView() string {
	var b strings.Builder
	now := time.Now()

	for _, zone := range m.zones {
		t := now.In(zone.Location)

		// Format: Zone Name    HH:MM:SS    Day, Mon DD
		timeStr := t.Format("15:04:05")
		dateStr := t.Format("Mon, Jan 02")
		offset := t.Format("-07:00")

		line := fmt.Sprintf("%-15s  %s  %s  (UTC%s)",
			zone.Name, timeStr, dateStr, offset)
		b.WriteString(clockStyle.Render(line))
		b.WriteString("\n")
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
