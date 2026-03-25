package main

import tea "charm.land/bubbletea/v2"

func LoadingView() tea.View {
	v := tea.NewView("Loading...")
	v.AltScreen = true
	return v
}
