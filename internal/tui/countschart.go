package tui

import (
	"fmt"
	"strings"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
)

// calculateCountsContentLines calculates lines needed for counts chart content
func (m *DashboardModel) calculateCountsContentLines() int {
	if len(m.countsHistory) == 0 {
		return 1 // "No data available"
	}

	// Bar chart height - updated to match renderCountsContent logic
	chartHeight := 8 // Increased to accommodate TOTAL row in legend
	if m.width < 80 {
		chartHeight = 6 // Reduced proportionally for narrow terminals
	}
	return chartHeight // Chart now includes legend on right side, no extra line needed
}

// renderCountsChart renders the line counts chart
func (m *DashboardModel) renderCountsChart(width, height int) string {
	// Use MaxHeight instead of Height to prevent empty space
	style := sectionStyle.Width(width).Height(height)
	if m.activeSection == SectionCounts {
		style = activeSectionStyle.Width(width).Height(height)
	}

	// Create header with title on left and min/max/latest stats on right
	var headerText string
	if len(m.countsHistory) > 0 {
		latest := m.countsHistory[len(m.countsHistory)-1]

		// Calculate min/max totals across history
		minTotal, maxTotal := latest.Total, latest.Total
		for _, counts := range m.countsHistory {
			if counts.Total < minTotal {
				minTotal = counts.Total
			}
			if counts.Total > maxTotal {
				maxTotal = counts.Total
			}
		}

		// Create left and right parts of header
		leftTitle := "Log Counts"
		rightStats := fmt.Sprintf("Min: %d | Max: %d", minTotal, maxTotal)

		// Calculate available space (account for borders and padding)
		availableWidth := width - 4
		rightStatsWidth := len(rightStats)
		leftTitleWidth := len(leftTitle)
		spacerWidth := availableWidth - leftTitleWidth - rightStatsWidth

		if spacerWidth > 0 {
			headerText = leftTitle + strings.Repeat(" ", spacerWidth) + rightStats
		} else {
			// Fallback if not enough space - just show title
			headerText = leftTitle
		}
	} else {
		headerText = "Log Counts"
	}

	title := chartTitleStyle.Render(headerText)

	var content string
	if len(m.countsHistory) > 0 {
		content = m.renderCountsContent(width)
	} else {
		content = helpStyle.Render("No data available")
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left, title, content))
}

// renderCountsContent renders a stacked bar chart for log counts by severity
func (m *DashboardModel) renderCountsContent(chartWidth int) string {
	if len(m.countsHistory) == 0 {
		return helpStyle.Render("No data available")
	}

	// Calculate total logs for debugging if needed
	totalLogs := 0
	for _, counts := range m.countsHistory {
		totalLogs += counts.Total
	}

	// Reserve space for legend on the right side
	// Legend needs: "CRITICAL: " (10) + 6 digits (6) + padding (2) = 18 chars
	legendWidth := 18

	// Determine chart dimensions - use dynamic width but reserve space for legend
	chartHeight := 8                                 // Increased by 1 more for TOTAL row in legend
	actualChartWidth := chartWidth - legendWidth - 2 // Reserve space for legend + separator
	if actualChartWidth < 20 {
		actualChartWidth = 20 // Minimum width for chart
	}
	if m.width < 80 {
		chartHeight = 6 // Reduced proportionally for narrow terminals
	}

	// Prepare data for stacked bar chart
	dataPoints := len(m.countsHistory)
	maxBars := actualChartWidth / 3 // Conservative spacing

	// Always show the most recent data points that fit in the available space
	var paddingCount int
	var dataStartIdx int

	if dataPoints < maxBars {
		// Not enough data - pad with zeros on the left
		paddingCount = maxBars - dataPoints
		dataStartIdx = 0
	} else {
		// Enough data - take the most recent points
		paddingCount = 0
		dataStartIdx = dataPoints - maxBars
	}

	// Create ntcharts bar chart with stacked bars
	bc := barchart.New(actualChartWidth, chartHeight,
		barchart.WithBarGap(1),   // Gap of 1 for better visibility
		barchart.WithBarWidth(1), // Conservative width for readability
		barchart.WithNoAxis(),    // Remove axis lines
	)

	// Define severity colors - using more contrasting colors for better stacking visibility
	severityColors := map[string]lipgloss.Style{
		"TRACE":    lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Background(lipgloss.Color("240")), // Dark gray
		"DEBUG":    lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Background(lipgloss.Color("244")), // Medium gray
		"INFO":     lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Background(lipgloss.Color("39")),   // Bright blue
		"WARN":     lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Background(lipgloss.Color("208")), // Bright orange
		"ERROR":    lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Background(lipgloss.Color("196")), // Bright red
		"FATAL":    lipgloss.NewStyle().Foreground(lipgloss.Color("201")).Background(lipgloss.Color("201")), // Magenta/Pink (very distinct)
		"CRITICAL": lipgloss.NewStyle().Foreground(lipgloss.Color("201")).Background(lipgloss.Color("201")), // Magenta/Pink (very distinct)
		"UNKNOWN":  lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("250")), // Light gray
	}

	// Add zero padding (empty bars on the left)
	for i := 0; i < paddingCount; i++ {
		barData := barchart.BarData{
			Label: "",
			Values: []barchart.BarValue{
				{
					Name:  "EMPTY",
					Value: 0,
					Style: severityColors["UNKNOWN"],
				},
			},
		}
		bc.Push(barData)
	}

	// Add actual data (most recent points on the right) - stacked bar approach
	actualDataCount := min(dataPoints, maxBars-paddingCount)
	for i := 0; i < actualDataCount; i++ {
		counts := m.countsHistory[dataStartIdx+i]

		// Create stacked bars with multiple severity levels
		var barValues []barchart.BarValue

		// Define severity order (bottom to top): TRACE → DEBUG → INFO → WARN → ERROR → FATAL
		// Only add segments for severities that have count > 0
		severityData := []struct {
			name  string
			count int
			style lipgloss.Style
		}{
			{"TRACE", counts.Trace, severityColors["TRACE"]},
			{"DEBUG", counts.Debug, severityColors["DEBUG"]},
			{"INFO", counts.Info, severityColors["INFO"]},
			{"WARN", counts.Warn, severityColors["WARN"]},
			{"ERROR", counts.Error, severityColors["ERROR"]},
			{"FATAL", counts.Fatal + counts.Critical, severityColors["FATAL"]}, // Combine FATAL and CRITICAL
		}

		// Add bar segments for each severity level with count > 0
		for _, sev := range severityData {
			if sev.count > 0 {
				barValues = append(barValues, barchart.BarValue{
					Name:  sev.name,
					Value: float64(sev.count),
					Style: sev.style,
				})
			}
		}

		// If no logs in this interval, fill an empty bar to keep logs shifting
		if len(barValues) == 0 {
			barValues = append(barValues, barchart.BarValue{
				Name:  "EMPTY",
				Value: 0.0, // Minimal value to show something
				Style: severityColors["UNKNOWN"],
			})
		}

		barData := barchart.BarData{
			Label:  "",
			Values: barValues,
		}
		bc.Push(barData)
	}

	bc.Draw()
	chartOutput := bc.View()

	// Create vertical legend on the right side
	var legend string
	if len(m.countsHistory) > 0 {
		latest := m.countsHistory[len(m.countsHistory)-1]

		// Define severity levels in priority order with full names
		// Match the stacking order: TRACE → DEBUG → INFO → WARN → ERROR → FATAL (bottom to top)
		// Use same colors as chart for consistency
		severityLevels := []struct {
			name  string
			count int
			color string
		}{
			{"FATAL", latest.Fatal + latest.Critical, "201"}, // Magenta/Pink (matches chart)
			{"ERROR", latest.Error, "196"},                   // Bright red (matches chart)
			{"WARN", latest.Warn, "208"},                     // Bright orange (matches chart)
			{"INFO", latest.Info, "39"},                      // Bright blue (matches chart)
			{"DEBUG", latest.Debug, "244"},                   // Medium gray (matches chart)
			{"TRACE", latest.Trace, "240"},                   // Dark gray (matches chart)
			{"─────", 0, "7"},                                // Separator
			{"TOTAL", latest.Total, "7"},                     // Light gray for total
		}

		// Create legend lines with proper justification
		var legendLines []string

		for _, sev := range severityLevels {
			if sev.name == "─────" {
				// Add separator line without count
				colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(sev.color))
				line := colorStyle.Render("─────────────")
				legendLines = append(legendLines, line)
			} else {
				// Left-justify label, right-justify value with proper spacing
				label := fmt.Sprintf("%-6s:", sev.name) // 6 chars for label + colon
				value := fmt.Sprintf("%6d", sev.count)  // Right-align value in 6 chars

				// Apply color to the entire line
				colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(sev.color))
				line := colorStyle.Render(label + value)
				legendLines = append(legendLines, line)
			}
		}

		// Pad legend to match chart height
		for len(legendLines) < chartHeight {
			legendLines = append(legendLines, strings.Repeat(" ", legendWidth-2))
		}

		// Join legend lines vertically
		legend = strings.Join(legendLines, "\n")
	} else {
		// Empty legend when no data
		legend = strings.Repeat("\n", chartHeight-1)
	}

	// Combine chart and legend horizontally with separator
	separator := strings.Repeat(" ", 2) // 2-space separator

	// Split chart output into lines to align with legend
	chartLines := strings.Split(chartOutput, "\n")

	// Pad chart lines to match desired height
	for len(chartLines) < chartHeight {
		chartLines = append(chartLines, "")
	}

	// Combine each chart line with corresponding legend line
	var combinedLines []string
	legendSplit := strings.Split(legend, "\n")

	for i := 0; i < chartHeight; i++ {
		chartLine := ""
		legendLine := ""

		if i < len(chartLines) {
			chartLine = chartLines[i]
		}
		if i < len(legendSplit) {
			legendLine = legendSplit[i]
		}

		// Ensure chart line has consistent width
		if len(chartLine) < actualChartWidth {
			chartLine += strings.Repeat(" ", actualChartWidth-len(chartLine))
		}

		combined := chartLine + separator + legendLine
		combinedLines = append(combinedLines, combined)
	}

	return strings.Join(combinedLines, "\n")
}
