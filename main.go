package main

import (
	"aeon/timezones"
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
	zones        []Zone
	currentView  viewType
	selectedZone int
	width        int
	height       int

	// Add zone state
	addZoneActive bool
	addZoneInput  textinput.Model

	// Convert view state
	convertInput  textinput.Model
	convertResult string
	convertActive bool

	// Meeting view state
	meetingInput  textinput.Model
	meetingResult string
	meetingActive bool

	err error
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

	clockSelectedStyle = lipgloss.NewStyle().
				Padding(0, 1).
				MarginBottom(1).
				Background(lipgloss.Color("235"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	resultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")).
			Padding(1, 0)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")).
			Bold(true)
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

	// Setup add zone input
	azi := textinput.New()
	azi.Placeholder = "Enter city name or timezone (e.g., NYC, Berlin, Asia/Tokyo)"
	azi.CharLimit = 100
	azi.Width = 60

	// Setup convert input
	ti := textinput.New()
	ti.Placeholder = "e.g., 3pm NYC to Berlin"
	ti.CharLimit = 100
	ti.Width = 50

	// Setup meeting input
	mi := textinput.New()
	mi.Placeholder = "e.g., NYC, London, Tokyo"
	mi.CharLimit = 100
	mi.Width = 50

	return model{
		zones:        []Zone{local, utcZone},
		currentView:  clockView,
		selectedZone: 0,
		addZoneInput: azi,
		convertInput: ti,
		meetingInput: mi,
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
		// Handle add zone input mode
		if m.addZoneActive {
			switch msg.String() {
			case "esc":
				m.addZoneActive = false
				m.addZoneInput.Blur()
				m.addZoneInput.SetValue("")
				return m, nil
			case "enter":
				zoneName := strings.TrimSpace(m.addZoneInput.Value())
				if zoneName != "" {
					if err := m.addZone(zoneName); err != nil {
						m.err = err
					} else {
						m.err = nil
					}
				}
				m.addZoneActive = false
				m.addZoneInput.Blur()
				m.addZoneInput.SetValue("")
				return m, nil
			default:
				m.addZoneInput, cmd = m.addZoneInput.Update(msg)
				return m, cmd
			}
		}

		// Handle convert input mode
		if m.convertActive {
			switch msg.String() {
			case "esc":
				m.convertActive = false
				m.convertInput.Blur()
				return m, nil
			case "enter":
				m.convertResult = m.processConversion(m.convertInput.Value())
				m.convertActive = false
				m.convertInput.Blur()
				m.convertInput.SetValue("")
				return m, nil
			default:
				m.convertInput, cmd = m.convertInput.Update(msg)
				return m, cmd
			}
		}

		// Handle meeting input mode
		if m.meetingActive {
			switch msg.String() {
			case "esc":
				m.meetingActive = false
				m.meetingInput.Blur()
				return m, nil
			case "enter":
				m.meetingResult = m.processMeeting(m.meetingInput.Value())
				m.meetingActive = false
				m.meetingInput.Blur()
				m.meetingInput.SetValue("")
				return m, nil
			default:
				m.meetingInput, cmd = m.meetingInput.Update(msg)
				return m, cmd
			}
		}

		// Global shortcuts (when not in input mode)
		switch msg.String() {
		case "q":
			return m, tea.Quit

		case "tab", "right":
			m.currentView = (m.currentView + 1) % 3
			m.selectedZone = 0

		case "shift+tab", "left":
			m.currentView = (m.currentView + 2) % 3
			m.selectedZone = 0

		case "1":
			m.currentView = clockView
			m.selectedZone = 0

		case "2", "c":
			m.currentView = convertView

		case "3", "m":
			m.currentView = meetingView

		case "a":
			if m.currentView == clockView {
				m.addZoneActive = true
				m.addZoneInput.Focus()
				m.err = nil
			}

		case "d":
			if m.currentView == clockView && len(m.zones) > 1 {
				m.deleteSelectedZone()
			}

		case "up", "k":
			if m.currentView == clockView && m.selectedZone > 0 {
				m.selectedZone--
			}

		case "down", "j":
			if m.currentView == clockView && m.selectedZone < len(m.zones)-1 {
				m.selectedZone++
			}

		case "enter":
			if m.currentView == convertView {
				m.convertActive = true
				m.convertInput.Focus()
			} else if m.currentView == meetingView {
				m.meetingActive = true
				m.meetingInput.Focus()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		return m, tickCmd()
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
		b.WriteString(m.renderMeetingView())
	}

	// Error message
	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	// Help text
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(m.getHelpText()))

	return b.String()
}

func (m model) renderClockView() string {
	var b strings.Builder
	now := time.Now()

	if m.addZoneActive {
		b.WriteString("Add timezone:\n\n")
		b.WriteString(m.addZoneInput.View())
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press Enter to add, Esc to cancel"))
		return b.String()
	}

	for i, zone := range m.zones {
		t := now.In(zone.Location)

		// Format: Zone Name    HH:MM:SS    Day, Mon DD
		timeStr := t.Format("15:04:05")
		dateStr := t.Format("Mon, Jan 02")
		offset := t.Format("-07:00")

		line := fmt.Sprintf("%-15s  %s  %s  (UTC%s)",
			zone.Name, timeStr, dateStr, offset)

		if i == m.selectedZone {
			b.WriteString(clockSelectedStyle.Render(line))
		} else {
			b.WriteString(clockStyle.Render(line))
		}
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

func (m model) renderMeetingView() string {
	var b strings.Builder

	if m.meetingActive {
		b.WriteString("Enter time zones (comma-separated):\n\n")
		b.WriteString(m.meetingInput.View())
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press Enter to find meeting slots, Esc to cancel"))
	} else {
		b.WriteString("Press Enter to find meeting slots\n\n")
		if m.meetingResult != "" {
			b.WriteString(resultStyle.Render(m.meetingResult))
		}
	}

	return b.String()
}

func (m model) getHelpText() string {
	switch m.currentView {
	case clockView:
		if m.addZoneActive {
			return ""
		}
		return "a: Add zone  •  d: Delete  •  ↑/↓: Select  •  ←/→: Switch views  •  q: Quit"
	case convertView:
		return "Enter: Convert  •  ←/→: Switch views  •  q: Quit"
	case meetingView:
		return "Enter: Find slots  •  ←/→: Switch views  •  q: Quit"
	}
	return ""
}

func (m *model) addZone(name string) error {
	loc, err := timezones.Resolve(name)
	if err != nil {
		return err
	}

	// Check if zone already exists
	for _, z := range m.zones {
		if z.Location.String() == loc.String() {
			return fmt.Errorf("zone '%s' already added", name)
		}
	}

	m.zones = append(m.zones, Zone{
		Name:     name,
		Location: loc,
	})

	return nil
}

func (m *model) deleteSelectedZone() {
	if m.selectedZone >= 0 && m.selectedZone < len(m.zones) {
		m.zones = append(m.zones[:m.selectedZone], m.zones[m.selectedZone+1:]...)
		if m.selectedZone >= len(m.zones) {
			m.selectedZone = len(m.zones) - 1
		}
		if m.selectedZone < 0 {
			m.selectedZone = 0
		}
	}
}

func (m model) processConversion(input string) string {
	// Simple parser for "3pm NYC to Berlin" format
	input = strings.TrimSpace(input)
	if input == "" {
		return errorStyle.Render("Error: Empty input")
	}

	// Basic parsing (MVP - simple approach)
	parts := strings.Split(input, " to ")
	if len(parts) != 2 {
		return errorStyle.Render("Error: Use format '3pm NYC to Berlin'")
	}

	sourcePart := strings.TrimSpace(parts[0])
	targetZone := strings.TrimSpace(parts[1])

	// Parse source time and zone
	sourceWords := strings.Fields(sourcePart)
	if len(sourceWords) < 2 {
		return errorStyle.Render("Error: Specify time and source zone")
	}

	timeStr := sourceWords[0]
	sourceZone := strings.Join(sourceWords[1:], " ")

	// Load locations using the timezones package
	sourceLoc, err := timezones.Resolve(sourceZone)
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", err))
	}

	targetLoc, err := timezones.Resolve(targetZone)
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", err))
	}

	// Parse time (simple 12/24 hour format)
	sourceTime, err := parseTime(timeStr)
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: Invalid time format '%s'", timeStr))
	}

	// Convert
	now := time.Now()
	source := time.Date(now.Year(), now.Month(), now.Day(),
		sourceTime.Hour(), sourceTime.Minute(), 0, 0, sourceLoc)
	target := source.In(targetLoc)

	return fmt.Sprintf("%s in %s  →  %s in %s",
		source.Format("3:04 PM Mon Jan 02"),
		sourceZone,
		target.Format("3:04 PM Mon Jan 02"),
		targetZone,
	)
}

func (m model) processMeeting(input string) string {
	zones := strings.Split(input, ",")
	if len(zones) < 2 {
		return errorStyle.Render("Error: Enter at least 2 zones")
	}

	// For MVP, show business hours for each zone
	var b strings.Builder
	b.WriteString("Business Hours (9 AM - 5 PM local):\n\n")

	now := time.Now()
	for _, zoneName := range zones {
		zoneName = strings.TrimSpace(zoneName)
		loc, err := timezones.Resolve(zoneName)
		if err != nil {
			b.WriteString(fmt.Sprintf("⚠️  %s: %v\n", zoneName, err))
			continue
		}

		t := now.In(loc)
		start := time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, loc)
		end := time.Date(t.Year(), t.Month(), t.Day(), 17, 0, 0, 0, loc)

		b.WriteString(fmt.Sprintf("%-15s: %s - %s\n",
			zoneName,
			start.In(time.Local).Format("3:04 PM"),
			end.In(time.Local).Format("3:04 PM"),
		))
	}

	return b.String()
}

func parseTime(s string) (time.Time, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	// Try common formats
	formats := []string{
		"3pm",
		"3:04pm",
		"15:04",
		"3PM",
		"3:04PM",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid time format")
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
