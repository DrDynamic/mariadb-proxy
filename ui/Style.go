package ui

import "github.com/charmbracelet/lipgloss"

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var focusStyle = baseStyle.
	BorderForeground(lipgloss.Color("245"))

func getStyle(focus bool) lipgloss.Style {
	if focus {
		return focusStyle
	} else {
		return baseStyle
	}
}
