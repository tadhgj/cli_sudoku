package main

import (
	tea "charm.land/bubbletea/v2"
)

func (s *SudokuGameWrapperState) WinBackToDifficulty() {
	s.shownPage = difficultyPage
}

func (s *SudokuGameWrapperState) WinUpdate(msg tea.Msg) {
	// if keypress is s
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			s.WinBackToDifficulty()
		}

	}

}

func (s *SudokuGameWrapperState) WinView() tea.View {

	var pauseText string

	gameName := "Tadhg-doku"

	pauseText += gameName + "\n"

	pauseText += "\n"

	pauseText += "You won!\n"

	pauseText += "\n"

	// print type of game (ez. med. hard.)
	// print how long user has been playing game

	// print controls
	pauseSelectControl := Control{
		keys: "enter",
		desc: "back to diff.",
	}
	pauseText += s.RenderControlList([]Control{pauseSelectControl, QuitControl()})

	v := tea.NewView(pauseText)
	v.AltScreen = true
	return v
}
