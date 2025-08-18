package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Chart calculation helpers

// calculateRequiredChartsHeight calculates how much vertical space the charts need
func (m *DashboardModel) calculateRequiredChartsHeight() int {
	// Calculate content needs for each chart
	wordsContentLines := m.calculateWordsContentLines()
	attrsContentLines := m.calculateAttributesContentLines()
	distContentLines := m.calculateDistributionContentLines()
	countsContentLines := m.calculateCountsContentLines()

	// More precise height calculation: title + content + borders
	// Each chart needs: title(1) + content(N) + top/bottom borders(2) = N+3
	wordsHeight := wordsContentLines + 3
	attrsHeight := attrsContentLines + 3
	distHeight := distContentLines + 3
	countsHeight := countsContentLines + 3

	// Row heights: use actual maximum needed by each row
	topRowHeight := max(wordsHeight, attrsHeight)
	bottomRowHeight := max(distHeight, countsHeight)

	// Total height: both rows + minimal spacing
	totalRequired := topRowHeight + bottomRowHeight

	// Ensure reasonable bounds but prioritize showing all content
	if totalRequired < 14 {
		totalRequired = 14 // Minimum for functional charts
	}
	if totalRequired > 35 {
		totalRequired = 35 // Increased maximum since we need more space
	}

	return totalRequired
}

// Chart rendering functions

// renderChartsGrid renders the 2x2 grid of charts
func (m *DashboardModel) renderChartsGrid(height int) string {
	if m.width < 20 {
		return "Terminal too narrow"
	}

	// DYNAMIC CHART SIZING: Use the same calculation as in calculateRequiredChartsHeight
	// Calculate actual content needs
	wordsContentLines := m.calculateWordsContentLines()
	attrsContentLines := m.calculateAttributesContentLines()
	distContentLines := m.calculateDistributionContentLines()
	countsContentLines := m.calculateCountsContentLines()

	// Calculate precise heights for each chart: content + title + borders
	wordsHeight := wordsContentLines + 1
	attrsHeight := attrsContentLines + 1
	distHeight := distContentLines + 1
	countsHeight := countsContentLines + 1

	// 2x2 grid: top row height = max(left, right), bottom row height = max(left, right)
	topRowHeight := max(wordsHeight, attrsHeight)
	bottomRowHeight := max(distHeight, countsHeight)

	// Use nearly all available width for charts
	// Account for: borders(4) only = 4 chars per chart
	chartWidth := (m.width / 2) - 2 // Use almost all available space
	if chartWidth < 25 {
		chartWidth = 25 // Reasonable minimum for readability
	}

	// Top row with calculated heights
	wordsChart := m.renderWordsChart(chartWidth, topRowHeight)
	attrsChart := m.renderAttributesChart(chartWidth, topRowHeight)
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, wordsChart, attrsChart)

	// Bottom row with calculated heights
	// Use drain3 chart instead of distribution chart, but keep distribution code for future use
	distChart := m.renderDrain3Chart(chartWidth, bottomRowHeight)
	countsChart := m.renderCountsChart(chartWidth, bottomRowHeight)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, distChart, countsChart)

	// Combine rows - apply strict height constraint to prevent overflow
	result := lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)

	// Don't force height - let content determine size
	constrainedStyle := lipgloss.NewStyle().
		MaxHeight(height). // Prevent exceeding allocation but don't force
		Width(m.width)

	return constrainedStyle.Render(result)
}
