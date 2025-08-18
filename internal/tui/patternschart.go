package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderDrain3Chart renders the drain3 pattern extraction chart
func (m *DashboardModel) renderDrain3Chart(width, height int) string {
	// Use MaxHeight instead of Height to prevent empty space
	style := sectionStyle.Width(width).Height(height)
	if m.activeSection == SectionDistribution {
		style = activeSectionStyle.Width(width).Height(height)
	}

	// Get pattern stats for the title
	patternCount, totalLogs := 0, 0
	if m.drain3Manager != nil {
		patternCount, totalLogs = m.drain3Manager.GetStats()
	}

	// Build title with stats
	titleText := "Log Patterns"
	if patternCount > 0 {
		titleText = fmt.Sprintf("Log Patterns (%d patterns from %d logs)", patternCount, totalLogs)
	}
	title := chartTitleStyle.Render(titleText)

	var content string
	if m.drain3Manager != nil && patternCount > 0 {
		content = m.renderDrain3Content(width)
	} else {
		content = helpStyle.Render("Extracting patterns")
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left, title, content))
}

// renderDrain3Content renders the drain3 pattern list
func (m *DashboardModel) renderDrain3Content(chartWidth int) string {
	if m.drain3Manager == nil {
		return helpStyle.Render("Pattern extraction not available")
	}

	// Get top patterns - limit to 8 for display
	patterns := m.drain3Manager.GetTopPatterns(8)

	// Calculate the maximum count for bar scaling
	maxCount := 0
	for _, p := range patterns {
		if p.Count > maxCount {
			maxCount = p.Count
		}
	}

	// Build the pattern list with mini bar charts
	var lines []string

	// Calculate available width for the template text
	// Format: [bar] count% | template
	// Reserve space for: bar(15) + count%(8) + separators(3) = 26
	templateWidth := chartWidth - 26
	if templateWidth < 20 {
		templateWidth = 20
	}

	// Always show exactly 8 lines to match log counts chart
	const displayLines = 8

	for i := 0; i < displayLines; i++ {
		if i < len(patterns) {
			// We have a pattern to display
			pattern := patterns[i]

			// Create a mini bar for each pattern
			barWidth := 12
			fillWidth := int(float64(pattern.Count) * float64(barWidth) / float64(maxCount))
			if fillWidth == 0 && pattern.Count > 0 {
				fillWidth = 1 // Ensure at least 1 char for non-zero counts
			}

			// Create the bar visual
			bar := strings.Repeat("█", fillWidth) + strings.Repeat("░", barWidth-fillWidth)

			// Format percentage
			percentage := fmt.Sprintf("%5.1f%%", pattern.Percentage)

			// Truncate template if needed
			template := pattern.Template
			if len(template) > templateWidth {
				template = template[:templateWidth-3] + "..."
			}

			// Color code based on frequency (high frequency = more important)
			var barColor lipgloss.Style
			if i < 3 {
				// Top 3 patterns - red/orange for high frequency
				barColor = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
			} else if i < 6 {
				// Middle patterns - yellow
				barColor = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
			} else {
				// Lower patterns - blue/gray
				barColor = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
			}

			// Format the line
			line := fmt.Sprintf("%s %s │ %s",
				barColor.Render(bar),
				lipgloss.NewStyle().Foreground(ColorGray).Render(percentage),
				lipgloss.NewStyle().Foreground(ColorWhite).Render(template),
			)

			lines = append(lines, line)
		} else {
			// No pattern for this line - add empty line to maintain 7-line format
			emptyBar := strings.Repeat("░", 12)
			grayStyle := lipgloss.NewStyle().Foreground(ColorGray)
			line := fmt.Sprintf("%s %s  │ %s",
				grayStyle.Render(emptyBar),
				grayStyle.Render("     "),
				grayStyle.Render("(no pattern)"),
			)
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}
