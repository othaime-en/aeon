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
	zones         []Zone
	currentView   viewType
	convertInput  textinput.Model
	convertResult string
	convertActive bool
	err           error
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

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	resultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")).
			Padding(1, 0)
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

	// Setup convert input
	ti := textinput.New()
	ti.Placeholder = "e.g., 3pm NYC to Berlin"
	ti.CharLimit = 100
	ti.Width = 50

	return model{
		zones:        []Zone{local, utcZone},
		currentView:  clockView,
		convertInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), textinput.Blink)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			if m.convertActive {
				m.convertActive = false
				m.convertInput.Blur()
				return m, nil
			}
			return m, tea.Quit
		case "tab", "right":
			if !m.convertActive {
				m.currentView = (m.currentView + 1) % 3
			}
		case "shift+tab", "left":
			if !m.convertActive {
				m.currentView = (m.currentView + 2) % 3
			}
		case "1":
			if !m.convertActive {
				m.currentView = clockView
			}
		case "2", "c":
			if !m.convertActive {
				m.currentView = convertView
			}
		case "3":
			if !m.convertActive {
				m.currentView = meetingView
			}
		case "enter":
			if m.currentView == convertView && m.convertActive {
				// TODO: process
				m.convertResult = "TODO: Conversion result"
				m.convertActive = false
				m.convertInput.Blur()
				m.convertInput.SetValue("")
			} else if m.currentView == convertView {
				m.convertActive = true
				m.convertInput.Focus()
			}
		}
	case tickMsg:
		return m, tickCmd()
	}

	if m.convertActive {
		m.convertInput, cmd = m.convertInput.Update(msg)
	}

	return m, cmd
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
		b.WriteString(m.renderConvertView())
	case meetingView:
		b.WriteString("Meeting view (TODO)\n")
	}

	// Help text
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(m.getHelpText()))

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

func (m model) renderConvertView() string {
	var b strings.Builder

	if m.convertActive {
		b.WriteString("Enter conversion query:\n\n")
		b.WriteString(m.convertInput.View())
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press Enter to convert, Esc to cancel"))
	} else {
		b.WriteString("Press Enter to start a conversion\n\n")
		if m.convertResult != "" {
			b.WriteString(resultStyle.Render(m.convertResult))
		}
	}

	return b.String()
}

func (m model) getHelpText() string {
	switch m.currentView {
	case clockView:
		return "←/→ or Tab: Switch views  •  1/2/3: Jump to view  •  q: Quit"
	case convertView:
		return "c: Convert  •  ←/→: Switch views  •  q: Quit"
	case meetingView:
		return "←/→: Switch views  •  q: Quit"
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
