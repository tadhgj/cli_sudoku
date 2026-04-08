package main

import tea "charm.land/bubbletea/v2"

// enumerate difficulties

func (s *SudokuGameWrapperState) BackToGame() {
	s.shownPage = gamePage
}

func (s *SudokuGameWrapperState) PauseUpdate(msg tea.Msg) {
	// if keypress is s
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "p":
			s.BackToGame()
		}

	}

}

func (s *SudokuGameWrapperState) PauseView() tea.View {

	var pauseText string

	gameName := "Tadhg-doku"

	pauseText += gameName + "\n\n"

	pauseText += "Paused"

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
