package main

import (
	"strconv"

	tea "charm.land/bubbletea/v2"
)

// enumerate difficulties

func (s *SudokuGameWrapperState) BackToGame() {
	s.shownPage = gamePage
}

func (s *SudokuGameWrapperState) ToggleCenterOfUniverse() {
	s.userOptions.centerOfUniverseRenderingStyle = !s.userOptions.centerOfUniverseRenderingStyle
}

func (s *SudokuGameWrapperState) PauseUpdate(msg tea.Msg) {
	// if keypress is s
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "c":
			s.ToggleCenterOfUniverse()

		case "p":
			s.BackToGame()
		}

	}

}

func (s *SudokuGameWrapperState) PauseView() tea.View {

	var pauseText string

	gameName := "Tadhg-doku"

	pauseText += gameName + "\n"

	pauseText += "\n"

	pauseText += "Paused\n"

	pauseText += "\n"

	// print type of game (ez. med. hard.)
	// print how long user has been playing game

	// print user settings
	pauseText += "User Settings:\n"

	pauseText += "Center of Universe (press c to toggle): "
	pauseText += strconv.FormatBool(s.userOptions.centerOfUniverseRenderingStyle)
	pauseText += "\n"

	pauseText += "\n"
	// print controls
	pauseSelectControl := Control{
		keys: "p",
		desc: "play",
	}
	pauseText += s.RenderControlList([]Control{pauseSelectControl, QuitControl()})

	v := tea.NewView(pauseText)
	v.AltScreen = true
	return v
}
