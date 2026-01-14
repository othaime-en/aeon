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
	zones       []Zone
	currentView viewType
	width       int
	height      int

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

	// Setup meeting input
	mi := textinput.New()
	mi.Placeholder = "e.g., NYC, London, Tokyo"
	mi.CharLimit = 100
	mi.Width = 50

	return model{
		zones:        []Zone{local, utcZone},
		currentView:  clockView,
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
		// Global shortcuts
		switch msg.String() {
		case "q", "esc":
			if m.convertActive {
				m.convertActive = false
				m.convertInput.Blur()
				return m, nil
			}
			if m.meetingActive {
				m.meetingActive = false
				m.meetingInput.Blur()
				return m, nil
			}
			return m, tea.Quit

		case "tab", "right":
			if !m.convertActive && !m.meetingActive {
				m.currentView = (m.currentView + 1) % 3
			}

		case "shift+tab", "left":
			if !m.convertActive && !m.meetingActive {
				m.currentView = (m.currentView + 2) % 3
			}

		case "1":
			if !m.convertActive && !m.meetingActive {
				m.currentView = clockView
			}

		case "2", "c":
			if !m.convertActive && !m.meetingActive {
				m.currentView = convertView
			}

		case "3", "m":
			if !m.convertActive && !m.meetingActive {
				m.currentView = meetingView
			}

		case "enter":
			if m.currentView == convertView && m.convertActive {
				m.convertResult = m.processConversion(m.convertInput.Value())
				m.convertActive = false
				m.convertInput.Blur()
				m.convertInput.SetValue("")
			} else if m.currentView == meetingView && m.meetingActive {
				m.meetingResult = m.processMeeting(m.meetingInput.Value())
				m.meetingActive = false
				m.meetingInput.Blur()
				m.meetingInput.SetValue("")
			} else if m.currentView == convertView {
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

	// Update inputs
	if m.convertActive {
		m.convertInput, cmd = m.convertInput.Update(msg)
	}
	if m.meetingActive {
		m.meetingInput, cmd = m.meetingInput.Update(msg)
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
		return "←/→ or Tab: Switch views  •  1/2/3: Jump to view  •  q: Quit"
	case convertView:
		return "c: Convert  •  ←/→: Switch views  •  q: Quit"
	case meetingView:
		return "m: Meeting  •  ←/→: Switch views  •  q: Quit"
	}
	return ""
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

	// Load locations
	sourceLoc, err := loadLocation(sourceZone)
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: Unknown source zone '%s'", sourceZone))
	}

	targetLoc, err := loadLocation(targetZone)
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: Unknown target zone '%s'", targetZone))
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
		loc, err := loadLocation(zoneName)
		if err != nil {
			b.WriteString(fmt.Sprintf("⚠️  %s: Unknown zone\n", zoneName))
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

// cityToTimezone maps common city names and abbreviations to IANA timezone identifiers
var cityToTimezone = map[string]string{
	// North America
	"nyc":           "America/New_York",
	"new york":      "America/New_York",
	"newyork":       "America/New_York",
	"la":            "America/Los_Angeles",
	"los angeles":   "America/Los_Angeles",
	"losangeles":    "America/Los_Angeles",
	"sf":            "America/Los_Angeles",
	"san francisco": "America/Los_Angeles",
	"chicago":       "America/Chicago",
	"denver":        "America/Denver",
	"phoenix":       "America/Phoenix",
	"seattle":       "America/Los_Angeles",
	"boston":        "America/New_York",
	"miami":         "America/New_York",
	"toronto":       "America/Toronto",
	"vancouver":     "America/Vancouver",
	"montreal":      "America/Toronto",
	"mexico city":   "America/Mexico_City",
	"mexicocity":    "America/Mexico_City",

	// South America
	"sao paulo":    "America/Sao_Paulo",
	"saopaulo":     "America/Sao_Paulo",
	"buenos aires": "America/Argentina/Buenos_Aires",
	"buenosaires":  "America/Argentina/Buenos_Aires",
	"lima":         "America/Lima",
	"bogota":       "America/Bogota",
	"santiago":     "America/Santiago",

	// Europe
	"london":     "Europe/London",
	"paris":      "Europe/Paris",
	"berlin":     "Europe/Berlin",
	"madrid":     "Europe/Madrid",
	"rome":       "Europe/Rome",
	"amsterdam":  "Europe/Amsterdam",
	"brussels":   "Europe/Brussels",
	"vienna":     "Europe/Vienna",
	"zurich":     "Europe/Zurich",
	"stockholm":  "Europe/Stockholm",
	"oslo":       "Europe/Oslo",
	"copenhagen": "Europe/Copenhagen",
	"dublin":     "Europe/Dublin",
	"lisbon":     "Europe/Lisbon",
	"athens":     "Europe/Athens",
	"istanbul":   "Europe/Istanbul",
	"moscow":     "Europe/Moscow",
	"warsaw":     "Europe/Warsaw",
	"prague":     "Europe/Prague",
	"budapest":   "Europe/Budapest",

	// Asia
	"tokyo":        "Asia/Tokyo",
	"beijing":      "Asia/Shanghai",
	"shanghai":     "Asia/Shanghai",
	"hong kong":    "Asia/Hong_Kong",
	"hongkong":     "Asia/Hong_Kong",
	"singapore":    "Asia/Singapore",
	"dubai":        "Asia/Dubai",
	"mumbai":       "Asia/Kolkata",
	"delhi":        "Asia/Kolkata",
	"bangalore":    "Asia/Kolkata",
	"kolkata":      "Asia/Kolkata",
	"bangkok":      "Asia/Bangkok",
	"jakarta":      "Asia/Jakarta",
	"manila":       "Asia/Manila",
	"seoul":        "Asia/Seoul",
	"taipei":       "Asia/Taipei",
	"kuala lumpur": "Asia/Kuala_Lumpur",
	"karachi":      "Asia/Karachi",
	"tehran":       "Asia/Tehran",
	"riyadh":       "Asia/Riyadh",

	// Africa
	"cairo":         "Africa/Cairo",
	"johannesburg":  "Africa/Johannesburg",
	"lagos":         "Africa/Lagos",
	"nairobi":       "Africa/Nairobi",
	"casablanca":    "Africa/Casablanca",
	"algiers":       "Africa/Algiers",
	"tunis":         "Africa/Tunis",
	"accra":         "Africa/Accra",
	"addis ababa":   "Africa/Addis_Ababa",
	"addisababa":    "Africa/Addis_Ababa",
	"dar es salaam": "Africa/Dar_es_Salaam",
	"kinshasa":      "Africa/Kinshasa",

	// Oceania
	"sydney":     "Australia/Sydney",
	"melbourne":  "Australia/Melbourne",
	"brisbane":   "Australia/Brisbane",
	"perth":      "Australia/Perth",
	"adelaide":   "Australia/Adelaide",
	"auckland":   "Pacific/Auckland",
	"wellington": "Pacific/Auckland",

	// Common timezone abbreviations
	"est":  "America/New_York",
	"cst":  "America/Chicago",
	"mst":  "America/Denver",
	"pst":  "America/Los_Angeles",
	"utc":  "UTC",
	"gmt":  "Europe/London",
	"bst":  "Europe/London",
	"cet":  "Europe/Paris",
	"ist":  "Asia/Kolkata",
	"jst":  "Asia/Tokyo",
	"aest": "Australia/Sydney",
}

func loadLocation(name string) (*time.Location, error) {
	// Normalize the input
	normalized := strings.ToLower(strings.TrimSpace(name))
	normalized = strings.ReplaceAll(normalized, "_", " ")

	// First, check the city mapping
	if tz, ok := cityToTimezone[normalized]; ok {
		return time.LoadLocation(tz)
	}

	// Try variations with the original name
	variations := []string{
		name,
		strings.ReplaceAll(name, " ", "_"),
		"America/" + strings.ReplaceAll(name, " ", "_"),
		"Europe/" + strings.ReplaceAll(name, " ", "_"),
		"Asia/" + strings.ReplaceAll(name, " ", "_"),
		"Africa/" + strings.ReplaceAll(name, " ", "_"),
		"Australia/" + strings.ReplaceAll(name, " ", "_"),
		"Pacific/" + strings.ReplaceAll(name, " ", "_"),
	}

	for _, v := range variations {
		if loc, err := time.LoadLocation(v); err == nil {
			return loc, nil
		}
	}

	// Provide helpful error message with suggestions
	suggestions := getSuggestions(normalized)
	if len(suggestions) > 0 {
		return nil, fmt.Errorf("unknown location '%s'. Did you mean: %s?", name, strings.Join(suggestions, ", "))
	}

	return nil, fmt.Errorf("unknown location: %s", name)
}

// getSuggestions finds similar city names for typos
func getSuggestions(input string) []string {
	var suggestions []string
	for city := range cityToTimezone {
		if strings.HasPrefix(city, input) || strings.Contains(city, input) {
			suggestions = append(suggestions, city)
			if len(suggestions) >= 3 {
				break
			}
		}
	}
	return suggestions
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
